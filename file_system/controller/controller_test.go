package controller

import (
	"dfs/config"
	"dfs/storage"
	"os"
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
	defer os.Remove("fmdt")

	var file string

	testConfig, err := config.ConfigFromPath("../config.json")
	if err != nil {
		t.Errorf("Unable to open config")
		return
	}
	testConfig.ControllerHost = "localhost"
	testConfig.ControllerPort = 12060

	//bit for the controller
	go func(file *string) {
		controller := NewController("testId", testConfig)
		go controller.Start()
		controller.fileMetadata.Upload("/foo/something/file.txt")
		controller.fileMetadata.Upload("/foo/something2/file2.txt")
		controller.SaveFileMetadata()
		controller.shutdown()
		testConfig.ControllerPort = 12061
		controller = NewController("testId", testConfig)
		go controller.Start()
		controller.LoadFileMetadata()
		*file = controller.fileMetadata.Ls("/foo/something/")
	}(&file)

	time.Sleep(time.Second * 6)

	if file != "file.txt" {
		t.Fatalf("loaded file metadata is incorrect, %s", file)
	}

	return
}
