package main

import (
	"dfs/client"
	"dfs/config"
	"dfs/controller"
	"dfs/storage"
	"flag"
	"log"
)

// Runs the application
func main() {
	nodeType := flag.String("type", "", "The type of node to run controller, client, storage")
	id := flag.String("id", "", "The identifier of the controller")
	command := flag.String("command", "", "The command for the client")
	host := flag.String("host", "", "The host")
	port := flag.Int("port", -1, "The port")

	// storage size is in MB, chunk size in config is in bytes
	storageSize := flag.Int64("storage_size", -1, "The storage size of the storage node")
	localPath :=  flag.String("local_path", "", "The path of the local file on disk")
	remotePath :=  flag.String("remote_path", "", "The path of the remote file on the dfs")
	flag.Parse()

	configFile, err := config.ConfigFromPath("./config.json")
	if err != nil {
		log.Fatalln("Failed to open config on controller ", err)
	}

	switch *nodeType {
	case "controller":
		runController(id, configFile)
	case "client":
		runClient(configFile, command, localPath, remotePath)
	case "storage":
		runStorage(id, storageSize, host, port, configFile, localPath)
	default:
		log.Fatalln("Unsupported node type")
	}
	log.Println("Shutting down ", *nodeType)
}

func runStorage(id *string, size *int64, host *string, port *int, config *config.Config, localPath *string) {
	if len(*id) < 1 || *size < 0 || len(*host) < 1 || *port < 0 {
		log.Fatalln("Invalid args for storage node")
	}
	newStorageNode := storage.NewStorageNode(*id, *size, *host, int32(*port), config, *localPath)
	newStorageNode.Start()
}

func runClient(config *config.Config, command *string, localPath *string, remotePath *string) {
	if len(*command) < 1 {
		log.Fatalln("Invalid args for client")
	}
	newClient := client.NewClient(config, *command, *localPath, *remotePath)
	newClient.Start()
}

func runController(id *string, config *config.Config) {
	if len(*id) < 1 {
		log.Fatalln("Invalid args for controller")
	}
	newController := controller.NewController(*id, config)
	newController.Start()
}
