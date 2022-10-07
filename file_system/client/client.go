package client

import (
	"crypto/md5"
	"crypto/sha256"
	. "dfs/atomic_data"
	. "dfs/config"
	"dfs/connection"
	"errors"
	"fmt"
	"hash"
	"log"
	"os"
	"reflect"
	"sort"
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
	name  string
	size  int64
	Nodes []*connection.Node
	num   int32
}

type Node struct {
	Id       string
	Hostname string
	Port     int32
}

type BlockingConnection struct {
	conHandler *connection.ConnectionHandler
	mu         sync.Mutex
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
		"ls":    connection.MessageType_LS,
		"rm":    connection.MessageType_RM,
		"put":   connection.MessageType_PUT,
		"get":   connection.MessageType_GET,
		"stats": connection.MessageType_STATS,
	}
)

func (client *Client) Run() (*string, error) {
	messageType, ok := commandMap[strings.ToLower(client.command)]
	if !ok {
		log.Println("Command not recognized ", client.command)
		return nil, errors.New("command not recognized")
	}
	message := &connection.FileData{}
	message.MessageType = messageType

	if messageType == connection.MessageType_RM ||
		messageType == connection.MessageType_GET ||
		messageType == connection.MessageType_PUT ||
		messageType == connection.MessageType_LS {
		if len(client.args) < 1 {
			log.Println("Missing arguments")
			return nil, errors.New("missing arguments")
		}
		client.remotePath = client.args[0]
		message.Path = client.remotePath
		if messageType == connection.MessageType_GET || messageType == connection.MessageType_PUT {
			if len(client.args) < 2 {
				log.Println("Missing local path for get or put")
				return nil, errors.New("missing local path for get or put")
			}
			client.localPath = client.args[1]
		}
	}
	connectionHandler, err := connection.NewClient(client.controllerHost, client.controllerPort)
	if err != nil {
		log.Printf("Error client unable to connect to Controller %s\n", err)
		return nil, errors.New(fmt.Sprintf("Error client unable to connect to Controller %s", err))
	}
	log.Printf("Client connected to controller")
	client.connectionHandler = connectionHandler

	// send command
	return client.sendToController(err, message, connectionHandler)
}

func (client *Client) sendToController(err error, message *connection.FileData, connectionHandler *connection.ConnectionHandler) (*string, error) {
	// TODO refactor to immediately enter method for specific messageType
	err = client.connectionHandler.Send(message)
	if err != nil {
		log.Println("Error sending data to controller")
		return nil, errors.New("error sending data to controller")
	}
	log.Printf("Client sent command to server")
	result, err := connectionHandler.Receive()
	if err != nil {
		log.Println("Error receiving data from controller on the client")
		return nil, errors.New("error receiving data from controller on the client")
	}
	if result.MessageType == connection.MessageType_LS {
		log.Printf("Client received ls message back from controller")
		return client.ls(result, connectionHandler)
	} else if result.MessageType == connection.MessageType_PUT {
		log.Println("Client received put message")
		return nil, client.put(result, connectionHandler)
	} else if result.MessageType == connection.MessageType_GET {
		return nil, client.get(result, connectionHandler)
	} else if result.MessageType == connection.MessageType_RM {
		return nil, client.rm(result, connectionHandler)
	} else if result.MessageType == connection.MessageType_ERROR {
		log.Println("Error: ", result.Data)
		return nil, errors.New(fmt.Sprintf("Error: %s", result.Data))
	} else if result.MessageType == connection.MessageType_STATS {
		return client.stats(result, connectionHandler)
	} else {
		log.Println("Error client unable to get result from controller")
		return nil, errors.New("error client unable to get result from controller")
	}
}

func (client *Client) ls(result *connection.FileData, connectionHandler *connection.ConnectionHandler) (*string, error) {
	fmt.Println(result.Data)
	ackLS := &connection.FileData{}
	ackLS.MessageType = connection.MessageType_ACK_LS
	err := connectionHandler.Send(ackLS)
	if err != nil {
		log.Println("Error sending ack ls to controller")
		return nil, errors.New("error client unable to send ls ack")
	}
	return &result.Data, nil
}

