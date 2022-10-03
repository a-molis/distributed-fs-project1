package connection

import (
	"fmt"
	"log"
	"net"
)

func NewClient(host string, port int) (*ConnectionHandler, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Printf("received error while creating the client %s", err)
		return nil, err
	}
	log.Printf("connected to host %s, at port %d", host, port)
	return NewConnectionHandler(conn), nil
}
