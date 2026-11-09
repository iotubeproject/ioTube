#!/bin/bash

# Colour codes
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

defaultdatadir="$HOME/iotex-witness"
# SQL client inside witness-db: MariaDB 11 images ship only `mariadb`, MySQL only `mysql`
DB_CLI='if command -v mariadb >/dev/null 2>&1; then exec mariadb "$@"; else exec mysql "$@"; fi'
# witness image version seeded into a new .env (a release tag, never latest)
DEFAULT_WITNESS_TAG="7.2.0"
CURL="curl -Ss"
DB_ROOT_PASSWORD="${DB_ROOT_PASSWORD:-}"   # no hardcoded default; resolveDbPassword() fills this in
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
    git submodule update --remote --merge
    copyFile "docker-compose-witness.yml" "docker-compose.yml" 1
    copyFile "configs/witness-config-iotex-payload.yaml" "witness-config-iotex-payload.yaml" 1
    copyFile "configs/witness-config-iotex.secret.yaml" "witness-config-iotex.secret.yaml" 0
    copyFile "configs/witness-config-ethereum-payload.yaml" "witness-config-ethereum-payload.yaml" 1
    copyFile "configs/witness-config-ethereum.secret.yaml" "witness-config-ethereum.secret.yaml" 0
    copyFile "configs/witness-config-bsc-payload.yaml" "witness-config-bsc-payload.yaml" 1
    copyFile "configs/witness-config-bsc.secret.yaml" "witness-config-bsc.secret.yaml" 0
    copyFile "configs/witness-config-matic-payload.yaml" "witness-config-matic-payload.yaml" 1
    copyFile "configs/witness-config-matic.secret.yaml" "witness-config-matic.secret.yaml" 0
    copyFile "configs/witness-config-solana.yaml" "witness-config-solana.yaml" 1
    copyFile "configs/witness-config-solana.secret.yaml" "witness-config-solana.secret.yaml" 0
    copyFile "configs/witness-config-iotex-solana.yaml" "witness-config-iotex-solana.yaml" 1
    copyFile "configs/witness-config-iotex-solana.secret.yaml" "witness-config-iotex-solana.secret.yaml" 0
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
    if grep -q "^DB_ROOT_PASSWORD=" ${envFile}; then
        sed "s/^DB_ROOT_PASSWORD=.*//" ${envFile} > ${envFile}.tmp
        mv ${envFile}.tmp $envFile
    fi
    echo "DB_ROOT_PASSWORD=$DB_ROOT_PASSWORD" >> ${envFile}
    if ! grep -q "^WITNESS_TAG=" ${envFile}; then
        echo "WITNESS_TAG=${WITNESS_TAG:-$DEFAULT_WITNESS_TAG}" >> ${envFile}
    fi
    if grep -q "^WITNESS_DB_IMAGE=" ${envFile}; then
        sed "s/^WITNESS_DB_IMAGE=.*//" ${envFile} > ${envFile}.tmp
        mv ${envFile}.tmp $envFile
    fi
    echo "WITNESS_DB_IMAGE=$WITNESS_DB_IMAGE" >> ${envFile}
    sed "/^$/d" ${envFile} > ${envFile}.tmp
    mv ${envFile}.tmp $envFile
    cp -f $PROJECT_ABS_DIR/crontab ${IOTEX_WITNESS}/etc/crontab
    cp -f $PROJECT_ABS_DIR/backup_witness ${IOTEX_WITNESS}/etc/backup
}

function stopService() {
    if [[ -f ${IOTEX_WITNESS}/etc/docker-compose.yml ]]; then
        pushd $IOTEX_WITNESS/etc
        docker compose down
    fi
}

function makeWorkspace() {
    if ! mkdir -p "${IOTEX_WITNESS}" "${IOTEX_WITNESS}/etc" "${IOTEX_WITNESS}/data/mysql" "${IOTEX_WITNESS}/backup"; then
        echo -e "${RED} Cannot create the witness directories under ${IOTEX_WITNESS}. ${NC}"
        exit 1
    fi
    downloadConfigFile
}

