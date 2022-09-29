# Project 1: Distributed File System

Please see: https://www.cs.usfca.edu/~mmalensek/cs677/assignments/project-1.html

You can structure your project however you'd like, but be sure to include instructions here on how to build and run it.


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
`./dfs -type=controller -host=localhost -port=12999 -id=controller1`

Run a storage node with 1GB in size
`./dfs -type=storage -host=localhost -port=12998 -id=storage0 -storage_size=1000000