#!/bin/bash

# Colour codes
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

defaultdatadir="$HOME/iotex-witness"
CURL="curl -Ss"
DB_ROOT_PASSWORD="kdfjjrU64fjK58H"
PROJECT_ABS_DIR=$(cd "$(dirname "$0")";pwd)

pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null
}

function checkDockerPermissions() {
    docker ps > /dev/null
    if [ $? = 1 ];then
        echo -e "your $RED [$USER] $NC not privilege docker"
        echo -e "please run $RED [sudo bash] $NC first"
        echo -e "Or docker not install "
        exit 1
    fi
}

function checkDockerCompose() {
    docker compose --version > /dev/null 2>&1
    if [ $? -eq 127 ];then
        echo -e "$RED docker compose command not found $NC"
        echo -e "Please install it first"
        exit 1
    fi
}

function determinIotexWitness() {
    if [[ ! $IOTEX_WITNESS ]];then
        ##Input Data Dir
        echo "Input IOTEX_WITNESS directory, Service witness will copy config file into this dir."
        echo "The current user of the input directory must have write permission!!!"
        echo -e "${RED} Input your directory \$IOTEX_WITNESS !!! ${NC}"
     
        #while True: do
        read -p "Input your \$IOTEX_WITNESS [e.g., $defaultdatadir]: " inputdir
        IOTEX_WITNESS=${inputdir:-"$defaultdatadir"}
    fi
}

function confirmEnvironmentVariable() {
    echo -e "IOTEX_WITNESS directory: ${RED} ${IOTEX_WITNESS} ${NC}, Service witness will copy config file into this dir."
}

