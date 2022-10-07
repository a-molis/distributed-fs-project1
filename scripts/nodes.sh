controller="orion03"

# Node list. You can comment out nodes that you don't want to use with '#'
nodes=(
"orion01"
"orion02"
"orion03"
"orion04"
"orion05"
"orion06"
"orion07"
"orion08"
"orion09"
"orion10"
"orion11"
"orion12"
)

package_dir="${HOME}/Classes/cs677/project1/P1-go-distributed-file-system"

file_system_dir="${package_dir}/file_system"
controller_port=12500
storage_node_port=12501
data_path="/bigdata/$(whoami)"
