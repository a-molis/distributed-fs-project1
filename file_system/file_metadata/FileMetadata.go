package file_metadata

type FileMetadata struct {
	rootNode *Node
}

func newFileMetaData() *FileMetadata {
	rootNode := &Node{}
	rootNode.path = "/"
	rootNode.dirs = make(map[string]Node)
	rootNode.files = make(map[string]File)
	return &FileMetadata{rootNode: rootNode}
}

func (fileMetadata *FileMetadata) upload(path string) (string[], error) {

}

type Node struct {
	path  string
	dirs  map[string]Node
	files map[string]File
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