func (client *Client) put(result *connection.FileData, connectionHandler *connection.ConnectionHandler) error {
	chunkMetaMap, err := client.getChunkMeta()
	if err != nil {
		log.Println("Error getting chunk data from file on client ", err)
		return errors.New(fmt.Sprintf("error getting chunk data from file on client %s\n", err))
	}
	// The method below fills chunkMetaMap with the correct storage node info
	err = client.sendChunkInfoController(connectionHandler, chunkMetaMap)
	if err != nil {
		log.Println("Error sending chunk info from client to server ", err)
		return errors.New(fmt.Sprintf("Error sending chunk info from client to server %s", err))
	}

	//TODO call function that sends chunks to the storage nodes defined in chunkMetaMap
	// map with storage node and connection handler
	// generate the chunks-data as byte arrays
	// send the chunk data to corresponding node
	checkSum, err := client.sendChunksToNodes(chunkMetaMap)
	if err != nil {
		return err
	}
	err = client.sendCheckSumToController(checkSum, connectionHandler)
	if err != nil {
		log.Println("Error sending checksum to controller ", err)
		return errors.New(fmt.Sprintf("Error sending checksum to controller %s", err))
	}
	return nil
}

func (client *Client) sendCheckSumToController(checkSum []byte, handler *connection.ConnectionHandler) error {
	message := &connection.FileData{}
	message.MessageType = connection.MessageType_CHECKSUM
	message.Checksum = checkSum
	message.Path = client.remotePath
	err := handler.Send(message)
	if err != nil {
		return err
	}
	ack, err := handler.Receive()
	if err != nil || ack.MessageType != connection.MessageType_ACK {
		log.Println("Error receiving ack from controller for file checksum ", err)
		return errors.New("error sending checksum to controller")
	}
	return nil
}

func (client *Client) getChunkMeta() (map[string]*chunkMeta, error) {
	fileSize := getFileSize(client.localPath)
	numChunks := int(fileSize / client.config.ChunkSize)
	chunkMap := make(map[string]*chunkMeta)
	filePrefix := strings.Replace(strings.TrimPrefix(client.remotePath, "/"), "/", "_", -1)
	filePrefix = strings.ReplaceAll(filePrefix, ".", "")
	for i := 0; i < numChunks; i++ {
		chunkName := fmt.Sprintf("%s_%d", filePrefix, i)
		chunkMeta := &chunkMeta{}
		chunkMeta.size = client.config.ChunkSize
		chunkMeta.name = chunkName
		chunkMeta.num = int32(i)
		chunkMap[chunkName] = chunkMeta
	}
	remainingSize := fileSize % client.config.ChunkSize
	if remainingSize != 0 {
		lastChunkName := fmt.Sprintf("%s_%d", filePrefix, numChunks)
		chunkMeta := &chunkMeta{}
		chunkMeta.size = remainingSize
		chunkMeta.name = lastChunkName
		chunkMeta.num = int32(numChunks)
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
			log.Println("Error getting chunk data from controller on client")
			return errors.New("error getting chunk data from controller on client")
		}
		initialChunk.Nodes = protoChunk.Nodes
	}
	return nil
}

func (client *Client) sendChunksToNodes(chunkMetaMap map[string]*chunkMeta) ([]byte, error) {

	blockingHandlerMap := make(map[string]*BlockingConnection)
	errorMap := make(map[string]error)
	errorLock := &sync.Mutex{}

	file, err := os.Open(client.localPath)
	if err != nil {
		log.Println("Cannot open target file")
		return nil, err
	}
	checksum := sha256.New()
	for chunkName := range chunkMetaMap {

		chunkData, err := getChunkData(chunkMetaMap, chunkName, file)
		if err != nil {
			log.Println("Error getting chunk data for chunk ", chunkName)
			return nil, err
		}
		checksum.Write(chunkData)
		//TODO do this bit as a go routine
		go client.sendChunk(chunkMetaMap, chunkName, blockingHandlerMap, chunkData, errorMap, errorLock)
	}
	for id, e := range errorMap {
		if e != nil {
			log.Printf("Error recorded for node %s, error: %s  \n", id, e)
			return nil, e
		}
	}
	return checksum.Sum(nil), nil
}

