package storage

import (
	"P1-go-distributed-file-system/connection"
	"log"
)

type StorageNode struct {
	id                string
	controllerPort    int
	controllerHost    string
	connectionHandler *connection.ConnectionHandler
}

func NewStorageNode(id string, host string, port int) *StorageNode {
	storageNode := &StorageNode{}
	storageNode.id = id
	storageNode.controllerPort = port
	storageNode.controllerHost = host
	return storageNode
}

func (storageNode *StorageNode) Start() {
	connectionHandler, err := connection.NewClient(storageNode.controllerHost, storageNode.controllerPort)
	if err != nil {
		log.Printf("Error while startig the storageNode %s", err)
		return
	}
	storageNode.connectionHandler = connectionHandler

	storageNode.register()
}

func (storageNode *StorageNode) register() {
	message := &connection.FileData{}
	message.MessageType = connection.MessageType_REGISTRATION
	message.SenderId = storageNode.id
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
