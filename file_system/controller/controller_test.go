package controller

import (
	"P1-go-distributed-file-system/storage"
	"testing"
	"time"
)

func TestController(t *testing.T) {

	var port int = 12020
	host := "localhost"

	testId := "storageTestId"

	var members []string

	var size int32 = 10

	//bit for the controller
	go func(port int, receivedMembers *[]string) {
		controller := NewController("testId", host, port)
		controller.Start()
		time.Sleep(time.Second * 3)
		*receivedMembers = controller.List()
	}(port, &members)

	time.Sleep(time.Second * 1)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
	}(port, testId)

	time.Sleep(time.Second * 6)

	if len(members) == 0 || members[0] != testId {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestHeartBeatInactive(t *testing.T) {

	var port int = 12021
	host := "localhost"

	testId := "storageTestId"

	var members []string

	var size int32 = 10

	//bit for the controller
	go func(port int, receivedMembers *[]string) {
		controller := NewController("testId", host, port)
		controller.Start()
		time.Sleep(time.Second * 15)
		*receivedMembers = controller.List()
	}(port, &members)

	time.Sleep(time.Second * 1)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
		storageNode.Shutdown()
	}(port, testId)

	time.Sleep(time.Second * 20)

	if len(members) != 0 {
		t.Fatalf("the registered node did not switch to inactive")
	}

	return
}

func TestHeartBeat(t *testing.T) {

	var port int = 12022
	host := "localhost"

	testId := "storageTestId"

	var members []string

	var size int32 = 10

	//bit for the controller
	go func(port int, receivedMembers *[]string) {
		controller := NewController("testId", host, port)
		controller.Start()
		time.Sleep(time.Second * 15)
		*receivedMembers = controller.List()
	}(port, &members)

	time.Sleep(time.Second * 1)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
	}(port, testId)

	time.Sleep(time.Second * 20)

	if len(members) == 0 || members[0] != testId {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}

func TestHeartBeatMultiNode(t *testing.T) {

	var port int = 12023
	host := "localhost"

	testId := "storageTestId"
	testId2 := "storageTestId2"

	var members []string

	var size int32 = 10

	//bit for the controller
	go func(port int, receivedMembers *[]string) {
		controller := NewController("testId", host, port)
		controller.Start()
		time.Sleep(time.Second * 20)
		*receivedMembers = controller.List()
	}(port, &members)

	time.Sleep(time.Second * 1)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
		time.Sleep(time.Second * 1)
		storageNode.Shutdown()
	}(port, testId)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
	}(port, testId2)

	time.Sleep(time.Second * 25)

	if len(members) == 0 || members[0] != testId2 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}
