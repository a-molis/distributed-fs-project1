package client

import (
	"dfs/config"
	"dfs/controller"
	"dfs/storage"
	"fmt"
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
		storageNode := storage.NewStorageNode(testId, storageSize, storageHost, storagePort0, testConfig)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}(testId)

	go func(id string) {
		storageNode := storage.NewStorageNode(testId2, storageSize, storageHost, storagePort1, testConfig)
		storageNode.Start()
	}(testId2)

	go func(id string) {
		testClient := NewClient(testConfig, "ls")
		testClient.Start()
	}(testId2)
	time.Sleep(time.Second * 5)

	if 1 != 1 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestClientUpload(t *testing.T) {
	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12027
	var storagePort0 int32 = 12028
	var storagePort1 int32 = 12029
	var storagePort2 int32 = 12034
	var storagePort3 int32 = 12035
	var storagePort4 int32 = 12036
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
	testStorageNode4 := "testStorageNode4"
	testClientId0 := "clientId0"

	var members []string

	go func(receivedMembers *[]string) {
		testController := controller.NewController("testId", testConfig)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode0, size, storageHost, storagePort0, testConfig)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode1, size, storageHost, storagePort1, testConfig)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode2, size, storageHost, storagePort2, testConfig)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode3, size, storageHost, storagePort3, testConfig)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode4, size, storageHost, storagePort4, testConfig)
		storageNode.Start()
	}()

	time.Sleep(time.Second * 3)
	go func(port int, id string) {
		testClient := NewClient(testConfig, "put", uploadPath, localPath)
		testClient.Start()
	}(controllerPort, testClientId0)

	time.Sleep(time.Second * 2)

	// TODO complete test to validate file is saved
	if 1 != 1 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestClientUploadData(t *testing.T) {
	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12050
	var storagePort0 int32 = 12051

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
	testClientId0 := "clientId0"

	var members []string

	go func(receivedMembers *[]string) {
		testController := controller.NewController("testId", testConfig)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode0, size, storageHost, storagePort0, testConfig)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}()

	time.Sleep(time.Second * 3)
	go func(port int, id string) {
		testClient := NewClient(testConfig, "put", uploadPath, localPath)
		testClient.Start()
	}(controllerPort, testClientId0)

	time.Sleep(time.Second * 10)

	name := "this_test_path_footxt_"
	for i := 0; i<10 ; i++ {
		//newName := name + string(i)
		newName := fmt.Sprintf("%s%d", name, i)
		os.Remove(newName)
	}

	// TODO complete test to validate file is saved
	if 1 != 1 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

}



func TestClientDownloadSimple(t *testing.T) {
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
	savePath := "testdata/downloadCopySimple.txt"

	testStorageNode0 := "testStorageNode0"
	testStorageNode1 := "testStorageNode1"
	testStorageNode2 := "testStorageNode2"
	testStorageNode3 := "testStorageNode3"
	testStorageNode4 := "testStorageNode4"
	testClientId0 := "clientId0"

	var members []string

	go func(receivedMembers *[]string) {
		testController := controller.NewController("testId", testConfig)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode0, size, storageHost, storagePort0, testConfig)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode1, size, storageHost, storagePort1, testConfig)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode2, size, storageHost, storagePort2, testConfig)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode3, size, storageHost, storagePort3, testConfig)
		storageNode.Start()
	}()

	go func() {
		storageNode := storage.NewStorageNode(testStorageNode4, size, storageHost, storagePort4, testConfig)
		storageNode.Start()
	}()

	time.Sleep(time.Second * 3)
	go func(port int, id string) {
		testClient := NewClient(testConfig, "put", remotePath, localPath)
		testClient.Start()
	}(controllerPort, testClientId0)

	time.Sleep(time.Second * 2)

	go func(port int, id string) {
		testClient := NewClient(testConfig, "get", remotePath, savePath)
		testClient.Start()
	}(controllerPort, testClientId0)

	time.Sleep(time.Second * 4)

	// TODO complete test to validate file is saved
	if 1 != 1 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}
	return
}