function resolveDbPassword() {
    local envFile="${IOTEX_WITNESS}/etc/.env"
    if [[ -z "$DB_ROOT_PASSWORD" && -f "$envFile" ]]; then
        DB_ROOT_PASSWORD=$(grep -E '^DB_ROOT_PASSWORD=' "$envFile" | tail -1 | cut -d= -f2- | tr -d '\r')
    fi
    if [[ -z "$DB_ROOT_PASSWORD" ]]; then
        DB_ROOT_PASSWORD=$(openssl rand -base64 24 | tr -dc 'A-Za-z0-9' | head -c 24)
        echo -e "${YELLOW} Generated a new random DB password for the local database. ${NC}"
    fi
}

function exportAll() {
    # WITNESS_TAG is exported so docker compose uses the checked value: the shell
    # environment takes precedence over .env, whatever else .env contains.
    export IOTEX_WITNESS DB_ROOT_PASSWORD PROJECT_ABS_DIR WITNESS_TAG WITNESS_DB_IMAGE
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
            docker exec witness-db sh -c "$DB_CLI" -- -uroot "-p${DB_ROOT_PASSWORD}" -e "\q" > /dev/null 2>&1
            if [ $? -eq 0 ];then
                break
            fi
            retryTimes=$((retryTimes+1))
            sleep 4
        done
        popd
        echo -e "$YELLOW Success! $NC"
        docker exec witness-db sh -c "$DB_CLI" -- -uroot "-p${DB_ROOT_PASSWORD}" -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%'"  > /dev/null 2>&1
        docker exec witness-db sh -c "$DB_CLI" -- -uroot "-p${DB_ROOT_PASSWORD}" -e "CREATE DATABASE witness"  > /dev/null 2>&1
        $WHITE_LINE
        touch $IOTEX_WITNESS/data/mysql/.inited
    fi
 }

