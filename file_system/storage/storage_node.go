package storage

import (
	"bufio"
	"dfs/config"
	"dfs/connection"
	file_io "dfs/files_io"
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

}

func NewStorageNode(id string, size int64, host string, port int32, config *config.Config) *StorageNode {
	storageNode := &StorageNode{}
	storageNode.id = id
	storageNode.storageNodeHost = host
	storageNode.storageNodePort = port
	storageNode.config = config

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

	//prepare ack
	response := &connection.FileData{}
	response.MessageType = connection.MessageType_ACK_PUT
	connectionHandler.Send(response)

	//read from stream
	data := make([]byte, size)
	connectionHandler.ReadN(data)

	//save file
	file, err := os.Create("./" + path)
	defer file.Close()
	if err != nil {
		log.Println("Error opening file ", path)
		return
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		log.Println("Error writing to file ", path,  err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Println("Error flushing writer to file ", path)
		return
	}

	//send back ack
	response = &connection.FileData{}
	response.MessageType = connection.MessageType_ACK_PUT
	connectionHandler.Send(response)

	// TODO send some heartbeat with info

	//begin replication

}

func (storageNode *StorageNode) replicationHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	//chunkMetaMap[chunkName].Nodes = chunkMetaMap[chunkName].Nodes[1:] //peel of the used node
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
