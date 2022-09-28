package controller

import (
	"P1-go-distributed-file-system/config"
	"P1-go-distributed-file-system/storage"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	testId := "storageTestId"
	var storagePort0 int32 = 12041

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
		controller := NewController("testId", testConfig.ControllerHost, testConfig.ControllerPort, testConfig)
		controller.Start()
		time.Sleep(time.Second * 3)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig)
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
		controller := NewController("testId", testConfig.ControllerHost, testConfig.ControllerPort, testConfig)
		controller.Start()
		time.Sleep(time.Second * 15)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig)
		storageNode.Start()
		storageNode.Shutdown()
	}(testId)

	time.Sleep(time.Second * 20)

	if len(members) != 0 {
		t.Fatalf("the registered node did not switch to inactive")
	}

	return
}

func TestHeartBeat(t *testing.T) {
	testId := "storageTestId"
	var storagePort0 int32 = 12037

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
		controller := NewController("testId", testConfig.ControllerHost, testConfig.ControllerPort, testConfig)
		controller.Start()
		time.Sleep(time.Second * 15)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func( id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig)
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
		controller := NewController("testId", testConfig.ControllerHost, testConfig.ControllerPort, testConfig)
		controller.Start()
		time.Sleep(time.Second * 20)
		*receivedMembers = controller.List()
	}(&members)

	time.Sleep(time.Second * 1)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort0, testConfig)
		storageNode.Start()
		time.Sleep(time.Second * 1)
		storageNode.Shutdown()
	}(testId)

	go func(id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", storagePort1, testConfig)
		storageNode.Start()
	}(testId2)

	time.Sleep(time.Second * 25)

	if len(members) == 0 || members[0] != testId2 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}
