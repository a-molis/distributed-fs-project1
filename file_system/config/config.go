package config

import (
	"P1-go-distributed-file-system/files_io"
	"encoding/json"
	"log"
)

type Config struct {
	ChunkSize int `json:"chunk_size"`
}

func ConfigFromPath(path string) (*Config, error) {
	configBytes, err := file_io.OpenFile(path)
	if err != nil {
		log.Printf("Error reading config %s %s\n", path, err)
	}
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Printf("Error converting config %s\n", err)
		return nil, err
	}
	return &config, nil
}
