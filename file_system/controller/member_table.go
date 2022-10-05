package controller

import (
	"dfs/connection"
	"errors"
	"log"
	"math/big"
	"sync"
	"time"
)

type MemberTable struct {
	members map[string]Member
	lock    sync.RWMutex
	running bool
}

type Member struct {
	status        bool
	lastBeat      time.Time
	availableSize *big.Int
	port          int32
	host          string
	handler       *connection.ConnectionHandler
}

func NewMemberTable() *MemberTable {
	memberTable := &MemberTable{}
	members := make(map[string]Member)
	memberTable.members = members
	memberTable.lock = sync.RWMutex{}
	memberTable.running = true
	go memberTable.failureDetection()
	return memberTable
}

func (memberTable *MemberTable) Register(id string, size *big.Int, host string, port int32, handler *connection.ConnectionHandler) error {
	memberTable.lock.Lock()
	defer memberTable.lock.Unlock()
	storageNode := Member{}
	storageNode.status = true
	storageNode.availableSize = size
	storageNode.port = port
	storageNode.host = host
	storageNode.handler = handler
	if memberTable.members[id].status {
		log.Printf("Member with id %s already exists", id)
		return errors.New("storage node already registered")
	}
	memberTable.members[id] = storageNode
	return nil
}

func (memberTable *MemberTable) List() []string {
	memberTable.lock.RLock()
	defer memberTable.lock.RUnlock()
	result := make([]string, 0)
	for member := range memberTable.members {
		if memberTable.members[member].status {
			result = append(result, member)
		}
	}
	return result
}

func (memberTable *MemberTable) RecordBeat(id string) {
	memberTable.lock.Lock()
	defer memberTable.lock.Unlock()
	member := memberTable.members[id]
	member.lastBeat = time.Now()
	member.status = true
	memberTable.members[id] = member
}

func (memberTable *MemberTable) failureDetection() {
	sleepTimer := time.Second * 5
	threshold := time.Second * 10
	for memberTable.running {
		for member := range memberTable.members {
			memberTable.lock.Lock()
			duration := time.Now().Sub(memberTable.members[member].lastBeat)
			if duration > threshold && memberTable.members[member].status {
				deadMember := memberTable.members[member]
				deadMember.status = false
				memberTable.members[member] = deadMember
				log.Println("Deactivated storage node: ", member)
			}
			memberTable.lock.Unlock()
		}
		time.Sleep(sleepTimer)
	}
}

func (memberTable *MemberTable) Shutdown() {
	log.Println("Shutting down member table")
	memberTable.running = false
	memberTable.lock.Lock()
	defer memberTable.lock.Unlock()
	for member := range memberTable.members {
		m := memberTable.members[member]
		if m.handler != nil {
			err := m.handler.Close()
			if err != nil {
				log.Printf("Failed to close socket for member %s on port %d \n", m.host, m.port)
			}
		}
	}
}
