package client

import (
	"dfs/config"
	"dfs/controller"
	"dfs/storage"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		testClient := NewClient(testConfig, "ls", "/")
		_, clientLsError = testClient.Run()
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
	err = os.Remove(testConfig.ControllerPath)
	if err != nil {
		log.Println("Unable to remove controller backup ", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)
	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	testFile := "foo.txt"
	lsPath := "/this/test/path"
	uploadPath := filepath.Join(lsPath, testFile)
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
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 3)
	var clientLsError0 error = nil
	var lsResult0 *string
	go func() {
		testClient := NewClient(testConfig, "ls", lsPath)
		result, err := testClient.Run()
		clientLsError0 = err
		lsResult0 = result
	}()

	time.Sleep(time.Second * 10)

	if clientPutError != nil || clientLsError0 != nil {
		t.Fatalf("client test failed %s %s ", clientPutError, clientLsError0)
	}
	if !strings.Contains(*lsResult0, testFile) {
		t.Fatalf("Result %s should contain %s ", *lsResult0, testFile)
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
				log.Println(fullPath)
				dirName := fmt.Sprintf("./sn%d", j)
				err := os.RemoveAll(dirName)
				if err != nil {
					log.Printf("Error removing %s %s", dirName, err)
				}
				log.Println("removing dir", dirName)
				pwd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				log.Println(pwd)
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

	err = os.Remove(testConfig.ControllerPath)
	if err != nil {
		log.Println("Unable to remove controller backup ", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)

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
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 20)

	var clientGetError error = nil
	go func() {
		testClient := NewClient(testConfig, "get", remotePath, savePath)
		_, clientGetError = testClient.Run()
	}()

	time.Sleep(time.Second * 10)

	// TODO complete test to validate file is saved
	if clientGetError != nil || clientPutError != nil {
		t.Fatalf("client test failed %s, %s", clientGetError, clientPutError)
	}
	return
}

func TestClientRmSimple(t *testing.T) {
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
	controllerPort := 12067
	var storagePort0 int32 = 12068
	var storagePort1 int32 = 12069
	var storagePort2 int32 = 12070
	var storagePort3 int32 = 12071
	var storagePort4 int32 = 12072
	var size int64 = 1000000
	var chunkSize int64 = 5000000

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
			log.Println("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)
	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	testFile := "bar.txt"
	lsPath := "/this/test/path"
	remotePath := filepath.Join(lsPath, testFile)
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
		_, clientPutError = testClient.Run()
	}()
	time.Sleep(time.Second * 6)

	var clientLsError0 error = nil
	var lsResult0 *string
	go func() {
		testClient := NewClient(testConfig, "ls", lsPath)
		result, err := testClient.Run()
		clientLsError0 = err
		lsResult0 = result
	}()

	time.Sleep(time.Second * 3)

	var clientRmError error = nil
	go func() {
		testClient := NewClient(testConfig, "rm", remotePath)
		_, clientRmError = testClient.Run()
	}()

	time.Sleep(time.Second * 3)

	var clientLsError1 error = nil
	var lsResult1 *string
	go func() {
		testClient := NewClient(testConfig, "ls", lsPath)
		result, err := testClient.Run()
		clientLsError1 = err
		lsResult1 = result
	}()

	time.Sleep(time.Second * 10)

	if !strings.Contains(*lsResult0, testFile) {
		t.Fatalf("Result %s should contain %s ", *lsResult0, testFile)
	}

	if strings.Contains(*lsResult1, testFile) {
		t.Fatalf("Result %s should not contain %s ", *lsResult1, testFile)
	}

	// TODO complete test to validate file is saved
	if clientPutError != nil || clientRmError != nil || clientLsError0 != nil || clientLsError1 != nil {
		t.Fatalf("client test failed %s, %s, %s, %s",
			clientRmError, clientPutError, clientLsError0, clientLsError1)
	}
	return
}

func TestClientLs(t *testing.T) {
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
	controllerPort := 12082
	var storagePort0 int32 = 12081
	var storagePort1 int32 = 12080
	var storagePort2 int32 = 12077
	var storagePort3 int32 = 12079

	var size int64 = 1000000
	var chunkSize int64 = 5000000

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

	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	testFile := "foo.txt"
	dog := "dog"
	testPath := "/this/test"
	lsPath := filepath.Join(testPath, dog)
	uploadPath := filepath.Join(lsPath, testFile)
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
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 3)
	var clientLsError0 error = nil
	var lsResult0 *string
	go func() {
		testClient := NewClient(testConfig, "ls", lsPath)
		result, err := testClient.Run()
		clientLsError0 = err
		lsResult0 = result
	}()

	var clientLsError1 error = nil
	var lsResult1 *string
	go func() {
		testClient := NewClient(testConfig, "ls", testPath)
		result, err := testClient.Run()
		clientLsError1 = err
		lsResult1 = result
	}()

	time.Sleep(time.Second * 10)

	if clientPutError != nil || clientLsError0 != nil || clientLsError1 != nil {
		t.Fatalf("client test failed %s %s %s ",
			clientPutError, clientLsError0, clientLsError1)
	}
	if !strings.Contains(*lsResult0, testFile) {
		t.Fatalf("Result %s should contain %s ", *lsResult0, testFile)
	}

	if !strings.Contains(*lsResult1, dog) {
		t.Fatalf("Result %s should contain %s ", *lsResult1, dog)
	}
}

