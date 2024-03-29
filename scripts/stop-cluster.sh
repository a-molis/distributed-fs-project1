#!/usr/bin/env bash

script_dir="$(cd "$(dirname "$0")" && pwd)"
log_dir="${script_dir}/logs"

source "${script_dir}/nodes.sh"

echo "Stopping controller..."
ssh "${controller}" 'pkill -u $(whoami) dfs'

echo "Stopping Storage Nodes..."
for node in ${nodes[@]}; do
    echo "${node}"
    ssh "${node}" 'pkill -u $(whoami) dfs'
done

echo "Done!"
