package client

import (
	. "P1-go-distributed-file-system/config"
	"P1-go-distributed-file-system/connection"
	"fmt"
	"log"
	"strings"
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

func NewClient(config *Config, controllerHost string, controllerPort int, command string, args ...string) *Client {
	client := &Client{}
	client.config = config
	client.command = command
	client.args = args
	client.controllerHost = controllerHost
	client.controllerPort = controllerPort
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
	fileExists, err := connectionHandler.Receive()
	if err != nil {
		log.Fatalln("Error getting data from controller on client")
	}
	if fileExists.MessageType == connection.MessageType_ERROR {
		log.Fatalln("Path on remote file system already exists")
	} else if fileExists.MessageType != connection.MessageType_ACK {
		log.Fatalln("Received invalid message from controller on client put")
	}
}
