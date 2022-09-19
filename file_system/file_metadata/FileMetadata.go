package file_metadata

type FileMetadata struct {
	rootNode Node
}

type Node struct {
	path  string
	dirs  map[string]Node
	files map[string]File
}

type File struct {
	name string
}
