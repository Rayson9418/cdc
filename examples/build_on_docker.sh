#!/bin/bash

set -euo pipefail

trap cleanup SIGINT SIGTERM ERR EXIT

usage() {
    cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-t] [-r]...

Available options:

-h, --help              Print this help and exit
-t, --test              Install your own develop middleware with docker-compose.yml
-r, --reload            Reload base i18n file, and recreate new docker-compose.yml

Examples:

# Create network_[namespace], Install your own develop middleware, then enter the compile contailer
$(basename "${BASH_SOURCE[0]}") -t

# Reload base i18n file if in container, and recreate new docker-compose.yml
$(basename "${BASH_SOURCE[0]}") -r

EOF
    exit

}

cleanup() {
    [[ arg_help -ne 0 ]] && return

    trap - SIGINT SIGTERM ERR EXIT
}

setup_colors() {
    if [[ -t 2 ]] && [[ -z "${NO_COLOR-}" ]] && [[ "${TERM-}" != "dumb" ]]; then
        NOFORMAT='\033[0m' RED='\033[0;31m' GREEN='\033[0;32m' ORANGE='\033[0;33m' BLUE='\033[0;34m' PURPLE='\033[0;35m' CYAN='\033[0;36m' YELLOW='\033[1;33m'
    else
        NOFORMAT='' RED='' GREEN='' ORANGE='' BLUE='' PURPLE='' CYAN='' YELLOW=''
    fi
}

setup_colors

msg() {
    echo >&2 -e "${GREEN} [$(date '+%Y-%m-%d %H:%M:%S')] ${1-}${NOFORMAT}"
}

die() {
    local msg=$1
    local code=${2-1} # default exit status 1
    msg "$msg"
    exit "$code"
}

BASE_PATH=$(cd "$(dirname "${BASH_SOURCE[0]}}")" &>/dev/null && pwd -P)
IMAGE_NAME="go-compiler"
IMAGE_VERSION="v1"
COMPILE_IMAGE=$IMAGE_NAME:$IMAGE_VERSION

NAMESPACE=$(date +'%Y%m%d%H%M%S')
arg_help=0
arg_test=0
arg_reload=0
NETWORK_NAME="cdc_network"

parse_params() {

    while :; do
        case "${1-}" in 
        -h | --help)
            arg_help=1
            usage
            ;;
        -t | --test)
            arg_test=1
            ;;
        -r | --reload)
            arg_reload=1
            ;;
        -?*) die "Unknown option: $1" ;;
        *) break ;;
        esac
        shift
    done

    return 0
}

build_image() {
  docker build --network=host --build-arg VERSION=$IMAGE_VERSION -t $COMPILE_IMAGE -f Dockerfile .
}

reload() {
    local new_path=${BASE_PATH}
    
    [[ 0 -ne $(env | grep HOST_BASE_PATH | wc -l) ]] && {
        new_path=${HOST_BASE_PATH}
    }

    cp -f ${BASE_PATH}/tests/docker-compose.template ${BASE_PATH}/tests/docker-compose.yml
    sed -i 's|{{BASE_PATH}}|'${new_path}'|g' ${BASE_PATH}/tests/docker-compose.yml
    exit 0
}

install_test_env() {
    prepare

    [[ 0 -eq $arg_test ]] && return
    msg "prepare test environment"
    docker-compose -f ${BASE_PATH}/tests/docker-compose.yml up -d
    msg "prepare done"
}

prepare() {
    cp -f ${BASE_PATH}/tests/docker-compose.template ${BASE_PATH}/tests/docker-compose.yml
    sed -i 's|{{BASE_PATH}}|'${BASE_PATH}'|g' ${BASE_PATH}/tests/docker-compose.yml
}

remove_test_env() {
    [[ 0 -eq $arg_test ]] && return

    if [[ -f ${BASE_PATH}/tests/docker-compose.yml ]]; then
        docker-compose -f ${BASE_PATH}/tests/docker-compose.yml down
        rm -rf ${BASE_PATH}/tests/mysql/data/*
    fi
}

create_network() {
    if [[ 0 -eq $(docker network ls | grep "${NETWORK_NAME}" | wc -l) ]]; then
        msg "create docker network: ${NETWORK_NAME}"
        docker network create ${NETWORK_NAME}
    fi
}

main() {
    parse_params "$@"

    if [[ 1 -eq $arg_reload ]]; then
        reload
    fi

    build_image
    create_network
    install_test_env

    docker run --rm -it \
        --network ${NETWORK_NAME} \
        -v /root/.ssh:/root/.ssh \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v ${BASE_PATH}/../../cdc/:/test/cdc \
        -e HOST_BASE_PATH=${BASE_PATH} \
        -e UNIT_TEST_ENV=docker \
        -w /test/cdc/examples/ \
        ${COMPILE_IMAGE} bash

    remove_test_env
    docker network prune -f
}

main $*
exit 0