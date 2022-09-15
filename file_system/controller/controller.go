package controller

import (
	"P1-go-distributed-file-system/connection"
	"log"
)

type Controller struct {
	id          string
	port        int
	server      *connection.Server
	memberTable *MemberTable
	running     bool
}

func NewController(id string, port int) *Controller {
	controller := &Controller{}
	controller.id = id
	controller.port = port
	controller.memberTable = NewMemberTable()
	controller.running = true
	return controller
}

func (controller *Controller) Start() {
	controller.server = connection.NewServer(controller.port)
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
