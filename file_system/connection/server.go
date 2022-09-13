package connection

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
}

func NewServer(port int) *Server {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Unable to create server on port %d", port)
	}
	log.Printf("listening at port %d", port)
	return &Server{
		listener: listener,
	}
}

func (server *Server) NextConnectionHandler() (*ConnectionHandler, error) {
	conn, err := server.listener.Accept()
	if err != nil {
		log.Println("Unable to get next connection handler")
		return nil, err
	}
	return NewConnectionHandler(conn), nil
}
