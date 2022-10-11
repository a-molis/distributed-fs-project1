#!/usr/bin/env bash

node1="orion05"
node2="orion08"

echo "Stopping Storage Nodes..."
echo "${node1}"
ssh "${node1}" 'pkill -u $(whoami) dfs'
echo "${node2}"
ssh "${node2}" 'pkill -u $(whoami) dfs'

echo "Done!"
