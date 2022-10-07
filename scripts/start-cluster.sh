#!/usr/bin/env bash

script_dir="$(cd "$(dirname "$0")" && pwd)"
log_dir="${script_dir}/logs"

source "${script_dir}/nodes.sh"

echo "script_dir :" ${script_dir}
echo "log_dir :" ${log_dir}
echo "file_system_dir :" ${file_system_dir}
echo "controller_port :" ${storage_node_port}
echo "storage_node_port :" ${storage_node_port}


echo "Building project..."
cd ${file_system_dir} || exit 1
go clean
go build || exit 1
cd ${script_dir} || exit 1
echo "Done!"

echo "Creating log directory: ${log_dir}"
mkdir -pv "${log_dir}"

echo "Starting Controller..."
ssh "${controller}" "cd ${file_system_dir} && pwd && ${file_system_dir}/dfs -type=controller -host=${HOSTNAME} -port=${controller_port} -id=${controller} &> "${log_dir}/controller.log" &"


echo "Starting Storage Nodes..."
for node in ${nodes[@]}; do
    echo "${node}"
    ssh "${node}" "cd ${file_system_dir} && "${file_system_dir}/dfs -type=storage -host=${HOSTNAME} -port=${storage_node_port} -id=${node} -storage_size=1000000 -local_path=${data_path}/sn" &> "${log_dir}/${node}.log" &"
    storage_node_port="$((storage_node_port + 1))"
done

echo "Startup complete!"
