package storage

import (
	"P1-go-distributed-file-system/config"
	"P1-go-distributed-file-system/connection"
	"log"
	"math/big"
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