func TestClientLsFilePath(t *testing.T) {
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
	controllerPort := 12062
	var storagePort0 int32 = 12063
	var storagePort1 int32 = 12064
	var storagePort2 int32 = 12065
	var storagePort3 int32 = 12066

	var size int64 = 1000000
	var chunkSize int64 = 5000000

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
			log.Println("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)

	testConfig.ChunkSize = chunkSize
	testConfig.ControllerHost = controllerHost
	testConfig.ControllerPort = controllerPort

	testFile := "foo.txt"
	cat := "cat"
	testPath := "/this/test"
	lsPath := filepath.Join(testPath, cat)
	uploadPath := filepath.Join(lsPath, testFile)
	badPath := filepath.Join(lsPath, "fakeFile.txt")
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
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 3)
	var clientLsError0 error = nil
	var lsResult0 *string
	go func() {
		testClient := NewClient(testConfig, "ls", badPath)
		result, err := testClient.Run()
		clientLsError0 = err
		lsResult0 = result
	}()

	var clientLsError1 error = nil
	var lsResult1 *string
	go func() {
		testClient := NewClient(testConfig, "ls", uploadPath)
		result, err := testClient.Run()
		clientLsError1 = err
		lsResult1 = result
	}()

	time.Sleep(time.Second * 10)

	if clientPutError != nil || clientLsError0 != nil || clientLsError1 != nil {
		t.Fatalf("client test failed %s %s %s ",
			clientPutError, clientLsError0, clientLsError1)
	}

	if clientPutError != nil || clientLsError1 != nil {
		t.Fatalf("client test failed %s %s ",
			clientPutError, clientLsError1)
	}
	notFound := "No such file or directory"
	if !strings.Contains(*lsResult0, notFound) {
		t.Fatalf("Result %s should contain %s ", *lsResult0, "not found")
	}
	log.Printf("found file %s", *lsResult0)

	if !strings.Contains(*lsResult1, uploadPath) {
		t.Fatalf("Result %s should contain %s ", *lsResult1, uploadPath)
	}
}

func TestClientStats(t *testing.T) {
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
				log.Println(fullPath)
				dirName := fmt.Sprintf("./sn%d", j)
				err := os.RemoveAll(dirName)
				if err != nil {
					log.Printf("Error removing %s %s", dirName, err)
				}
				log.Println("removing dir", dirName)
				pwd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				log.Println(pwd)
			}
		}
	}()

	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12083
	var storagePort0 int32 = 12084
	var storagePort1 int32 = 12085
	var storagePort2 int32 = 12086
	var storagePort3 int32 = 12087
	var storagePort4 int32 = 12088
	var size int64 = 1000000
	var chunkSize int64 = 5000000

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
			log.Println("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)

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
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 15)

	var clientGetError error = nil
	go func() {
		testClient := NewClient(testConfig, "get", remotePath, savePath)
		_, clientGetError = testClient.Run()
	}()

	time.Sleep(time.Second * 10)

	var statsResult0 *string
	var clientStatsError error = nil
	go func() {
		testClient := NewClient(testConfig, "stats")
		result, err := testClient.Run()
		statsResult0 = result
		clientStatsError = err
		if err != nil {
			return
		}
	}()

	time.Sleep(time.Second * 2)

	if clientGetError != nil || clientPutError != nil || clientStatsError != nil {
		t.Fatalf("client test failed %s, %s, %s\n", clientGetError, clientPutError, clientStatsError)
	}

	if !strings.Contains(*statsResult0, testStorageNode4) {
		t.Fatalf("Result %s should contain %s ", *statsResult0, testStorageNode4)
	}
	return
}

func TestClientDownloadFailingNodes(t *testing.T) {
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
				log.Println(fullPath)
				dirName := fmt.Sprintf("./sn%d", j)
				err := os.RemoveAll(dirName)
				if err != nil {
					log.Printf("Error removing %s %s", dirName, err)
				}
				log.Println("removing dir", dirName)
				pwd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				log.Println(pwd)
			}
		}
	}()

	controllerHost := "localhost"
	storageHost := "localhost"
	controllerPort := 12090
	var storagePort0 int32 = 12091
	var storagePort1 int32 = 12092
	var storagePort2 int32 = 12093
	var storagePort3 int32 = 12094
	var storagePort4 int32 = 12095
	var size int64 = 1000000
	var chunkSize int64 = 5000000

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
			log.Println("Unable to remove controller backup")
		}
	}(testConfig.ControllerPath)

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

	storageNode1 := storage.NewStorageNode(testStorageNode3, size, storageHost, storagePort3, testConfig, savePathStorageNode3)
	go storageNode1.Start()
	go func (storageNode *storage.StorageNode) {
		time.Sleep(time.Second * 20)
		storageNode.Shutdown()
	}(storageNode1)


	storageNode2 := storage.NewStorageNode(testStorageNode4, size, storageHost, storagePort4, testConfig, savePathStorageNode4)
	go storageNode2.Start()
	go func (storageNode *storage.StorageNode) {
		time.Sleep(time.Second * 20)
		storageNode.Shutdown()
	}(storageNode2)

	time.Sleep(time.Second * 3)
	var clientPutError error = nil
	go func() {
		testClient := NewClient(testConfig, "put", remotePath, localPath)
		_, clientPutError = testClient.Run()
	}()

	time.Sleep(time.Second * 40)

	var clientGetError error = nil
	go func() {
		testClient := NewClient(testConfig, "get", remotePath, savePath)
		_, clientGetError = testClient.Run()
	}()

	time.Sleep(time.Second * 10)

	// TODO complete test to validate file is saved
	if clientGetError != nil || clientPutError != nil {
		t.Fatalf("client test failed %s, %s", clientGetError, clientPutError)
	}
	return
}
