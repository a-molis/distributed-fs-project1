package controller

import (
	"P1-go-distributed-file-system/storage"
	"testing"
	"time"
)

func TestController(t *testing.T) {

	var port int = 12000

	testId := "storageTestId"

	var members []string

	//bit for the controller
	go func(port int, receivedMembers *[]string) {
		controller := NewController("testId", port)
		controller.Start()
		time.Sleep(time.Second * 3)
		*receivedMembers = controller.List()
	}(port, &members)

	time.Sleep(time.Second * 1)

	go func(port int, id string) {
		storageNode := storage.NewStorageNode(id, "localHost", port)
		storageNode.Start()
	}(port, testId)

	time.Sleep(time.Second * 6)

	if len(members) == 0 || members[0] != testId {
		t.Fatalf("the registered node id doesnt match %s", members[0])
	}

	return
}
