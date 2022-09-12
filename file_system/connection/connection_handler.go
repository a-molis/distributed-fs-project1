package connection

import (
	"encoding/binary"
	"net"

	"google.golang.org/protobuf/proto"
)

type ConnectionHandler struct {
	conn net.Conn
}

func NewConnectionHandler(conn net.Conn) *ConnectionHandler {
	m := &ConnectionHandler{
		conn: conn,
	}

	return m
}

func (m *ConnectionHandler) readN(buf []byte) error {
	bytesRead := uint64(0)
	for bytesRead < uint64(len(buf)) {
		n, err := m.conn.Read(buf[bytesRead:])
		if err != nil {
			return err
		}
		bytesRead += uint64(n)
	}
	return nil
}

func (m *ConnectionHandler) writeN(buf []byte) error {
	bytesWritten := uint64(0)
	for bytesWritten < uint64(len(buf)) {
		n, err := m.conn.Write(buf[bytesWritten:])
		if err != nil {
			return err
		}
		bytesWritten += uint64(n)
	}
	return nil
}

func (m *ConnectionHandler) Send(fileData *FileData) error {
	serialized, err := proto.Marshal(fileData)
	if err != nil {
		return err
	}

	prefix := make([]byte, 8)
	binary.LittleEndian.PutUint64(prefix, uint64(len(serialized)))
	m.writeN(prefix)
	m.writeN(serialized)

	return nil
}

func (m *ConnectionHandler) Receive() (*FileData, error) {
	prefix := make([]byte, 8)
	m.readN(prefix)

	payloadSize := binary.LittleEndian.Uint64(prefix)
	payload := make([]byte, payloadSize)
	m.readN(payload)

	fileData := &FileData{}
	err := proto.Unmarshal(payload, fileData)
	return fileData, err
}

func (m *ConnectionHandler) Close() {
	m.conn.Close()
}
