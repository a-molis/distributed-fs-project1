package controller

import (
	"dfs/connection"
	"dfs/file_metadata"
	"errors"
	"math/big"
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
		newChunk.Num = c.Num
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
		newChunk.Num = c.Num
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
func findAvailableNodes(chunks []*file_metadata.Chunk, memberTable *MemberTable, numReplicas int) error {
	listOfNodes := make([]string, 0)
	if len(memberTable.members) < 1 {
		err := errors.New("Erorr finding any available nodes member table empty")
		if err != nil {
			return err
		}
	}
	for key := range memberTable.members {
		listOfNodes = append(listOfNodes, key)
	}
	// TODO update so this throws an error instead
	if numReplicas > len(memberTable.members) {
		numReplicas = len(memberTable.members)
	}
	for _, chunk := range chunks {

		// TODO handle case when storage nodes are full
		for j := 0; j < numReplicas; {
			idx := rand.Intn(len(listOfNodes))
			if memberTable.members[listOfNodes[idx]].availableSize.Cmp(big.NewInt(chunk.Size)) >= 0 && !contains(listOfNodes[idx], chunk.StorageNodes) {
				chunk.StorageNodes = append(chunk.StorageNodes, listOfNodes[idx])
				member := memberTable.members[listOfNodes[idx]]
				member.availableSize.Set(member.availableSize.Sub(member.availableSize, big.NewInt(chunk.Size)))
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
