package connection

import (
	"fmt"
	"log"
	"net"
	"time"
)

func NewClient(host string, port int) (*ConnectionHandler, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	return createClient(host, port, err, conn)
}

func NewClientTimeout(host string, port int, timeout time.Duration) (*ConnectionHandler, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	return createClient(host, port, err, conn)
}

func createClient(host string, port int, err error, conn net.Conn) (*ConnectionHandler, error) {
	if err != nil {
		log.Printf("received error while creating the client %s", err)
		return nil, err
	}
	log.Printf("connected to host %s, at port %d", host, port)
	return NewConnectionHandler(conn), nil
}