function copyFile() {
    srcFile=$1
    tgtFile=$2
    if [[ ! -f ${IOTEX_WITNESS}/etc/$tgtFile || $# -ge 3 && $3 == 1 ]]; then
        echo -e "copy file ${srcFile} to ${tgtFile}"
        cp -f $PROJECT_ABS_DIR/$srcFile ${IOTEX_WITNESS}/etc/$tgtFile
         if [ $? -ne 0 ];then
             echo "Get config error"
             exit 2
         fi
    else
       echo "skip copy file ${srcFile} to ${tgtFile}"
    fi
}

function downloadConfigFile() {
    copyFile "docker-compose-witness.yml" "docker-compose.yml" 1
    copyFile "witness-config-iotex.yaml" "witness-config-iotex.yaml" 1
    copyFile "witness-config-iotex.secret.yaml" "witness-config-iotex.secret.yaml" 0
    copyFile "witness-config-ethereum.yaml" "witness-config-ethereum.yaml" 1
    copyFile "witness-config-ethereum.secret.yaml" "witness-config-ethereum.secret.yaml" 0
    #copyFile "witness-config-heco.yaml" "witness-config-heco.yaml" 1
    #copyFile "witness-config-heco.secret.yaml" "witness-config-heco.secret.yaml" 0
    #copyFile "witness-config-polis.yaml" "witness-config-polis.yaml" 1
    #copyFile "witness-config-polis.secret.yaml" "witness-config-polis.secret.yaml" 0
    copyFile "witness-config-bsc.yaml" "witness-config-bsc.yaml" 1
    copyFile "witness-config-bsc.secret.yaml" "witness-config-bsc.secret.yaml" 0
    copyFile "witness-config-matic.yaml" "witness-config-matic.yaml" 1
    copyFile "witness-config-matic.secret.yaml" "witness-config-matic.secret.yaml" 0
    copyFile "witness-config-bitcoin-testnet.yaml" "witness-config-bitcoin-testnet.yaml" 1
    copyFile "witness-config-bitcoin-testnet.secret.yaml" "witness-config-bitcoin-testnet.secret.yaml" 0
    envFile=${IOTEX_WITNESS}/etc/.env
    if [[ ! -f ${envFile} ]]; then
        touch ${envFile}
    fi
    if grep -q "^IOTEX_WITNESS=" ${envFile}; then
        sed "s/^IOTEX_WITNESS=.*//" ${envFile} > ${envFile}.tmp
        mv ${envFile}.tmp ${envFile}
    fi
    echo "" >> $envFile
    echo "IOTEX_WITNESS=$IOTEX_WITNESS" >> ${envFile}
    if grep -q "/^DB_ROOT_PASSWORD=" ${envFile}; then
        sed "/^DB_ROOT_PASSWORD=.*//" ${envFile} > ${envFile}.tmp
        mv ${envFile}.tmp $envFile
    fi
    echo "DB_ROOT_PASSWORD=$DB_ROOT_PASSWORD" >> ${envFile}
    sed "/^$/d" ${envFile} > ${envFile}.tmp
    mv ${envFile}.tmp $envFile
    cp -f $PROJECT_ABS_DIR/crontab ${IOTEX_WITNESS}/etc/crontab
    cp -f $PROJECT_ABS_DIR/backup_witness ${IOTEX_WITNESS}/etc/backup
}

function makeWorkspace() {
    mkdir -p ${IOTEX_WITNESS}
    mkdir -p ${IOTEX_WITNESS}/etc
    mkdir -p ${IOTEX_WITNESS}/data/mysql
    mkdir -p ${IOTEX_WITNESS}/backup
    downloadConfigFile
}

function exportAll() {
    export IOTEX_WITNESS DB_ROOT_PASSWORD PROJECT_ABS_DIR
}

function grantPrivileges() {
    if [[ ! -f $IOTEX_WITNESS/data/mysql/.inited ]];then
        echo -e "$YELLOW Starting database...$NC"
        # maxRetryTime * sleeptime = timeout
        retryTimes=0
        maxRetryTime=10
        pushd $IOTEX_WITNESS/etc
        docker compose up -d database

        echo -e "$YELLOW Waiting for the mysqld daemon in the witness-db container to successful... $NC"
        while true;do
            if [ $retryTimes -gt $maxRetryTime ];then
                echo -e "$RED Start mysql server container faild. $NC"
                echo -e "$RED Please check its logs by command \"docker logs iotex-db\" $NC"
                exit 1
            fi
            docker exec witness-db mysql -uroot -p${DB_ROOT_PASSWORD} -e "\q" > /dev/null 2>&1
            if [ $? -eq 0 ];then
                break
            fi
            retryTimes=$((retryTimes+1))
            sleep 4
        done
        popd
        echo -e "$YELLOW Success! $NC"
        docker exec witness-db mysql -uroot -p${DB_ROOT_PASSWORD} -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%'"  > /dev/null 2>&1
        docker exec witness-db mysql -uroot -p${DB_ROOT_PASSWORD} -e "CREATE DATABASE witness"  > /dev/null 2>&1
        $WHITE_LINE
        touch $IOTEX_WITNESS/data/mysql/.inited
    fi
 }

function buildService() {
    pushd $PROJECT_ABS_DIR
    docker build . -f Dockerfile.witness -t witness:latest || exit 2
}

function startup() {
    echo -e "$YELLOW Start witness and it's database. $NC"
    pushd $IOTEX_WITNESS/etc
    docker compose up -d
    docker compose restart
    docker system prune -a
    if [ $? -eq 0 ];then
        echo -e "${YELLOW} Service on. ${NC}"
    fi
    popd
}

function cleanAll() {
    echo -e "$YELLOW Starting clean all containers... $NC"
    pushd $IOTEX_WITNESS/etc
    docker compose rm -s -f -v
    popd
    echo -e "${YELLOW} Done. ${NC}"

    echo -e "${YELLOW} Starting delete all files... ${NC}"
    if [ "${IOTEX_WITNESS}X" = "X" ] || [ "${IOTEX_WITNESS}X" = "/X" ];then
        echo -e "${RED} \$IOTEX_WITNESS: ${IOTEX_WITNESS} is wrong. ${NC}"
        ## For safe.
        return
    fi

    $RM $IOTEX_WITNESS

    echo -e "${YELLOW} Done. ${NC}"
}

function main() {
    checkDockerPermissions
    checkDockerCompose

    determinIotexWitness
    confirmEnvironmentVariable

    makeWorkspace

    exportAll

    buildService
    grantPrivileges

    startup
}

main $@
