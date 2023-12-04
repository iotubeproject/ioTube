chain=$1
shift
declare -A chainmap
chainmap["iotex"]="9081"
chainmap["ethereum"]="9083"
chainmap["matic"]="9085"
chainmap["bsc"]="9087"
echo "refetch from port:${port}"
curl -X POST localhost:$port/fetch -d '{"heights": '$@'}'
echo "...done"