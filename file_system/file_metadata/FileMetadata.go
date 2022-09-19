package file_metadata

import (
	"errors"
	"strings"
)

type FileMetadata struct {
	rootNode *Node
}

func newFileMetaData() *FileMetadata {
	path := "/"
	rootNode := newNode(path)
	return &FileMetadata{rootNode: rootNode}
}

func (fileMetadata *FileMetadata) upload(path string) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, "/"+fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath)
	_, ok := directoryNode.files[fileName]
	if !ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.name = fileName
	file.status = Pending
	directoryNode.files[fileName] = file
	return nil
}

func getNode(node *Node, path string) *Node {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) == 1 {
		return node
	}
	directoryName := pathSplit[0]
	directoryNode, ok := node.dirs[directoryName]
	if !ok {
		node.dirs[directoryName] = newNode(directoryName)
		directoryNode = node.dirs[directoryName]
	}
	return getNode(directoryNode, strings.TrimPrefix(path, directoryName+"/"))
}

type Node struct {
	path  string
	dirs  map[string]*Node
	files map[string]*File
}

func newNode(path string) *Node {
	node := &Node{}
	node.path = path
	node.dirs = make(map[string]*Node)
	node.files = make(map[string]*File)
	return node
}

type File struct {
	name         string
	checksum     string
	status       Status
	storageNodes string[]
}

type Status int32

const (
	Pending Status = iota
	Complete
)
