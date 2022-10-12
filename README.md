# Project 1: Distributed File System

Please see: https://www.cs.usfca.edu/~mmalensek/cs677/assignments/project-1.html

You can structure your project however you'd like, but be sure to include instructions here on how to build and run it.

## Running the cluster with start scripts

Starting the controller and storage nodes 
```bash 
cd scripts
sh start-cluster.sh
```

Running ls with client
```bash
./dfs -type=client -remote_path=/some/random/path/ -command=ls
```

Running put with client
```bash
/dfs -type=client -remote_path=/some/random/path/test_file_1.bin -local_path=/home/ykarakus/proj1-test-data/test_file_1.bin -command=put
```

Running get with client
```bash
./dfs -type=client -remote_path=/some/random/path/test_file_5.bin -local_path=/home/ykarakus/proj1-downloaded-data/test_file_1.bin -command=get
```

Running rm with client
```bash
./dfs -type=client -remote_path=/some/random/path/test_file_1.bin -command=rm
```

Running stats with client
```bash
./dfs -type=client -command=stats
```

## Manual Building and Running

Building the program
```bash
cd file_system
go build
```

Running ls with client
```bash
./dfs -type=client -local_path=dfs -remote_path=/foo/bar/dfs -command=ls
```

Run controller
```bash
./dfs -type=controller -host=localhost -port=12999 -id=controller1
```

Run a storage node with 1GB in size
```bash
./dfs -type=storage -host=localhost -port=12998 -id=storage0 -storage_size=1000000 -local_path=/path/to/data
```

## Design 

Controller
 
 - Controller has two data structures. 
 - Member table to keep track of all storage nodes. 
 - File_metadata to keep a record of the metadata for the files that have been uploaded. 
 - Heartbeats are handled by the controller and recorded in the member table.
 - Heartbeats are sent every 5 seconds. The threshold for deactivating the storage node is 10 seconds. 

Upload Process 

 - Client splits the file into chunks and sends the chunk metadata to the controller. 
 - Controller responds with a list of 3 storage nodes for each chunk. It also reserves the filename in file metadata and sets the status of the file as pending
 - Client connects to one of the available nodes for each chunk and starts a separate go routine for each chunk. 
 - Client first sends the chunk metadata and the replication nodes to the storage node, upon receiving an ack from the storage node the client writes the payload to the connection. 
 - Storage node reads the data and saves it and responds with an ack. Client exits after receiving this ack
 - The replication happens asynchronously where the storage node peels off itself from the message and connects to one of the remaining nodes and sends the replication message. 
 - After saving the data storage node sends a heartbeat with the file info to the controller. 
 - Controller saves file metadata to disk after updating it with the heartbeat info
 - When the controller receives the heartbeats for every chunk it sets the file status as complete 

## Testing

Running the tests 
```bash
cd file_system 
go test ./…
```
