package controller

import (
	"P1-go-distributed-file-system/connection"
	"P1-go-distributed-file-system/file_metadata"
	"log"
)

type Controller struct {
	id          string
	host        string
	port        int
	server      *connection.Server
	memberTable *MemberTable
	fileMetadata *file_metadata.FileMetadata
	running     bool
}

func NewController(id string, host string, port int) *Controller {
	controller := &Controller{}
	controller.id = id
	controller.port = port
	controller.memberTable = NewMemberTable()
	controller.running = true
	controller.host = host
	controller.fileMetadata = file_metadata.NewFileMetaData()
	return controller
}

func (controller *Controller) Start() {
	controller.server = connection.NewServer(controller.host, controller.port)
	go controller.listen()
}

func (controller *Controller) listen() {
	for controller.running {
		connectionHandler, err := controller.server.NextConnectionHandler()
		if err != nil {
			log.Println("Cannot get the connection handler ", err)
		}
		go controller.handleConnection(connectionHandler)
	}
}

func (controller *Controller) handleConnection(connectionHandler *connection.ConnectionHandler) {
	connectionChan := make(chan *connection.FileData)
	for controller.running {
		message, err := connectionHandler.Receive()
		if err != nil {
			log.Println("Cannot receive message ", err)
		}
		if message.MessageType == connection.MessageType_REGISTRATION {
			go controller.registerHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_HEARTBEAT {
			go controller.heartbeatHandler(connectionHandler, message)
		} else if message.MessageType == connection.MessageType_LS {
			go controller.ls(connectionHandler, connectionChan, message)
		} else if message.MessageType == connection.MessageType_ACK_LS {
			connectionChan <- message
		} else if message.MessageType == connection.MessageType_PUT {
			go controller.uploadHandler(connectionHandler, message)
		}
	}
}

func (controller *Controller) shutdown() {
	controller.running = false
}

func (controller *Controller) registerHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	log.Println("Received registration message from ", message.SenderId)
	err := controller.memberTable.Register(message.SenderId)
	if err == nil {
		ack := &connection.FileData{}
		ack.MessageType = connection.MessageType_ACK
		connectionHandler.Send(ack)
	} else {
		log.Println("Error handling registration request")
	}
}

func (controller *Controller) heartbeatHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	controller.memberTable.RecordBeat(message.SenderId)
}

func (controller *Controller) List() []string {
	return controller.memberTable.List()
}

func (controller *Controller) uploadHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	filepath := message.GetData()
	checksum := message.GetChecksum()
	chunks := ProtoToChunk(message.GetChunk())
	err := findAvailableNodes(chunks, controller.memberTable)
	if err != nil {
		log.Println("could  not assign nodes to chunks, ", err)
	}
	err = controller.fileMetadata.UploadChunks(filepath, chunks, checksum)
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

func (controller *Controller) ls(handler *connection.ConnectionHandler, connectionChan <-chan *connection.FileData, message *connection.FileData) {
	sendMessage := &connection.FileData{}
	sendMessage.MessageType = connection.MessageType_LS

	// TODO update to have ls logic to get files in directory path
	sendMessage.Data = "No files found"
	err := handler.Send(sendMessage)
	if err != nil {
		log.Println("Error sending ls data to client ", err)
	}
	ack := <-connectionChan
	if ack.MessageType != connection.MessageType_ACK_LS {
		// TODO retry if need to
		log.Println("Invalid ack for controller ls")
	}
}

func (controller *Controller) upload(handler *connection.ConnectionHandler) {
	// TODO add logic
	// 1. client connect and requests to upload chunks and names of chunks
	// filename, chunk name, chunk size, overall checksum, chunk checksum
	// 2. check in FileMetadata if path exists - if exists send error back
	// 3. If path does not exist
	//  - reserve the path with a pending state
	// 4. go to member table to ask for nodes that have enough space. -> if there is not enough space on enough nodes return message saying out of space
	// 5. if enough space member table returns list of storage nodes per chunk
	// 6 reserve space on member table -- subtract space out
	// 7. return chunk info to client
	//    - Wait for ack from client that client got list of chunks
	// 8. Get heart beat message from storage node that chunk is uploaded to storage node
	// 9. When N number of storage nodes respond back with ack of file then set status to Complete
	// end.
}
