package connection

import (
	"log"
	"net"
	"testing"
	"time"
)

func TestBasicServer(t *testing.T) {

	path := "some/path"

	receivedPath := ""

	server := NewServer(9999)

	//bit for the client
	go func(path string) {
		message := &FileData{}
		message.Path = path
		message.MessageType = MessageType_PUT
		conn, err := net.Dial("tcp", "localHost:9999") // connect to localhost port 9999
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		connectionHandler := NewConnectionHandler(conn)
		connectionHandler.Send(message)
		return
	}(path)

	go func(receivedPath *string) {
		connHandler, err := server.NextConnectionHandler()
		if err != nil {
			t.Fatalf("got error %s", err)
		}

		receivedMessage, err :=connHandler.Receive()
		if err != nil {
			t.Fatalf("got error %s", err)
		}
		*receivedPath = receivedMessage.Path
	}(&receivedPath)

	time.Sleep(time.Second * 1)

	if path != receivedPath {
		t.Fatalf("the message path dont match %s", receivedPath)
	}
}

func TestBasicClient(t *testing.T) {

	path := "some/path"

	receivedPath := ""

	var port int = 9998

	server := NewServer(port)

	//bit for the client
	go func(path string) {
		message := &FileData{}
		message.Path = path
		message.MessageType = MessageType_PUT
		connHand, _ := NewClient("localHost", port)
		connHand.Send(message)
		return
	}(path)

	go func(receivedPath *string) {
		connHandler, err := server.NextConnectionHandler()
		if err != nil {
			t.Fatalf("got error %s", err)
		}

		receivedMessage, err :=connHandler.Receive()
		if err != nil {
			t.Fatalf("got error %s", err)
		}
		*receivedPath = receivedMessage.Path
	}(&receivedPath)

	time.Sleep(time.Second * 1)

	if path != receivedPath {
		t.Fatalf("the message path dont match %s", receivedPath)
	}

	return
}