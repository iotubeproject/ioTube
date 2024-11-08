chain=$1
shift
declare -A chainmap
chainmap["iotex"]="9281"
chainmap["ethereum"]="9283"
chainmap["matic"]="9285"
chainmap["bsc"]="9287"
port=${chainmap[$chain]}
echo "refetch from port:${port}"
curl -X POST http://localhost:$port/fetch -d '{"heights": "'$@'"}'
echo "...done"
