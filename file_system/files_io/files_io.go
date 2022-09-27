package file_io

import (
	"P1-go-distributed-file-system/connection"
	"io/ioutil"
	"log"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open file %s\n", path)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Unable to close file %s\n", path)
		}
	}(file)
	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return readBytes, nil
}

func SendError(connectionHandler *connection.ConnectionHandler, errorMessage string) error {
	sendError := &connection.FileData{}
	sendError.MessageType = connection.MessageType_ERROR
	sendError.Data = errorMessage
	err := connectionHandler.Send(sendError)
	if err != nil {
		log.Println("Error sending error message ", errorMessage)
		return err
	}
	return nil
}

func SendAck(connectionHandler *connection.ConnectionHandler) error {
	sendAck := &connection.FileData{}
	sendAck.MessageType = connection.MessageType_ACK
	err := connectionHandler.Send(sendAck)
	if err != nil {
		log.Println("Error sending Ack message ")
		return err
	}
	return nil
}

func SendMessage(handler *connection.ConnectionHandler, messageType connection.MessageType) error {
	message := &connection.FileData{}
	message.MessageType = messageType
	err := handler.Send(message)
	if err != nil {
		log.Println("Error sending message of type ", messageType)
		return err
	}
	return nil
}
