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
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath)
	_, ok := directoryNode.files[fileName]
	if ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.name = fileName
	file.status = Pending
	directoryNode.files[fileName] = file
	return nil
}

func (fileMetadata *FileMetadata) uploadChunks(path string, chunks []Chunk) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath)
	_, ok := directoryNode.files[fileName]
	if ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.name = fileName
	file.status = Pending
	file.chunks = chunks
	file.pendingChunks = len(chunks)
	directoryNode.files[fileName] = file
	return nil
}

func getNode(node *Node, path string) *Node {
	if path == "/" {
		return node
	}
	pathSplit := strings.Split(path, "/")
	directoryName := pathSplit[1]
	directoryNode, ok := node.dirs[directoryName]
	if !ok {
		node.dirs[directoryName] = newNode(directoryName)
		directoryNode = node.dirs[directoryName]
	}
	return getNode(directoryNode, strings.TrimPrefix(path, "/" + directoryName))
}

func (fileMetadata *FileMetadata) ls(path string) string {
	directoryNode := getNode(fileMetadata.rootNode, path)
	res := ""
	for dir := range directoryNode.dirs {
	res = res + dir + " "
	}
	for file := range directoryNode.files {
		res = res + file + " "
	}
	return strings.TrimSuffix(res, " ")
}

func (fileMetadata *FileMetadata) download(path string) []Chunk {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath)
	file := directoryNode.files[fileName]
	return file.chunks
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

// TODO chunks to use set data type
type File struct {
	name          string
	chunks        []Chunk
	status        Status
	pendingChunks int
}

type Chunk struct {
	checksum     string
	status       Status
	storageNodes []string
}

type Status int32

const (
	Pending Status = iota
	Complete
)
