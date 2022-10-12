package atomic_data

import (
	"dfs/connection"
	"sync"
)

type BlockingConnection struct {
	ConHandler *connection.ConnectionHandler
	Mu         sync.Mutex
}

type BlockingHandlerMap struct {
	hashMap map[string]*BlockingConnection
	Lock    sync.Mutex
}

func NewBlockingHandlerMap() *BlockingHandlerMap {
	b := &BlockingHandlerMap{}
	b.hashMap = make(map[string]*BlockingConnection)
	return b
}

func (bMap *BlockingHandlerMap) Put(host string, conn *BlockingConnection) {
	bMap.hashMap[host] = conn
}

func (bMap *BlockingHandlerMap) Get(host string) (*BlockingConnection, bool) {
	blockingConnection, ok := bMap.hashMap[host]
	return blockingConnection, ok
}
