if [ $# -ne 2 ]; then
    echo $0 chain id
    exit
fi
platform=$1
id=$2
docker exec witness-db mysql -uroot -pkdfjjrU64fjK58H -e "DELETE FROM witness.${platform}_transfers WHERE id='$id';"