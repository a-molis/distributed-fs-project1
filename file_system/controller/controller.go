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
	running bool
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
	for controller.running {
		message, err := connectionHandler.Receive()
		if err != nil {
			log.Println("Cannot receive message ", err)
		}
		if message.MessageType == connection.MessageType_REGISTRATION {
			controller.registerHandler(connectionHandler, message)
		}
		if message.MessageType == connection.MessageType_HEARTBEAT {
			controller.heartbeatHandler(connectionHandler, message)
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

func (controller *Controller) List() []string{
	return controller.memberTable.List()
}
