package storage

import (
	"bytes"
	"crypto/md5"
	"dfs/config"
	"dfs/connection"
	file_io "dfs/files_io"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

type StorageNode struct {
	id                string
	size              big.Int
	storageNodeHost string
	storageNodePort int32
	config *config.Config
	connectionHandler *connection.ConnectionHandler
	running           bool
	server 			*connection.Server
	savePath string

}

func NewStorageNode(id string, size int64, host string, port int32, config *config.Config, savePath string) *StorageNode {
	storageNode := &StorageNode{}
	storageNode.id = id
	storageNode.storageNodeHost = host
	storageNode.storageNodePort = port
	storageNode.config = config
	storageNode.savePath = savePath

	// input size in MB convert to bytes and stored in Math/big.Int
	storageNode.size.SetInt64(size)
	storageNode.size.Set(storageNode.size.Mul(&storageNode.size, big.NewInt(1000)))
	log.Println("storage node size", storageNode.size)
	storageNode.running = true
	return storageNode
}

func (storageNode *StorageNode) Start() {
	connectionHandler, err := connection.NewClient(storageNode.config.ControllerHost, storageNode.config.ControllerPort)
	if err != nil {
		log.Printf("Error while startig the storageNode %s", err)
		return
	}
	storageNode.connectionHandler = connectionHandler

	storageNode.register()

	storageNode.running = true
	go storageNode.heartbeat()
	storageNode.server = connection.NewServer(storageNode.storageNodeHost, int(storageNode.storageNodePort))
	go storageNode.listen()
}

func (storageNode *StorageNode) Shutdown() {
	storageNode.running = false
}

func (storageNode *StorageNode) register() {
	message := &connection.FileData{}
	message.MessageType = connection.MessageType_REGISTRATION
	message.SenderId = storageNode.id
	message.Size = storageNode.size.Bytes()
	node := &connection.Node{}
	node.Port = storageNode.storageNodePort
	node.Hostname = storageNode.storageNodeHost
	message.Node = node
	err := storageNode.connectionHandler.Send(message)
	if err != nil {
		log.Println("Error sending registration to controller ", err)
	}
	ack, err := storageNode.connectionHandler.Receive()
	if err != nil || ack.MessageType != connection.MessageType_ACK {
		log.Println("Error getting ack")
	}
	log.Println("Received ack from controller for registration")
}

func (storageNode *StorageNode) heartbeat() {
	var heartBeatRate = time.Second * 5
	message := &connection.FileData{}
	message.MessageType = connection.MessageType_HEARTBEAT
	message.SenderId = storageNode.id
	message.Size = storageNode.size.Bytes()
	for storageNode.running {
		storageNode.connectionHandler.Send(message)
		time.Sleep(heartBeatRate)
	}
}

func (storageNode *StorageNode) listen() {
	for storageNode.running {
		connectionHandler, err := storageNode.server.NextConnectionHandler()
		if err != nil {
			log.Println("Cannot get the connection handler ", err)
		}
		go storageNode.handleConnection(connectionHandler)
	}
}

func (storageNode *StorageNode) handleConnection(connectionHandler *connection.ConnectionHandler) {
	getChan := make(chan *connection.FileData)
	for storageNode.running {
		message, err := connectionHandler.Receive()
		if err != nil {
			log.Println("Cannot receive message ", err)
		}
		if message.MessageType == connection.MessageType_PUT {
			storageNode.uploadHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_GET {
			go storageNode.downloadHandler(connectionHandler, message, getChan)
		} else if message.MessageType == connection.MessageType_ACK_GET {
			getChan <- message
		}
	}
}

func (storageNode *StorageNode) uploadHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	size := message.DataSize
	path := message.Path
	//nodes := message.Nodes
	//extract file path to be used for heartbeat

	//prepare ack
	response := &connection.FileData{}
	response.MessageType = connection.MessageType_ACK_PUT
	connectionHandler.Send(response)

	//read from stream
	data := make([]byte, size)
	connectionHandler.ReadN(data)

	//compare checksums
	sum := md5.Sum(data)
	if !bytes.Equal(sum[:], message.Checksum) {
		log.Printf("the received chunk data is corrupted for chunk %s \n", path)
		//TODO handle the sad path here
	}
	log.Println("storage node read the data and compared checksums")

	//save file
	dirname := storageNode.savePath
	err := os.Mkdir( dirname, 0700)
	file, err := os.Create("./" + dirname+ "/" + path)
	defer file.Close()
	if err != nil {
		log.Println("Error opening file ", path, err)
		return
	}
	file.Write(data)

	log.Printf("chunk %s saved \n", path)

	//send back ack
	response = &connection.FileData{}
	response.MessageType = connection.MessageType_ACK_PUT
	connectionHandler.Send(response)

	// TODO send some heartbeat with info

	//begin replication
	if len(message.Nodes) > 1 {
		storageNode.replicationHandler(data, message)
	}
}

func (storageNode *StorageNode) replicationHandler(chunkData []byte, message *connection.FileData) {
	//peel off current node from the message
	nodes := peelOffNode(message.Nodes, storageNode.id)

	//prepare message
	replicationMessage := &connection.FileData{}
	replicationMessage.Path = message.Path
	replicationMessage.MessageType = connection.MessageType_PUT
	replicationMessage.Nodes = nodes
	replicationMessage.DataSize = message.DataSize
	replicationMessage.Data = message.Data
	replicationMessage.Checksum = message.Checksum

	//pick destination node
	node := nodes[0]

	//dial the destination node
	conHandler, err := connection.NewClient(node.Hostname, int(node.Port))
	if err != nil {
		log.Println("error dialing the next node for replication")
	}

	//send metadata
	err = conHandler.Send(replicationMessage)
	if err != nil {
		fmt.Println("Error sending replication message to storage node")
	}

	// wait for ack
	result, err := conHandler.Receive()
	if err != nil || result.MessageType != connection.MessageType_ACK_PUT {
		log.Fatalf("Error receiving ack data for replication message from storage node %s to storage node %s \n", node.Id, storageNode.id)
	}

	// send chunk
	err = conHandler.WriteN(chunkData)
	if err != nil {
		fmt.Println("Error sending chunk payload for replication to storage node")
	}
	//wait for ack
	result, err = conHandler.Receive()
	if err != nil || result.MessageType != connection.MessageType_ACK_PUT {
		log.Fatalf("Error receiving ack data for replication message from storage node %s to storage node %s \n", node.Id, storageNode.id)
	}

	log.Println("successfully sent replication message and data ")
}

func peelOffNode(nodes []*connection.Node, id string) []*connection.Node{
	for i, n := range nodes {
		if n.Id == id {
			nodes = append(nodes[:i], nodes[i+1:]...)
			break
		}
	}
	return nodes
}

func (storageNode *StorageNode) downloadHandler(handler *connection.ConnectionHandler,
	message *connection.FileData, getChan <-chan *connection.FileData) {
	sendMessage := &connection.FileData{}
	log.Printf("Storage node %s received request to download %s", storageNode.id, message.Path)
	data, err := file_io.ReadFile(message.Path)
	if err != nil {
		sendMessage.MessageType = connection.MessageType_ERROR
		err := handler.Send(sendMessage)
		if err != nil {
			log.Printf("Error sending erorr message for download chunk on storage node %s", storageNode.id)
		}
	}
	// TODO change to message type of send data
	sendMessage.MessageType = connection.MessageType_GET
	sendMessage.DataSize = int64(len(data))
	err = handler.Send(sendMessage)
	if err != nil {
		log.Printf("Error sending initial data size to client on storage node %s", storageNode.id)
	}
	ack := <- getChan
	if ack.MessageType != connection.MessageType_ACK_GET {
		log.Printf("Error ack for sending initial data size to client on storage node %s", storageNode.id)
	}
	log.Println("Storage node received ack")
	err = handler.WriteN(data)
	if err != nil {
		log.Printf("Error writeN sending data to client download on storage node %s", storageNode.id)
	}
	lastAck := <- getChan
	if lastAck.MessageType != connection.MessageType_ACK_GET {
		log.Printf("Error ack for downlaod data client on storage node %s", storageNode.id)
	}
}
