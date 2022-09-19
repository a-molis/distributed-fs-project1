package controller

import (
	"errors"
	"log"
	"sync"
	"time"
)

type MemberTable struct {
	members map[string]Member
	lock    sync.RWMutex
}

type Member struct {
	status        bool
	lastBeat      time.Time
	availableSize int
}

func NewMemberTable() *MemberTable {
	memberTable := &MemberTable{}
	members := make(map[string]Member)
	memberTable.members = members
	memberTable.lock = sync.RWMutex{}
	go memberTable.failureDetection()
	return memberTable
}

func (memberTable *MemberTable) Register(id string) error {
	memberTable.lock.Lock()
	defer memberTable.lock.Unlock()
	storageNode := Member{}
	storageNode.status = true
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
	memberTable.members[id] = member
}

func (memberTable *MemberTable) failureDetection() {
	sleepTimer := time.Second * 5
	threshold := time.Second * 10
	for {
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
