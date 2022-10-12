package controller

import (
	"bufio"
	. "dfs/config"
	"dfs/connection"
	"dfs/file_metadata"
	fileio "dfs/files_io"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

type Controller struct {
	id           string
	host         string
	port         int
	server       *connection.Server
	memberTable  *MemberTable
	fileMetadata *file_metadata.FileMetadata
	running      bool
	config       *Config
}

func NewController(id string, config *Config) *Controller {
	controller := &Controller{}
	controller.id = id
	controller.port = config.ControllerPort
	controller.memberTable = NewMemberTable()
	controller.host = config.ControllerHost
	controller.config = config
	controller.fileMetadata = file_metadata.NewFileMetaData()
	loadErr := controller.LoadFileMetadata()
	if loadErr != nil {
		log.Println("Unable to load file metadata, stopping new controller")
	}
	return controller
}

func (controller *Controller) Start() {
	controller.running = true
	log.Println("Starting controller")
	server, err := connection.NewServer(controller.host, controller.port)
	if err != nil {
		log.Fatalln("Unable to create new server on controller")
	}
	controller.server = server
	controller.listen()
	log.Println("Controller shut down")
}

func (controller *Controller) listen() {
	log.Println("Starting controller listener")
	for controller.running {
		connectionHandler, err := controller.server.NextConnectionHandler()
		if err != nil {
			log.Println("Cannot get the connection handler ", err)
		}
		go controller.handleConnection(connectionHandler)
	}
	log.Println("Shutting down controller listen")
}

func (controller *Controller) handleConnection(connectionHandler *connection.ConnectionHandler) {
	lsChan := make(chan *connection.FileData)
	putChan := make(chan *connection.FileData)
	for controller.running {
		message, err := connectionHandler.Receive()
		if err != nil {
			log.Println("Cannot receive message ", err)
		}
		if message.MessageType == connection.MessageType_REGISTRATION {
			go controller.registerHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_HEARTBEAT ||
			message.MessageType == connection.MessageType_HEARTBEAT_CHUNK {
			go controller.heartbeatHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_LS {
			go controller.Ls(connectionHandler, lsChan, message)
		} else if message.MessageType == connection.MessageType_ACK_LS {
			lsChan <- message
		} else if message.MessageType == connection.MessageType_PUT {
			go controller.uploadHandler(connectionHandler, putChan, message)
		} else if message.MessageType == connection.MessageType_ACK_PUT {
			putChan <- message
		} else if message.MessageType == connection.MessageType_GET {
			go controller.getHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_DUMMY {
			// nothing to see here
			// TODO fix this nothing
		} else if message.MessageType == connection.MessageType_RM {
			go controller.rmHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_CHECKSUM {
			go controller.updateCheckSum(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_STATS {
			go controller.statsHandler(connectionHandler, message)
		}
	}
}

func (controller *Controller) updateCheckSum(handler *connection.ConnectionHandler, message *connection.FileData) {
	checkSum := message.Checksum
	path := message.Path
	err := controller.fileMetadata.UpdateChecksum(path, checkSum)
	ack := &connection.FileData{}
	if err != nil {
		log.Printf("Error saving checksum for file %s %s\n", path, err)
		ack.MessageType = connection.MessageType_ERROR
	} else {
		ack.MessageType = connection.MessageType_ACK
		log.Println("updated Checksum for file")
	}
	err = handler.Send(ack)
	if err != nil {
		log.Printf("Error sending message for check sum ack for %s %s", path, err)
	}
}

func (controller *Controller) shutdown() {
	controller.running = false
	controller.memberTable.Shutdown()
	log.Println("Shutting down controller listening server")
	err := controller.server.Shutdown()
	if err != nil {
		log.Println("Error shutting down controller listening server")
	}
}

func (controller *Controller) registerHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	log.Println("Received registration message from ", message.SenderId)
	size := new(big.Int).SetBytes(message.Size)
	err := controller.memberTable.Register(message.SenderId, size, message.Node.Hostname, message.Node.Port, connectionHandler)
	if err == nil {
		ack := &connection.FileData{}
		ack.MessageType = connection.MessageType_ACK
		connectionHandler.Send(ack)
	} else {
		log.Println("Error handling registration request")
	}
}

func (controller *Controller) heartbeatHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	log.Printf("Received heart beat from %s on %s\n", message.SenderId, controller.id)
	size := new(big.Int).SetBytes(message.Size)
	controller.memberTable.RecordBeat(message.SenderId, size, message.NumberOfRequests)

	//TODO update file metadata with info that is passed in heartbeat
	//^^ returns a boolean
	//depending on boolean save file metadata
	// we still update the chunk info so it might be best to save it regardless
	if message.MessageType == connection.MessageType_HEARTBEAT_CHUNK {
		// call some function here
		controller.fileMetadata.HeartbeatHandler(message.Path, message.Data, message.Checksum)
		controller.SaveFileMetadata()
	}
}

func (controller *Controller) List() []string {
	return controller.memberTable.List()
}

func (controller *Controller) uploadHandler(connectionHandler *connection.ConnectionHandler, connectionChan <-chan *connection.FileData, message *connection.FileData) {
	filepath := message.GetPath()
	if filepath == "" {
		// TODO send error back to client
	}
	exists := controller.fileMetadata.PathExists(filepath)

	if exists {
		err := fileio.SendError(connectionHandler, "File already exists")
		if err != nil {
			log.Println("Unable to send error to client")
		}
		return
	}
	err := fileio.SendMessage(connectionHandler, connection.MessageType_PUT)
	chunkData := <-connectionChan
	chunks := ProtoToChunk(chunkData.GetChunk())

	err = findAvailableNodes(chunks, controller.memberTable, controller.config.NumReplicas)
	if err != nil {
		log.Println("could  not assign nodes to chunks, ", err)
	}
	err = controller.fileMetadata.UploadChunks(filepath, chunks)
	if err != nil {
		log.Println("could  not upload chunks to filetree, ", err)
	}
	// prepare client response

	response := &connection.FileData{}
	response.MessageType = connection.MessageType_PUT
	response.Chunk = chunkToProto(chunks, controller.memberTable)
	response.Path = filepath
	connectionHandler.Send(response)
}

func (controller *Controller) Ls(handler *connection.ConnectionHandler, connectionChan <-chan *connection.FileData, message *connection.FileData) {
	log.Printf("LS for path %s\n", message.Path)
	sendMessage := &connection.FileData{}
	sendMessage.MessageType = connection.MessageType_LS
	data, err := controller.fileMetadata.Ls(message.Path)

	// TODO update to have ls logic to get files in directory path
	sendMessage.Data = data
	err = handler.Send(sendMessage)
	if err != nil {
		log.Println("Error sending ls data to client ", err)
	}
	ack := <-connectionChan
	if ack.MessageType != connection.MessageType_ACK_LS {
		// TODO retry if need to
		log.Println("Invalid ack for controller ls")
	}
}

func (controller *Controller) SaveFileMetadata() {
	bytes, err := controller.fileMetadata.GetBytes()
	if err != nil {
		log.Println("Error converting filemetadata to bytes")
		return
	}
	file, err := os.Create(controller.config.ControllerPath)
	defer file.Close()
	if err != nil {
		log.Println("Error opening file fmdt")
		return
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		log.Println("Error writing to file fmdt", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Println("Error flushing writer to file fmdt")
		return
	}
}

func (controller *Controller) LoadFileMetadata() error {
	_, err := os.Stat(controller.config.ControllerPath)
	if err != nil {
		log.Println("Error finding file for meta data load")
		return nil
	}

	file, err := os.Open(controller.config.ControllerPath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file for meta data load")
		}
	}(file)
	//if os.IsNotExist(err) {
	//	log.Println("No file metadata backup found on disk skipping...")
	//	return nil
	//}
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Error reading in metadata file")
		return err
	}
	err = controller.fileMetadata.LoadBytes(bytes)
	if err != nil {
		return err
	}
	return err
}

func (controller *Controller) getHandler(handler *connection.ConnectionHandler, message *connection.FileData) {
	sendMessage := &connection.FileData{}

	// TODO add check if in pending state //done
	chunks, checksum, err := controller.fileMetadata.Download(message.Path)
	controller.RemoveDeadNodes(chunks)
	if err != nil {
		sendMessage.Data = fmt.Sprintf("%s", err)
		sendMessage.MessageType = connection.MessageType_ERROR
		err := handler.Send(sendMessage)
		if err != nil {
			log.Println("Error sending error for get operation")
		}
		return
	}
	sendMessage.Checksum = checksum
	sendMessage.Chunk = chunkToProto(chunks, controller.memberTable)
	sendMessage.MessageType = connection.MessageType_GET
	err = handler.Send(sendMessage)
}

func (controller *Controller) RemoveDeadNodes(chunks map[string]*file_metadata.Chunk) {
	for _, v := range chunks {
		for i, n := range v.StorageNodes {
			if !controller.memberTable.members[n].status {
				// if the member is inactive remove it from the list
				if i >= len(v.StorageNodes)-1 {
					v.StorageNodes = v.StorageNodes[:len(v.StorageNodes)-1]
				} else {
					v.StorageNodes = append(v.StorageNodes[:i], v.StorageNodes[i+1:]...)
				}
			}
		}
	}
}

func (controller *Controller) rmHandler(handler *connection.ConnectionHandler, message *connection.FileData) {
	deleteMessage := &connection.FileData{}
	if message.Path == "" {
		deleteMessage.MessageType = connection.MessageType_ERROR
		err := handler.Send(deleteMessage)
		if err != nil {
			log.Println("Error sending error for missing path")
		}
		return
	}
	log.Printf("Removing path %s\n", message.Path)

	chunks, err := controller.fileMetadata.Remove(message.Path)
	deleteMessage.Chunk = chunkToProto(chunks, controller.memberTable)
	if err != nil {
		deleteMessage.Data = fmt.Sprintf("%s", err)
		deleteMessage.MessageType = connection.MessageType_ERROR
	} else {
		deleteMessage.Data = fmt.Sprintf("Removed %s from file index ", message.Path)
		deleteMessage.MessageType = connection.MessageType_RM
	}
	log.Println("Sending rm message with chunk data back to client")
	err = handler.Send(deleteMessage)
	if err != nil {
		log.Println("Error sending delete chunks to client")
	}
	log.Println("Controller finished removing file ", message.Path)
}

func (controller *Controller) statsHandler(handler *connection.ConnectionHandler, message *connection.FileData) {
	stats, err := controller.memberTable.Stats()
	statMessage := &connection.FileData{}
	if err != nil {
		statMessage.MessageType = connection.MessageType_ERROR
		statMessage.Data = "Server error getting stats"
	} else {
		statMessage.Data = stats
		statMessage.MessageType = connection.MessageType_STATS
	}

	err = handler.Send(statMessage)
	if err != nil {
		log.Println("Error sending stats to client from controller ", controller.id)
	}
}
