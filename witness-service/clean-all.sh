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
    docker-compose --version > /dev/null 2>&1
    if [ $? -eq 127 ];then
        echo -e "$RED docker-compose command not found $NC"
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

function cleanAll() {
    echo -e "$YELLOW Starting clean all containers... $NC"
    pushd $IOTEX_WITNESS/etc
    docker-compose rm -s -f -v
    popd
    echo -e "${YELLOW} Done. ${NC}"

    echo -e "${YELLOW} Starting delete all files... ${NC}"
    if [ "${IOTEX_WITNESS}X" = "X" ] || [ "${IOTEX_WITNESS}X" = "/X" ];then
        echo -e "${RED} \$IOTEX_WITNESS: ${IOTEX_WITNESS} is wrong. ${NC}"
        ## For safe.
        return
    fi

    sudo rm -rf $IOTEX_WITNESS

    echo -e "${YELLOW} Done. ${NC}"
}

function main() {
    checkDockerPermissions
    checkDockerCompose

    determinIotexWitness
    confirmEnvironmentVariable

    cleanAll
}

main $@