# A fixed image version: a release (7.2.0, 7.2.0-rc1) or a commit build
# (sha-<commit>), in Docker's tag syntax. Moving tags (latest, 7.2, master) are refused.
function isFixedTag() {
    [[ ${#1} -le 128 ]] && [[ $1 =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ || $1 =~ ^sha-[0-9a-f]{7,40}$ ]]
}

# Resolve WITNESS_TAG before anything is stopped, so a bad value never causes downtime.
function resolveWitnessTag() {
    local envFile="${IOTEX_WITNESS}/etc/.env"
    # Read only WITNESS_TAG from .env. Do NOT source the whole file: this runs
    # before downloadConfigFile normalizes .env, so a template-seeded file with
    # empty IOTEX_WITNESS= / DB_ROOT_PASSWORD= would clobber the already-computed
    # data dir (-> operating under /etc) and DB password (-> blank MySQL root).
    if [[ -f "$envFile" ]] && grep -qE '^[[:space:]]*WITNESS_TAG=' "$envFile"; then
        local line tag
        line=$(grep -E '^[[:space:]]*WITNESS_TAG=' "$envFile" | tail -1)
        tag=${line#*=}
        tag=${tag%$'\r'}                        # strip trailing CR
        tag=${tag#"${tag%%[![:space:]]*}"}      # ltrim
        # follow compose .env syntax: quoted value, or unquoted up to an inline
        # comment (docker tags contain no whitespace, so take the first token).
        if [[ $tag == \"*\"* ]]; then
            tag=${tag#\"}; tag=${tag%%\"*}
        elif [[ $tag == \'*\'* ]]; then
            tag=${tag#\'}; tag=${tag%%\'*}
        else
            tag=${tag%%[[:space:]]*}
        fi
        WITNESS_TAG="$tag"
    else
        WITNESS_TAG=${WITNESS_TAG:-$DEFAULT_WITNESS_TAG}
    fi
    if ! isFixedTag "$WITNESS_TAG"; then
        echo -e "${RED} WITNESS_TAG='${WITNESS_TAG}' in ${envFile} is not a fixed version. Set a release such as ${DEFAULT_WITNESS_TAG} (https://github.com/iotubeproject/ioTube/releases). Nothing was stopped. ${NC}"
        exit 1
    fi
    echo -e "Witness image: ghcr.io/iotubeproject/iotube-witness:${WITNESS_TAG}"
}

# The engine that owns the data dir: mariadb or mysql from files only that engine
# creates, "none" for a missing or empty dir, "unknown" when it cannot be read or
# is not recognized (the caller fails closed on unknown).
function dataDirEngine() {
    local dir="${IOTEX_WITNESS}/data/mysql" p entries
    if [[ ! -e "$dir" ]]; then
        # -e is also false when an ancestor cannot be traversed: the dir is only
        # known to be missing if its nearest existing ancestor can be searched.
        p=$dir
        while [[ ! -e "$p" && "$p" != "/" ]]; do p=$(dirname "$p"); done
        if [[ -d "$p" && -r "$p" && -x "$p" ]]; then echo none; else echo unknown; fi
        return
    fi
    if [[ ! -d "$dir" || ! -r "$dir" || ! -x "$dir" ]]; then
        echo unknown
    elif [[ -e "$dir/aria_log_control" ]]; then
        echo mariadb
    elif [[ -e "$dir/mysql.ibd" ]]; then
        echo mysql
    elif ! entries=$(ls -A "$dir" 2>/dev/null); then
        echo unknown
    elif [[ -z "$entries" ]]; then
        echo none
    else
        echo unknown
    fi
}

# Resolve WITNESS_DB_IMAGE (an existing .env entry wins over the environment, so
# a stray variable cannot switch engines) and refuse an image that cannot open the
# existing data dir. Runs before anything is stopped.
function resolveDbImage() {
    local envFile="${IOTEX_WITNESS}/etc/.env" line image engine want
    if [[ -f "$envFile" ]] && grep -qE '^[[:space:]]*WITNESS_DB_IMAGE=.' "$envFile"; then
        line=$(grep -E '^[[:space:]]*WITNESS_DB_IMAGE=' "$envFile" | tail -1)
        image=${line#*=}
        image=${image%$'\r'}
        image=${image//\"/}; image=${image//\'/}
        image=${image%%[[:space:]]*}
        if [[ -n "${WITNESS_DB_IMAGE:-}" && "$WITNESS_DB_IMAGE" != "$image" ]]; then
            echo -e "${YELLOW} Ignoring WITNESS_DB_IMAGE=${WITNESS_DB_IMAGE} from the environment; ${envFile} sets ${image}. ${NC}"
        fi
    else
        image=${WITNESS_DB_IMAGE:-mysql:8}
    fi
    case "${image##*/}" in
        mariadb*) want=mariadb ;;
        mysql*) want=mysql ;;
        *) echo -e "${RED} WITNESS_DB_IMAGE='${image}' is neither a mysql nor a mariadb image. Nothing was stopped. ${NC}"; exit 1 ;;
    esac
    engine=$(dataDirEngine)
    if [[ "$engine" == "unknown" ]]; then
        echo -e "${RED} Cannot tell which database engine created ${IOTEX_WITNESS}/data/mysql (not readable by $(id -un), or not a MySQL/MariaDB data dir). Run as a user that can read it. Nothing was stopped. ${NC}"
        exit 1
    fi
    if [[ "$engine" != "none" && "$engine" != "$want" ]]; then
        echo -e "${RED} ${IOTEX_WITNESS}/data/mysql was created by ${engine}, but WITNESS_DB_IMAGE='${image}'. Switching engines needs a dump and restore into a new data dir (see witness-service/README.md). Nothing was stopped. ${NC}"
        exit 1
    fi
    WITNESS_DB_IMAGE=$image
    echo -e "Witness database image: ${WITNESS_DB_IMAGE}"
}

function buildService() {
    docker pull "ghcr.io/iotubeproject/iotube-witness:${WITNESS_TAG}" || exit 2
}

function startup() {
    echo -e "$YELLOW Start witness and it's database. $NC"
    pushd $IOTEX_WITNESS/etc
    if ! docker compose up -d || ! docker compose restart; then
        echo -e "${RED} Failed to start the witness services; check 'docker compose ps' and the logs. ${NC}"
        exit 1
    fi
    docker image prune -f
    echo -e "${YELLOW} Service on. ${NC}"
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
    resolveDbPassword
    resolveWitnessTag
    resolveDbImage

    exportAll

    if [[ $# -ge 1 ]]; then
        if [[ $1 != 'restart' ]]; then
            echo -e "invalid parameter '$1'"
            exit
        fi
    else
        buildService
    fi

    stopService
    makeWorkspace
    grantPrivileges

    startup
}

main $@
