package controller

import (
	"P1-go-distributed-file-system/connection"
	"P1-go-distributed-file-system/file_metadata"
	"math/rand"
)

func ProtoToChunk(chunks []*connection.Chunk) []*file_metadata.Chunk {
	res := make([]*file_metadata.Chunk, 0)
	for _, c := range chunks {
		newChunk := &file_metadata.Chunk{}
		newChunk.Name = c.Name
		newChunk.Checksum = c.Checksum
		newChunk.Status = file_metadata.Pending
		newChunk.Size = c.Size
		res = append(res, newChunk)
	}
	return res
}

func chunkToProto(chunks []*file_metadata.Chunk, memberTable *MemberTable) []*connection.Chunk {
	res := make([]*connection.Chunk, 0)
	for _, c := range chunks {
		newChunk := &connection.Chunk{}
		newChunk.Name = c.Name
		newChunk.Checksum = c.Checksum
		newChunk.Size = c.Size
		newChunk.Nodes = generateNodes(c.StorageNodes, memberTable)
		res = append(res, newChunk)
	}
	return res
}

func generateNodes(ids []string, memberTable *MemberTable) []*connection.Node {
	res := make([]*connection.Node, 0)
	for _, id := range ids {
		newNode := &connection.Node{}
		newNode.Id = id
		newNode.Port = memberTable.members[id].port
		newNode.Hostname = memberTable.members[id].host
		res = append(res, newNode)
	}
	return res
}

// TODO move this to function in member table
func findAvailableNodes(chunks []*file_metadata.Chunk, memberTable *MemberTable) error {
	listOfNodes := make([]string, 0)
	for key := range memberTable.members {
		listOfNodes = append(listOfNodes, key)
	}

	for _, chunk := range chunks {

		for j := 0; j < 3; {
			idx := rand.Intn(len(listOfNodes))
			if memberTable.members[listOfNodes[idx]].availableSize > chunk.Size && !contains(listOfNodes[idx], chunk.StorageNodes) {
				chunk.StorageNodes = append(chunk.StorageNodes, listOfNodes[idx])
				member := memberTable.members[listOfNodes[idx]]
				member.availableSize = member.availableSize - chunk.Size
				memberTable.members[listOfNodes[idx]] = member
				j++
			}
		}
	}

	return nil
}

func contains(key string, list []string) bool {
	for _, word := range list {
		if word == key {
			return true
		}
	}
	return false
}
