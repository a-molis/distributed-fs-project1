package file_metadata

import (
	"encoding/json"
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

func (fileMetadata *FileMetadata) Upload(path string) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, true)
	// TODO once getNode returns error add error handling
	_, ok := directoryNode.Files[fileName]
	if ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.Name = fileName
	file.Status = Pending
	directoryNode.Files[fileName] = file
	return nil
}

func (fileMetadata *FileMetadata) UploadChunks(path string, chunks map[string]*Chunk) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, true)
	// TODO once getNode returns error add error handling
	_, ok := directoryNode.Files[fileName]
	if ok {
		return errors.New("file already exists")
	}
	file := &File{}
	file.Name = fileName
	file.Status = Pending
	file.Chunks = chunks
	file.PendingChunks = len(chunks) * 3 //TODO make this times 3
	directoryNode.Files[fileName] = file
	return nil
}

func chunkArrayToMap(chunks []*Chunk) map[string]*Chunk {
	res := make(map[string]*Chunk)
	for _, c := range chunks {
		res[c.Name] = c
	}
	return res
}

func getNode(node *Node, path string, write bool) *Node {
	if path == "/" {
		return node
	}
	pathSplit := strings.Split(path, "/")
	directoryName := pathSplit[1]
	directoryNode, ok := node.Dirs[directoryName]
	if !ok {
		if write {
			node.Dirs[directoryName] = newNode(directoryName)
			directoryNode = node.Dirs[directoryName]
		} else {
			// TODO refactor to return tuple with (*Node, error)
			return nil
		}
	}
	return getNode(directoryNode, strings.TrimPrefix(path, "/"+directoryName), write)
}

func (fileMetadata *FileMetadata) Ls(path string) string {
	directoryNode := getNode(fileMetadata.rootNode, path, false)
	res := ""
	for dir := range directoryNode.Dirs {
		res = res + dir + " "
	}
	for file := range directoryNode.Files {
		res = res + file + " "
	}
	return strings.TrimSuffix(res, " ")
}

func (fileMetadata *FileMetadata) Download(path string) (map[string]*Chunk, []byte, error) {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, false)
	if directoryNode == nil {
		return nil, nil, errors.New("File does not exist")
	}
	file := directoryNode.Files[fileName]
	return file.Chunks, file.Checksum, nil
}

func (fileMetadata *FileMetadata) PathExists(path string) bool {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, false)
	if directoryNode == nil {
		return false
	}
	_, ok := directoryNode.Files[fileName]
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
	directoryNode, ok := node.Dirs[directoryName]
	if !ok {
		node.Dirs[directoryName] = newNode(directoryName)
		directoryNode = node.Dirs[directoryName]
	}
	return getNode(directoryNode, strings.TrimPrefix(path, "/"+directoryName), false)
}

func (fileMetadata *FileMetadata) GetBytes() ([]byte, error) {
	res, err := json.Marshal(fileMetadata.rootNode)
	return res, err
}

func (fileMetadata *FileMetadata) LoadBytes(bytes []byte) error {
	err := json.Unmarshal(bytes, fileMetadata.rootNode)
	return err
}

func (fileMetadata *FileMetadata) UpdateChecksum(path string, sum []byte) error {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, false)
	if directoryNode == nil {
		return errors.New("File does not exist")
	}
	file := directoryNode.Files[fileName]
	file.Checksum = sum
	return nil
}

func (fileMetadata *FileMetadata) HeartbeatHandler(path string, chunk string) bool {
	pathSplit := strings.Split(path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	directoryPath := strings.Replace(path, fileName, "", -1)
	directoryNode := getNode(fileMetadata.rootNode, directoryPath, false)
	file := directoryNode.Files[fileName]
	file.Chunks[chunk].Status = Complete
	file.PendingChunks = file.PendingChunks - 1 //TODO This should really be threadsafe
	if file.PendingChunks <= 0 {
		file.Status = Complete
		return true
	}
	return false
}

type Node struct {
	Path  string
	Dirs  map[string]*Node
	Files map[string]*File
}

func newNode(path string) *Node {
	node := &Node{}
	node.Path = path
	node.Dirs = make(map[string]*Node)
	node.Files = make(map[string]*File)
	return node
}

// TODO chunks to use set data type
type File struct {
	Name          string
	Chunks        map[string]*Chunk //TODO Map //this is used for changing the status of chunk 246
	Status        Status
	Checksum      []byte
	PendingChunks int
}

type Chunk struct {
	Name         string
	Size         int64
	Checksum     int32
	Status       Status
	StorageNodes []string //TODO this needs to be map for better pending state
	Num          int32
}

type Status int32

const (
	Pending Status = iota
	Complete
)
