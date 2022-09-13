package controller

import (
	"errors"
	"log"
)

type MemberTable struct {
	members map[string]Member
}

type Member struct {
	status bool
}

func NewMemberTable() *MemberTable {
	memberTable := &MemberTable{}
	members := make(map[string]Member)
	memberTable.members = members
	return memberTable
}

func (memberTable *MemberTable) register(id string) error {
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
	result := make([]string, 0)
	for member := range memberTable.members {
		if memberTable.members[member].status {
			result = append(result, member)
		}
	}
	return result
}