func (client *Client) sendChunk(chunkMetaMap map[string]*chunkMeta, chunkName string, blockingHandlerMap map[string]*BlockingConnection, chunkData []byte, errorMap map[string]error, mutex *sync.Mutex) {
	// create message
	message := &connection.FileData{}
	message.MessageType = connection.MessageType_PUT
	message.Path = chunkName
	message.DataSize = chunkMetaMap[chunkName].size
	message.Nodes = chunkMetaMap[chunkName].Nodes
	message.Data = client.remotePath
	sum := md5.Sum(chunkData)
	message.Checksum = sum[:] //create a slice to convert from [16]byte to []byte
	var node *connection.Node
	var err error
	for _, node = range chunkMetaMap[chunkName].Nodes {
		blockingHandler, err := getConnectionHandler(node, blockingHandlerMap)
		if err != nil {
			log.Printf("had issue dialing node %s \n", node.Id)
			continue
		}
		//send metadata and payload
		blockingHandler.mu.Lock()
		err = client.tryNode(blockingHandler.conHandler, message, chunkData)
		blockingHandler.mu.Unlock()
		if err == nil {
			break
		}
		log.Printf("had issue sending chunk %s to node %s \n", chunkName, node.Id)
	}
	if err != nil {
		mutex.Lock()
		errorMap[node.Id] = err
		mutex.Unlock()
	}
	log.Printf("sent chunk %s info to storage node %s \n", chunkName, node.Id)
}

func (client *Client) tryNode(conHandler *connection.ConnectionHandler, message *connection.FileData, chunkData []byte) error {
	err := conHandler.Send(message)
	if err != nil {
		fmt.Println("Error sending chunk metadata to storage node")
		return err
	}
	log.Println("Client sent chunk metadata to storage node")
	// wait for ack
	result, err := conHandler.Receive()
	if err != nil || result.MessageType != connection.MessageType_ACK_PUT {
		log.Println("Error receiving ack data for put metadata from storage node on the client")
		return err
	}
	// send chunk
	err = conHandler.WriteN(chunkData)
	if err != nil {
		fmt.Println("Error sending chunk payload to storage node")
		return err
	}
	log.Println("Client wrote the chunk data to node stream")
	//wait for ack
	result, err = conHandler.Receive()
	if err != nil || result.MessageType != connection.MessageType_ACK_PUT {
		log.Println("Error receiving ack data for put payload from storage node on the client")
		return err
	}
	return nil
}

func getChunkData(chunkMetaMap map[string]*chunkMeta, chunkName string, file *os.File) ([]byte, error) {
	chunkData := make([]byte, chunkMetaMap[chunkName].size)
	numBytes, err := file.Read(chunkData)
	if err != nil {
		log.Println("Error reading bytes from target file", err)
		return nil, err
	}
	if int64(numBytes) != chunkMetaMap[chunkName].size {
		log.Println("Did not read same bytes as chunk size", chunkName)
	}
	return chunkData, err
}

func getConnectionHandler(node *connection.Node, handlerMap map[string]*BlockingConnection) (*BlockingConnection, error) {
	blockingConnection, ok := handlerMap[node.Id]
	if !ok {
		blockingConnection = &BlockingConnection{}
		conHandler, err := connection.NewClient(node.Hostname, int(node.Port))
		blockingConnection.conHandler = conHandler
		if err != nil {
			log.Println("error connecting to node ", node.Hostname, " ", node.Port)
			// TODO make this trigger switching to next node and repeat this iteration of the loop
			return nil, errors.New("cannot connect to node")
		}
		handlerMap[node.Id] = blockingConnection
	}

	return blockingConnection, nil
}

func (client *Client) get(result *connection.FileData, handler *connection.ConnectionHandler) error {
	savedChecksum := result.Checksum
	log.Println(len(savedChecksum))
	chunks := result.Chunk
	sort.SliceStable(chunks, func(i, j int) bool {
		return chunks[i].Num < chunks[j].Num
	})
	blockingHandlerMap := make(map[string]*BlockingConnection)

	// Using a channel as a blocking queue with size 10
	// bytes only written into queue when in order
	var downloadChan = make(chan []byte, 5)
	numChunks := int32(len(chunks))
	nextChunkNum := NewAtomicInt(&chunks[0].Num)
	saveLock := NewWaitNotify()
	var done *bool
	f := false
	done = &f
	checkSum := sha256.New()
	var saveError error = nil
	go client.saveData(numChunks, downloadChan, saveLock, checkSum, done, &saveError)

	for _, chunk := range chunks {
		go getChunkFromStorage(chunk, blockingHandlerMap, nextChunkNum, saveLock, downloadChan)
	}

	saveLock.Cond.L.Lock()
	for true {
		if !*done {
			saveLock.Cond.Wait()
		} else {
			if reflect.DeepEqual(savedChecksum, checkSum) {
				log.Println("Error saving file, checksums do not match")
				return errors.New("error saving file, checksums do not match")
			} else {
				log.Println("File checksums match")
				return nil
			}
		}
	}
	return nil
}

