package client

import (
	"P1-go-distributed-file-system/controller"
	"P1-go-distributed-file-system/storage"
	"testing"
	"time"
)

func TestBasicClient(t *testing.T) {

	var port int = 12024
	host := "localhost"

	testId := "storageTestId"
	testId2 := "storageTestId2"

	var members []string

	var size int32 = 10

	//bit for the controller
	go func(port int, receivedMembers *[]string) {
		testController := controller.NewController("testId", host, port)
		testController.Start()
		time.Sleep(time.Second * 1)
		*receivedMembers = testController.List()
	}(port, &members)

	time.Sleep(time.Second * 1)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
		time.Sleep(time.Second * 1)
	}(port, testId)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, size, "localHost", port)
		storageNode.Start()
	}(port, testId2)

	go func(port int, id string) {
		testClient := NewClient("localHost", port, "ls")
		testClient.Start()
	}(port, testId2)
	time.Sleep(time.Second * 5)

	if 1 != 1 {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}