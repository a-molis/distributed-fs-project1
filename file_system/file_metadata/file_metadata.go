package file_metadata

import (
	"errors"
	"strings"
)

type FileMetadata struct {
	rootNode *Node
}

func NewFileMetaData() *FileMetadata {
	path := "/"
	rootNode := newNode(path)
	return &FileMetadata{rootNode: rootNode}
}

func (fileMetadata *FileMetadata) upload(path string) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, true)
	_, ok := directoryNode.files[fileName]
	if ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.Name = fileName
	file.Status = Pending
	directoryNode.files[fileName] = file
	return nil
}

func (fileMetadata *FileMetadata) UploadChunks(path string, chunks []*Chunk, checksum int32) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, true)
	_, ok := directoryNode.files[fileName]
	if ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.Name = fileName
	file.Status = Pending
	file.Chunks = chunks
	file.Checksum = checksum
	file.PendingChunks = len(chunks)
	directoryNode.files[fileName] = file
	return nil
}

func getNode(node *Node, path string, write bool) *Node {
	if path == "/" {
		return node
	}
	pathSplit := strings.Split(path, "/")
	directoryName := pathSplit[1]
	directoryNode, ok := node.dirs[directoryName]
	if !ok {
		if write {
			node.dirs[directoryName] = newNode(directoryName)
		}
		directoryNode = node.dirs[directoryName]
	}
	return getNode(directoryNode, strings.TrimPrefix(path, "/"+directoryName), write)
}

func (fileMetadata *FileMetadata) ls(path string) string {
	directoryNode := getNode(fileMetadata.rootNode, path, false)
	res := ""
	for dir := range directoryNode.dirs {
		res = res + dir + " "
	}
	for file := range directoryNode.files {
		res = res + file + " "
	}
	return strings.TrimSuffix(res, " ")
}

func (fileMetadata *FileMetadata) download(path string) []*Chunk {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, false)
	file := directoryNode.files[fileName]
	return file.Chunks
}

func (fileMetadata *FileMetadata) PathExists(path string) bool {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, false)
	_, ok := directoryNode.files[fileName]
	if ok {
		return true
	}
	return false
}

func (fileMetadata *FileMetadata) checkDirectoryNode(node *Node, path string) *Node {
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
	return getNode(directoryNode, strings.TrimPrefix(path, "/"+directoryName), false)
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
	Name          string
	Chunks        []*Chunk
	Status        Status
	Checksum      int32
	PendingChunks int
}

type Chunk struct {
	Name         string
	Size         int32
	Checksum     int32
	Status       Status
	StorageNodes []string
}

type Status int32

const (
	Pending Status = iota
	Complete
)
