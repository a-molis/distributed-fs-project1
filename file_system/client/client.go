package client

import (
	. "dfs/config"
	"dfs/connection"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Client struct {
	config            *Config
	command           string
	args              []string
	remotePath        string
	localPath         string
	controllerHost    string
	controllerPort    int
	connectionHandler *connection.ConnectionHandler
}

type chunkMeta struct {
	name string
	size int64
	Nodes []*connection.Node
	data []byte
}

type Node struct {
	Id       string
	Hostname string
	Port     int32
}

func NewClient(config *Config, command string, args ...string) *Client {
	client := &Client{}
	client.config = config
	client.command = command
	client.args = args
	client.controllerHost = config.ControllerHost
	client.controllerPort = config.ControllerPort
	return client
}

var (
	commandMap = map[string]connection.MessageType{
		"ls":  connection.MessageType_LS,
		"rm":  connection.MessageType_RM,
		"put": connection.MessageType_PUT,
		"get": connection.MessageType_GET,
	}
)

func (client *Client) Start() {
	messageType, ok := commandMap[strings.ToLower(client.command)]
	if !ok {
		log.Fatalln("Command not recognized ", client.command)
	}
	message := &connection.FileData{}
	message.MessageType = messageType

	if messageType == connection.MessageType_RM ||
		messageType == connection.MessageType_GET ||
		messageType == connection.MessageType_PUT {
		if len(client.args) < 1 {
			log.Fatalln("Missing arguments")
		}
		client.remotePath = client.args[0]
		message.Path = client.remotePath
		if messageType == connection.MessageType_GET || messageType == connection.MessageType_PUT {
			if len(client.args) < 2 {
				log.Fatalln("Missing local path for get or put")
			}
			client.localPath = client.args[1]
		}
	}
	connectionHandler, err := connection.NewClient(client.controllerHost, client.controllerPort)
	if err != nil {
		log.Fatalf("Error client unable to connect to Controller %s", err)
		return
	}
	log.Printf("Client connected to controller")
	client.connectionHandler = connectionHandler

	// send command
	client.sendToController(err, message, connectionHandler)
}

func (client *Client) sendToController(err error, message *connection.FileData, connectionHandler *connection.ConnectionHandler) {
	// TODO refactor to immediately enter method for specific messageType
	err = client.connectionHandler.Send(message)
	if err != nil {
		log.Fatalln("Error sending data to controller")
	}
	log.Printf("Client sent command to server")
	result, err := connectionHandler.Receive()
	if err != nil {
		log.Fatalln("Error receiving data from controller on the client")
	}
	log.Printf("Client received message back from controller")

	if result.MessageType == connection.MessageType_LS {
		log.Printf("Client received ls message back from controller")
		client.ls(result, connectionHandler)
	} else if result.MessageType == connection.MessageType_PUT {
		log.Println("Client received put message")
		client.put(result, connectionHandler)
	} else if result.MessageType == connection.MessageType_ERROR {
		log.Fatalln("Error from controller: ", result.Data)
	} else {
		log.Fatalln("Error client unable to get result from controller")
	}
}

func (client *Client) ls(result *connection.FileData, connectionHandler *connection.ConnectionHandler) {
	fmt.Println(result.Data)
	ackLS := &connection.FileData{}
	ackLS.MessageType = connection.MessageType_ACK_LS
	err := connectionHandler.Send(ackLS)
	if err != nil {
		log.Println("Error sending ack ls to controller")
	}
}

func (client *Client) put(result *connection.FileData, connectionHandler *connection.ConnectionHandler) {
	chunkMetaMap, err := client.getChunkMeta()
	if err != nil {
		log.Fatalln("Error getting chunk data from file on client ", err)
	}
	// The method below fills chunkMetaMap with the correct storage node info
	err = client.sendChunkInfoController(connectionHandler, chunkMetaMap)
	if err != nil {
		log.Fatalln("Error sending chunk info from client to server ", err)
	}

	//TODO call function that sends chunks to the storage nodes defined in chunkMetaMap
	// map with storage node and connection handler
	// generate the chunks-data as byte arrays
	// send the chunk data to corresponding node
	//client.sendChunksToNodes(chunkMetaMap)
}

func (client *Client) getChunkMeta() (map[string]*chunkMeta, error) {
	fileSize := getFileSize(client.localPath)
	numChunks := int(fileSize / client.config.ChunkSize)
	chunkMap := make(map[string]*chunkMeta)
	filePrefix := strings.Replace(strings.TrimPrefix(client.remotePath, "/"), "/", "_", -1)
	for i := 0; i < numChunks; i++ {
		chunkName := fmt.Sprintf("%s_%d", filePrefix, i)
		chunkMeta := &chunkMeta{}
		chunkMeta.size = client.config.ChunkSize
		chunkMeta.name = chunkName
		chunkMap[chunkName] = chunkMeta
	}
	remainingSize := fileSize % client.config.ChunkSize
	if remainingSize != 0 {
		lastChunkName := fmt.Sprintf("%s_%d", filePrefix, numChunks)
		chunkMeta := &chunkMeta{}
		chunkMeta.size = remainingSize
		chunkMeta.name = lastChunkName
		chunkMap[lastChunkName] = chunkMeta
	}
	return chunkMap, nil
}

func (client *Client) sendChunkInfoController(handler *connection.ConnectionHandler, chunkMap map[string]*chunkMeta) error {
	fileData := &connection.FileData{}
	fileData.MessageType = connection.MessageType_ACK_PUT
	protoChunks := make([]*connection.Chunk, len(chunkMap))
	index := 0
	for _, chunk := range chunkMap {
		protoChunk := clientChunkToProto(chunk)
		protoChunks[index] = protoChunk
		index++
	}
	fileData.Chunk = protoChunks
	err := handler.Send(fileData)
	if err != nil {
		return err
	}
	message, err := handler.Receive() //TODO check for errors
	if err != nil {
		return err
	}
	for _, protoChunk := range message.Chunk {
		initialChunk, ok := chunkMap[protoChunk.Name]
		if !ok {
			log.Fatalln("Error getting chunk data from controller on client")
		}
		initialChunk.Nodes = protoChunk.Nodes
	}
	return nil
}

type BlockingConnection struct {
	conHandler *connection.ConnectionHandler
	mu *sync.Mutex
}

func (client *Client) sendChunksToNodes(chunkMetaMap map[string]*chunkMeta) {

	// get chunk byte array
	err := client.getChunkData(chunkMetaMap)
	if err != nil {
		log.Println("error getting chunk data")
		return
	}

	blockingHandlerMap := make(map[string]*BlockingConnection)
	//handlerMap := make(map[string]*connection.ConnectionHandler)
	// add handler mutex map
	// you need to aqcuire the lock before using the connecction handleer so that we dont accidentally send two chunks at the same time to same node
	//handlerMutexMap := make(map[string]*sync.Mutex)

	// file := os.Open(filename)
	for chunkName := range chunkMetaMap {
		//TODO do this bit as a go routine
		// bytes = readNbytes()
		//

		// get connection handler corresponding to storage node //this is gona block if con hand in use
		// node := chunkMetaMap[chunkName].Nodes[0]

		//mu := handlerMutexMap[node.Id] //this bit is to ensure that if multiple chunks are sent to the same node, it is in order
		//mu.Lock()

		blockingHandler := getConnectionHandler(chunkMetaMap, chunkName, blockingHandlerMap)

		blockingHandler.mu.Lock()
		conHandler := blockingHandler.conHandler
		// send metadata

		// create messaeg
		message := &connection.FileData{}
		message.MessageType = connection.MessageType_PUT
		message.Path = chunkName
		message.DataSize = chunkMetaMap[chunkName].size
		message.Nodes = chunkMetaMap[chunkName].Nodes
		//message.checksum TODO

		err := conHandler.Send(message)
		if err != nil {
			fmt.Println("Error sending chunk metadata to storage node")
		}
		// wait for ack
		result, err := conHandler.Receive()
		if err != nil || result.MessageType != connection.MessageType_ACK_PUT {
			log.Fatalln("Error receiving ack data for put metadata from storage node on the client")
		}
		// send chunk
		err = conHandler.WriteN(chunkMetaMap[chunkName].data)
		if err != nil {
			fmt.Println("Error sending chunk payload to storage node")
		}
		//wait for ack
		result, err = conHandler.Receive()
		if err != nil || result.MessageType != connection.MessageType_ACK_PUT {
			log.Fatalln("Error receiving ack data for put payload from storage node on the client")
		}
		blockingHandler.mu.Unlock()
	}
}

func getConnectionHandler(chunkMetaMap map[string]*chunkMeta, chunkName string, handlerMap map[string]*BlockingConnection) *BlockingConnection {
	node := chunkMetaMap[chunkName].Nodes[0]

	blockingConnection, ok := handlerMap[node.Id]
	if !ok {
		blockingConnection = &BlockingConnection{}
		conHandler, err := connection.NewClient(node.Hostname, int(node.Port))
		blockingConnection.conHandler = conHandler
		if err != nil {
			log.Println("error connecting to node ", node.Hostname, " ", node.Port)
			// TODO make this trigger switching to next node and repeat this iteration of the loop
		}
		handlerMap[node.Id] = blockingConnection
	}

	return blockingConnection
}

func (client *Client) getChunkData(chunkMetaMap map[string]*chunkMeta) error {
	file, err := os.Open(client.localPath)
	if err != nil {
		log.Println("Cannot open target file")
		return err
	}

	for chunkName := range chunkMetaMap {
		chunkData := make([]byte, chunkMetaMap[chunkName].size)
		numBytes, err := file.Read(chunkData)
		if err != nil {
			log.Println("Error reading bytes from target file", err)
			return err
		}
		if int64(numBytes) != chunkMetaMap[chunkName].size {
			log.Println("Did not read same bytes as chunk size", chunkName)
		}
		chunkMetaMap[chunkName].data = chunkData
	}
	return nil
}

func clientChunkToProto(chunkMeta *chunkMeta) *connection.Chunk {
	protoChunk := &connection.Chunk{}
	protoChunk.Name = chunkMeta.name
	protoChunk.Size = chunkMeta.size
	return protoChunk
}

func getFileSize(path string) int64 {
	fileStat, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Error getting file stat for %s\n%s\n", path, err)
	}
	return fileStat.Size()
}
