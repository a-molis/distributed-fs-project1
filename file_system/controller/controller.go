package controller

import (
	"P1-go-distributed-file-system/connection"
	"log"
)

type Controller struct {
	id     string
	port   int
	server *connection.Server
	memberTable *MemberTable
}

func NewController(id string, port int) *Controller {
	controller := &Controller{}
	controller.id = id
	controller.port = port
	controller.memberTable = NewMemberTable()
	return controller
}

func (controller *Controller) Start() {
	controller.server = connection.NewServer(controller.port)
	controller.listen()
}

func (controller *Controller) listen() {
	connectionHandler, err := controller.server.NextConnectionHandler()
	if err != nil {
		log.Println("Cannot get the connection handler ", err)
	}
	go controller.handleConnection(connectionHandler)
}

func (controller *Controller) handleConnection(connectionHandler *connection.ConnectionHandler) {
	message, err := connectionHandler.Receive()
	if err != nil {
		log.Println("Cannot receive message ", err)
	}
	if message.MessageType == connection.MessageType_REGISTRATION {
		controller.registerHandler(connectionHandler, message)
	}
}

func (controller *Controller) registerHandler(connectionHandler *connection.ConnectionHandler, message *connection.FileData) {
	log.Println("Received registration message from ", message.SenderId)
	err := controller.memberTable.register(message.SenderId)
	if err == nil {
		ack := &connection.FileData{}
		ack.MessageType = connection.MessageType_ACK
		connectionHandler.Send(ack)
	} else {
		log.Println("Error handling registration request")
	}
}

func (controller *Controller) List() []string{
	return controller.memberTable.List()
}
