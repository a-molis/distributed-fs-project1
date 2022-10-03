package client

import (
	"dfs/config"
	"dfs/controller"
	"dfs/storage"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

// TODO add test for duplicate file uploaded

func TestBasicClient(t *testing.T) {

	controllerHost := "localhost"
	storageHost := "localhost"
	var storagePort0 int32 = 12032
	var storagePort1 int32 = 12033

	testId := "storageTestId"
	testId2 := "storageTestId2"

	savePathStorageNode0 := "sn0"
	savePathStorageNode1 := "sn1"

	var members []string

	// Chunk size is in bytes, storage node size in MB
	var chunkSize int64 = 5000000
	var storageSize int64 = 1000000

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ChunkSize = chunkSize
	testConfig.ControllerPort = 12031
	testConfig.ControllerHost = controllerHost

	go func(receivedMembers *[]string) {
		testController := controller.NewController("testId", testConfig)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(testId, storageSize, storageHost, storagePort0, testConfig, savePathStorageNode0)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}(testId)

	go func(id string) {
		storageNode := storage.NewStorageNode(testId2, storageSize, storageHost, storagePort1, testConfig, savePathStorageNode1)
		storageNode.Start()
	}(testId2)

	var clientLsError error = nil
	go func() {
		testClient := NewClient(testConfig, "ls")
		clientLsError = testClient.Run()
	}()
	time.Sleep(time.Second * 5)

	if clientLsError != nil {
		t.Fatalf("client test failed %s", clientLsError)
	}

	return
}

func TestClientUploadData(t *testing.T) {
	defer func() {
		name := "this_test_path_footxt_"
		for i := 0; i < 10; i++ {
			newName := fmt.Sprintf("%s%d", name, i)
			for j := 0; j < 5; j++ {
				fullPath := fmt.Sprintf("./sn%d/%s", j, newName)
				os.Remove(fullPath)
				dirName := fmt.Sprintf("./sn%d", j)
				os.Remove(dirName)
			}
		}
	}()

	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12050
	var storagePort0 int32 = 12051
	var storagePort1 int32 = 12052
	var storagePort2 int32 = 12053
	var storagePort3 int32 = 12054

	var size int64 = 1000000
	var chunkSize int64 = 5000000

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	uploadPath := "/this/test/path/foo.txt"
	localPath := "testdata/testFile.txt"

	testStorageNode0 := "testStorageNode0"
	testStorageNode1 := "testStorageNode1"
	testStorageNode2 := "testStorageNode2"
	testStorageNode3 := "testStorageNode3"

	savePathStorageNode0 := "sn0"
	savePathStorageNode1 := "sn1"
	savePathStorageNode2 := "sn2"
	savePathStorageNode3 := "sn3"

	var members []string

	go func(receivedMembers *[]string) {
		testController := controller.NewController("testId", testConfig)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode0, size, storageHost, storagePort0, testConfig, savePathStorageNode0)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode1, size, storageHost, storagePort1, testConfig, savePathStorageNode1)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode2, size, storageHost, storagePort2, testConfig, savePathStorageNode2)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode3, size, storageHost, storagePort3, testConfig, savePathStorageNode3)
		storageNode.Start()
	}()

	time.Sleep(time.Second * 3)
	var clientPutError error = nil
	go func() {
		testClient := NewClient(testConfig, "put", uploadPath, localPath)
		clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 30)

	// TODO complete test to validate file is saved
	if clientPutError != nil {
		t.Fatalf("client test failed %s ", clientPutError)
	}

}

func TestClientDownloadSimple(t *testing.T) {
	savePath := "testdata/downloadCopySimple.txt"
	defer func(savePath string) {
		err := os.Remove(savePath)
		if err != nil {
			log.Printf("Unable to delete file at %s", savePath)
		}
	}(savePath)
	defer func() {
		name := "this_test_path_footxt_"
		for i := 0; i < 10; i++ {
			newName := fmt.Sprintf("%s%d", name, i)
			for j := 0; j < 5; j++ {
				fullPath := fmt.Sprintf("./sn%d/%s", j, newName)
				os.Remove(fullPath)
				dirName := fmt.Sprintf("./sn%d", j)
				os.Remove(dirName)
			}
		}
	}()

	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12043
	var storagePort0 int32 = 12044
	var storagePort1 int32 = 12045
	var storagePort2 int32 = 12046
	var storagePort3 int32 = 12047
	var storagePort4 int32 = 12048
	var size int64 = 1000000
	var chunkSize int64 = 5000000

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	remotePath := "/this/test/path/foo.txt"
	localPath := "testdata/testFile.txt"

	testStorageNode0 := "testStorageNode0"
	testStorageNode1 := "testStorageNode1"
	testStorageNode2 := "testStorageNode2"
	testStorageNode3 := "testStorageNode3"
	testStorageNode4 := "testStorageNode4"

	savePathStorageNode0 := "sn0"
	savePathStorageNode1 := "sn1"
	savePathStorageNode2 := "sn2"
	savePathStorageNode3 := "sn3"
	savePathStorageNode4 := "sn4"

	var members []string

	go func(receivedMembers *[]string) {
		testController := controller.NewController("testId", testConfig)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode0, size, storageHost, storagePort0, testConfig, savePathStorageNode0)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode1, size, storageHost, storagePort1, testConfig, savePathStorageNode1)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode2, size, storageHost, storagePort2, testConfig, savePathStorageNode2)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode3, size, storageHost, storagePort3, testConfig, savePathStorageNode3)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode4, size, storageHost, storagePort4, testConfig, savePathStorageNode4)
		storageNode.Start()
	}()

	time.Sleep(time.Second * 3)
	var clientPutError error = nil
	go func() {
		testClient := NewClient(testConfig, "put", remotePath, localPath)
		clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 2)

	var clientGetError error = nil
	go func() {
		testClient := NewClient(testConfig, "get", remotePath, savePath)
		clientGetError = testClient.Run()
	}()

	time.Sleep(time.Second * 10)

	// TODO complete test to validate file is saved
	if clientGetError != nil || clientPutError != nil {
		t.Fatalf("client test failed %s, %s", clientGetError, clientPutError)
	}
	return
}