func getChunkFromStorage(chunk *connection.Chunk, handlerMap map[string]*BlockingConnection,
	nextChunkNum *AtomicInt, saveLock *WaitNotify, downloadChan chan []byte) {
	// TODO should not just get the first node, should check other nodes if this fails

	var data []byte
	var err error
	for _, node := range chunk.Nodes {

		blockingHandler, err2 := getConnectionHandler(node, handlerMap)
		if err2 != nil {
			log.Println("Error dialing node")
			continue
		}
		blockingHandler.mu.Lock()
		conn := blockingHandler.conHandler

		data, err = getBytesFromStorage(conn, chunk)
		blockingHandler.mu.Unlock()
		if err == nil {
			break
		}

		log.Println("Error getting data from storage node")

	}


	saveLock.Cond.L.Lock()
	sent := false
	for !sent {
		if nextChunkNum.Get() == chunk.Num {
			log.Printf("Writing data for chunk num %d\n", chunk.Num)
			downloadChan <- data
			nextChunkNum.Increment()
			saveLock.Cond.L.Unlock()
			break
		}
		saveLock.Cond.Wait()
	}

}

func getBytesFromStorage(conn *connection.ConnectionHandler, chunk *connection.Chunk) ([]byte, error) {
	message := &connection.FileData{}
	message.Path = chunk.Name
	message.MessageType = connection.MessageType_GET
	err := conn.Send(message)
	if err != nil {
		log.Printf("Error getting data from storage node for chunk %s\n", chunk.Name)
		return nil, err
	}
	dataSize, err := conn.Receive()
	if err != nil || dataSize.MessageType == connection.MessageType_ERROR {
		log.Println("Error client getting size from storage node")
		return nil, err
	}
	ack := &connection.FileData{}
	ack.MessageType = connection.MessageType_ACK_GET
	err = conn.Send(ack)
	if err != nil {
		log.Printf("Error sending ack to storage node for size for chunk %sn", chunk.Name)
		return nil, err
	}
	byteData := make([]byte, dataSize.DataSize)
	err = conn.ReadN(byteData)
	if err != nil {
		log.Printf("Error reading byte data for chunk %s\n", chunk.Name)
		return nil, err
	}
	err = conn.Send(ack)
	if err != nil {
		log.Printf("Error sending ack to storage node for receiving raw bytes %s\n", chunk.Name)
		return nil, err
	}
	return byteData, nil
}

func (client *Client) saveData(numChunks int32, downloadChan chan []byte, saveLock *WaitNotify, hash hash.Hash, done *bool, e *error) {
	file, err := os.OpenFile(client.localPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open file on client for saving data ", err)
		*e = err
	}
	for i := numChunks; i > 0; i-- {
		data := <-downloadChan
		write, err := file.Write(data)
		hash.Write(data)
		if err != nil || write != len(data) {
			log.Println("Error writing to save file ", err)
			*e = err
		}
		log.Printf("recevied data of size %d\n", len(data))
		saveLock.Cond.Broadcast()
	}
	*done = true
	saveLock.Cond.Broadcast()
}

func (client *Client) rm(result *connection.FileData, handler *connection.ConnectionHandler) error {
	path := client.remotePath
	chunkData, err := handler.Receive()
	if err != nil {
		log.Printf("Error getting data for delete %s\n", err)
		return err
	}
	if chunkData.MessageType != connection.MessageType_RM {
		log.Printf("Unable to remove data from controller from client %s %s\n", path, err)
		return errors.New(chunkData.Data)
	}
	log.Printf("Number of chunks to remvoe %d\n", len(chunkData.Chunk))
	fmt.Printf("Removed %s\n", path)
	return nil
}

func (client *Client) stats(result *connection.FileData, handler *connection.ConnectionHandler) (*string, error) {
	fmt.Println(result.Data)
	return &result.Data, nil
}

func clientChunkToProto(chunkMeta *chunkMeta) *connection.Chunk {
	protoChunk := &connection.Chunk{}
	protoChunk.Name = chunkMeta.name
	protoChunk.Size = chunkMeta.size
	protoChunk.Num = chunkMeta.num
	return protoChunk
}

func getFileSize(path string) int64 {
	fileStat, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Error getting file stat for %s\n%s\n", path, err)
	}
	return fileStat.Size()
}
