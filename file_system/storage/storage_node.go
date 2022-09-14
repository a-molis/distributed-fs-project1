package storage

import (
	"P1-go-distributed-file-system/connection"
	"log"
	"time"
)

type StorageNode struct {
	id                string
	size				int32
	controllerPort    int
	controllerHost    string
	connectionHandler *connection.ConnectionHandler
}

func NewStorageNode(id string, size int32, host string, port int) *StorageNode {
	storageNode := &StorageNode{}
	storageNode.id = id
	storageNode.size = size
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

	alive := true
	go storageNode.heartbeat(&alive)
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

func (storageNode *StorageNode) heartbeat(alive *bool) {
	var heartBeatRate = time.Second * 5
	message := &connection.FileData{}
	message.MessageType = connection.MessageType_HEARTBEAT
	message.SenderId = storageNode.id
	message.Size = storageNode.size
	for *alive {
		storageNode.connectionHandler.Send(message)
		time.Sleep(heartBeatRate)
	}
}
