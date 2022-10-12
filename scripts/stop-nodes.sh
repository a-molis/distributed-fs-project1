#!/usr/bin/env bash

node1=$1
node2=$2

echo "Stopping Storage Nodes..."
echo "${node1}"
ssh "${node1}" 'pkill -u $(whoami) dfs'
echo "${node2}"
ssh "${node2}" 'pkill -u $(whoami) dfs'

echo "Done!"
