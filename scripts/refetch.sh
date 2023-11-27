chain=$1
shift
IOTEX_WITNESS=/path/to/witness_config/
set -o allexport
source $IOTEX_WITNESS/etc/.env
set +o allexport
cd /path/to/ioTube_home/
go run cmd/witness/main.go -config $IOTEX_WITNESS/etc/witness-config-${chain}.yaml -secret $IOTEX_WITNESS/etc/witness-config-${chain}.secret.yaml -blocks $@
