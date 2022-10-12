#!/usr/bin/env bash

script_dir="$(cd "$(dirname "$0")" && pwd)"
log_dir="${script_dir}/logs"

source "${script_dir}/nodes.sh"

echo "script_dir :" ${script_dir}
echo "log_dir :" ${log_dir}
echo "file_system_dir :" ${file_system_dir}
echo "controller_port :" ${storage_node_port}
echo "storage_node_port :" ${storage_node_port}
echo "Starting Controller..."
ssh "${controller}" "cd ${file_system_dir} && pwd && ${file_system_dir}/dfs -type=controller -host=${controller}.cs.usfca.edu -port=${controller_port} -id=${controller} &> "${log_dir}/controller.log" &"
