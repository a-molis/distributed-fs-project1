package controller

import (
	"dfs/client"
	"dfs/config"
	"dfs/storage"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	testId := "storageTestId"
	var storagePort0 int32 = 12041
	savePathStorageNode0 := "sn0"

	var members []string

	var size int64 = 10

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ControllerHost = "localhost"
	testConfig.ControllerPort = 12020

	//bit for the controller
	go func(receivedMembers *[]string) {
		controller := NewController("testId", testConfig)
		go controller.Start()
		time.Sleep(time.Second * 3)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig, savePathStorageNode0)
		storageNode.Start()
	}(testId)

	time.Sleep(time.Second * 6)

	if len(members) == 0 || members[0] != testId {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestHeartBeatInactive(t *testing.T) {
	testId := "storageTestId"
	var storagePort0 int32 = 12038
	savePathStorageNode0 := "sn0"

	var members []string

	var size int64 = 10

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ControllerHost = "localhost"
	testConfig.ControllerPort = 12021

	//bit for the controller
	go func(receivedMembers *[]string) {
		controller := NewController("testId", testConfig)
		go controller.Start()
		time.Sleep(time.Second * 20)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig, savePathStorageNode0)
		storageNode.Start()
		storageNode.Shutdown()
	}(testId)

	time.Sleep(time.Second * 25)

	if len(members) != 0 {
		t.Fatalf("the registered node did not switch to inactive")
	}

	return
}

func TestHeartBeat(t *testing.T) {
	testId := "storageTestId"
	var storagePort0 int32 = 12037
	savePathStorageNode0 := "sn0"

	var members []string

	var size int64 = 10

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ControllerHost = "localhost"
	testConfig.ControllerPort = 12022

	//bit for the controller
	go func(receivedMembers *[]string) {
		controller := NewController("testId", testConfig)
		go controller.Start()
		time.Sleep(time.Second * 15)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig, savePathStorageNode0)
		storageNode.Start()
	}(testId)

	time.Sleep(time.Second * 20)

	if len(members) == 0 || members[0] != testId {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestHeartBeatMultiNode(t *testing.T) {
	testId := "storageTestId"
	testId2 := "storageTestId2"
	var storagePort0 int32 = 12039
	var storagePort1 int32 = 12040
	savePathStorageNode0 := "sn0"
	savePathStorageNode1 := "sn1"

	var members []string

	var size int64 = 10

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ControllerHost = "localhost"
	testConfig.ControllerPort = 12023

	//bit for the controller
	go func(receivedMembers *[]string) {
		controller := NewController("testId", testConfig)
		go controller.Start()
		time.Sleep(time.Second * 20)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig, savePathStorageNode0)
		go storageNode.Start()
		time.Sleep(time.Second * 1)
		storageNode.Shutdown()
	}(testId)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort1, testConfig, savePathStorageNode1)
		storageNode.Start()
	}(testId2)

	time.Sleep(time.Second * 25)

	if len(members) == 0 || members[0] != testId2 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestSaveFileMetadata(t *testing.T) {
	var file string
	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	err = os.Remove(testConfig.ControllerPath)
	if err != nil {
		log.Println("Unable to remove controller backup ", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)
	testConfig.ControllerHost = "localhost"
	testConfig.ControllerPort = 12060

	//bit for the controller
	go func(file *string) {
		controller := NewController("testId", testConfig)
		go controller.Start()
		err := controller.fileMetadata.Upload("/foo/something/file.txt")
		if err != nil {
			t.Error("Error uploading data ", err)
			return
		}
		err = controller.fileMetadata.Upload("/foo/something2/file2.txt")
		if err != nil {
			t.Error("Error uploading data")
			return
		}
		controller.SaveFileMetadata()
		controller.shutdown()
		testConfig.ControllerPort = 12061
		controller = NewController("testId", testConfig)
		go controller.Start()
		err = controller.LoadFileMetadata()
		if err != nil {
			t.Error("Error loading file metadata")
			return
		}
		*file, _ = controller.fileMetadata.Ls("/foo/something/")
	}(&file)

	time.Sleep(time.Second * 6)

	if !strings.Contains(file, "file.txt") {
		t.Fatalf("loaded file metadata is incorrect, %s", file)
	}

	return
}

func TestControllerCrashRecovery(t *testing.T) {
	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	defer func() {
		name := "this_test_path_footxt_"
		for i := 0; i < 10; i++ {
			newName := fmt.Sprintf("%s%d", name, i)
			for j := 0; j < 5; j++ {
				fullPath := fmt.Sprintf("./sn%d/%s", j, newName)
				os.Remove(fullPath)
				dirName := fmt.Sprintf("./sn%d", j)
				os.RemoveAll(dirName)
			}
		}
	}()

	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12073
	var storagePort0 int32 = 12074
	var storagePort1 int32 = 12075
	var storagePort2 int32 = 12076
	var storagePort3 int32 = 12078

	var size int64 = 1000000
	var chunkSize int64 = 5000000

	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)

	testFile := "foo.txt"
	cat := "cat"
	testPath := "/this/test"
	lsPath := filepath.Join(testPath, cat)
	uploadPath := filepath.Join(lsPath, testFile)
	localPath := "../client/testdata/testFile.txt"

	testStorageNode0 := "testStorageNode0"
	testStorageNode1 := "testStorageNode1"
	testStorageNode2 := "testStorageNode2"
	testStorageNode3 := "testStorageNode3"

	savePathStorageNode0 := "sn0"
	savePathStorageNode1 := "sn1"
	savePathStorageNode2 := "sn2"
	savePathStorageNode3 := "sn3"
	testControllerId0 := "testController0"
	testControllerId1 := "testController1"

	var members []string

	testController0 := NewController(testControllerId0, testConfig)
	go func(receivedMembers *[]string) {
		testController0.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController0.List()
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
		testClient := client.NewClient(testConfig, "put", uploadPath, localPath)
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 3)
	var clientLsError0 error = nil
	var lsResult0 *string
	go func() {
		testClient := client.NewClient(testConfig, "ls", uploadPath)
		result, err := testClient.Run()
		clientLsError0 = err
		lsResult0 = result
	}()

	time.Sleep(time.Second * 10)
	testController0.shutdown()
	time.Sleep(time.Second * 15)
	go func() {
		testController1 := NewController(testControllerId1, testConfig)
		testController1.Start()
	}()
	time.Sleep(time.Second * 10)

	var clientLsError1 error = nil
	var lsResult1 *string
	go func() {
		testClient := client.NewClient(testConfig, "ls", uploadPath)
		result, err := testClient.Run()
		clientLsError1 = err
		lsResult1 = result
	}()

	time.Sleep(time.Second * 5)

	if clientPutError != nil || clientLsError0 != nil || clientLsError1 != nil {
		t.Fatalf("client test failed %s %s %s ",
			clientPutError, clientLsError0, clientLsError1)
	}

	if clientPutError != nil || clientLsError1 != nil {
		t.Fatalf("client test failed %s %s ",
			clientPutError, clientLsError1)
	}
	if !strings.Contains(*lsResult0, uploadPath) {
		t.Fatalf("Result %s should contain %s ", *lsResult0, uploadPath)
	}
	log.Printf("found file %s", *lsResult0)

	if !strings.Contains(*lsResult1, uploadPath) {
		t.Fatalf("Result %s should contain %s ", *lsResult1, uploadPath)
	}
}
