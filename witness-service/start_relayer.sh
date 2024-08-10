#!/bin/bash

# Colour codes
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

defaultdatadir="$HOME/iotex-relayer"
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

function determinIotexRelayer() {
    if [[ ! $IOTEX_RELAYER ]];then
        ##Input Data Dir
        echo "Input IOTEX_RELAYER directory, Service relayer will copy config file into this dir."
        echo "The current user of the input directory must have write permission!!!"
        echo -e "${RED} Input your directory \$IOTEX_RELAYER !!! ${NC}"
     
        #while True: do
        read -p "Input your \$IOTEX_RELAYER [e.g., $defaultdatadir]: " inputdir
        IOTEX_RELAYER=${inputdir:-"$defaultdatadir"}
    fi
}

function confirmEnvironmentVariable() {
    echo -e "IOTEX_RELAYER directory: ${RED} ${IOTEX_RELAYER} ${NC}, Service relayer will copy config file into this dir."
}

function downloadConfigFile() {
    if [[ ! -f ${IOTEX_RELAYER}/etc/docker-compose.yml ]];then
        cp -f $PROJECT_ABS_DIR/docker-compose-relayer.yml ${IOTEX_RELAYER}/etc/docker-compose.yml
        if [ $? -ne 0 ];then
            echo "Get docker-compose config error"
            exit 2
        fi
    fi
    
    if [[ ! -f ${IOTEX_RELAYER}/etc/relayer-config-iotex.yaml ]];then
        cp -f $PROJECT_ABS_DIR/configs/relayer-config-iotex.yaml ${IOTEX_RELAYER}/etc/relayer-config-iotex.yaml
        if [ $? -ne 0 ];then
            echo "Get config error"
            exit 2
        fi
    fi
    if [[ ! -f ${IOTEX_RELAYER}/etc/relayer-config-ethereum.yaml ]];then
        cp -f $PROJECT_ABS_DIR/configs/relayer-config-ethereum.yaml ${IOTEX_RELAYER}/etc/relayer-config-ethereum.yaml
        if [ $? -ne 0 ];then
            echo "Get config error"
            exit 2
        fi
    fi
#    if [[ ! -f ${IOTEX_RELAYER}/etc/relayer-config-heco.yaml ]];then
#        cp -f $PROJECT_ABS_DIR/configs/relayer-config-heco.yaml ${IOTEX_RELAYER}/etc/relayer-config-heco.yaml
#        if [ $? -ne 0 ];then
#            echo "Get config error"
#            exit 2
#        fi
#    fi
#    if [[ ! -f ${IOTEX_RELAYER}/etc/relayer-config-polis.yaml ]];then
#        cp -f $PROJECT_ABS_DIR/configs/relayer-config-polis.yaml ${IOTEX_RELAYER}/etc/relayer-config-polis.yaml
#        if [ $? -ne 0 ];then
#            echo "Get config error"
#            exit 2
#        fi
#    fi
    if [[ ! -f ${IOTEX_RELAYER}/etc/relayer-config-bsc.yaml ]];then
        cp -f $PROJECT_ABS_DIR/configs/relayer-config-bsc.yaml ${IOTEX_RELAYER}/etc/relayer-config-bsc.yaml
        if [ $? -ne 0 ];then
            echo "Get config error"
            exit 2
        fi
    fi
    if [[ ! -f ${IOTEX_RELAYER}/etc/relayer-config-matic.yaml ]];then
        cp -f $PROJECT_ABS_DIR/configs/relayer-config-matic.yaml ${IOTEX_RELAYER}/etc/relayer-config-matic.yaml
        if [ $? -ne 0 ];then
            echo "Get config error"
            exit 2
        fi
    fi
    [[ -f ${IOTEX_RELAYER}/etc/.env ]] || (echo "IOTEX_RELAYER=$IOTEX_RELAYER" > ${IOTEX_RELAYER}/etc/.env;echo "DB_ROOT_PASSWORD=$DB_ROOT_PASSWORD" >> ${IOTEX_RELAYER}/etc/.env)
    cp -f $PROJECT_ABS_DIR/crontab ${IOTEX_RELAYER}/etc/crontab
    cp -f $PROJECT_ABS_DIR/backup_relayer ${IOTEX_RELAYER}/etc/backup
}

function makeWorkspace() {
    mkdir -p ${IOTEX_RELAYER}
    mkdir -p ${IOTEX_RELAYER}/etc
    mkdir -p ${IOTEX_RELAYER}/data/mysql
    mkdir -p ${IOTEX_RELAYER}/data/backup
    downloadConfigFile
}

function exportAll() {
    export IOTEX_RELAYER DB_ROOT_PASSWORD PROJECT_ABS_DIR
}

function grantPrivileges() {
    if [[ ! -f $IOTEX_RELAYER/data/mysql/.inited ]];then
        echo -e "$YELLOW Starting database...$NC"
        # maxRetryTime * sleeptime = timeout
        retryTimes=0
        maxRetryTime=10
        pushd $IOTEX_RELAYER/etc
        docker compose up -d database
    
        echo -e "$YELLOW Waiting for the mysqld daemon in the relayer-db container to successful... $NC"
        while true;do
            if [ $retryTimes -gt $maxRetryTime ];then
                echo -e "$RED Start mysql server container faild. $NC"
                echo -e "$RED Please check its logs by command \"docker logs iotex-db\" $NC"
                exit 1
            fi
            docker exec relayer-db mysql -uroot -p${DB_ROOT_PASSWORD} -e "\q" > /dev/null 2>&1
            if [ $? -eq 0 ];then
                break
            fi
            retryTimes=$((retryTimes+1))
            sleep 4
        done
        popd
        echo -e "$YELLOW Success! $NC"
        docker exec relayer-db mysql -uroot -p${DB_ROOT_PASSWORD} -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%'"  > /dev/null 2>&1
        $WHITE_LINE
        touch $IOTEX_RELAYER/data/mysql/.inited
    fi
 }

function buildService() {
    pushd $PROJECT_ABS_DIR
    docker build . -f Dockerfile.relayer -t relayer:latest || exit 2
}

function startup() {
    echo -e "$YELLOW Start relayer and it's database. $NC"
    pushd $IOTEX_RELAYER/etc
    docker compose up -d
    if [ $? -eq 0 ];then
        echo -e "${YELLOW} Server port on 7000 & 7001. ${NC}"
    fi
    popd
}

function cleanAll() {
    echo -e "$YELLOW Starting clean all containers... $NC"
    pushd $IOTEX_RELAYER/etc
    docker compose rm -s -f -v
    popd
    echo -e "${YELLOW} Done. ${NC}"

    echo -e "${YELLOW} Starting delete all files... ${NC}"
    if [ "${IOTEX_RELAYER}X" = "X" ] || [ "${IOTEX_RELAYER}X" = "/X" ];then
        echo -e "${RED} \$IOTEX_RELAYER: ${IOTEX_RELAYER} is wrong. ${NC}"
        ## For safe.
        return
    fi

    $RM $IOTEX_RELAYER

    echo -e "${YELLOW} Done. ${NC}"
}

function main() {
    checkDockerPermissions
    checkDockerCompose

    determinIotexRelayer
    confirmEnvironmentVariable

    makeWorkspace

    exportAll

    buildService
    grantPrivileges

    startup
}

main $@
