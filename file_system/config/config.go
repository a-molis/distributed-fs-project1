package config

import (
	file_io "dfs/files_io"
	"encoding/json"
	"log"
)

type Config struct {
	ChunkSize      int64  `json:"chunk_size"`
	NumReplicas    int    `json:"num_replicas"`
	ControllerHost string `json:"controller_host"`
	ControllerPort int    `json:"controller_port"`
}

func ConfigFromPath(path string) (*Config, error) {
	configBytes, err := file_io.ReadFile(path)
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
