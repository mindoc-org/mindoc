#!/bin/bash
set -eux

if ! [ -x "$(command -v rclone)" ]; then
    apt install rclone -y
else
    whereis rclone
fi


export MINDOC_SYNC="${MINDOC_SYNC:=}"
export SYNC_LIST="${SYNC_LIST:=}"
export SYNC_ACTION="${SYNC_ACTION:=sync --dry-run}"
export HOST_DIR=/mindoc-sync-host
export DOCKER_DIR=/mindoc

function doSyncCopy() {
    rclone $SYNC_ACTION --progress --exclude .git* --exclude .git/** "${1}" "${2}"
}

function doSync() {
    case $MINDOC_SYNC in
        "docker2host")
            doSyncCopy "${DOCKER_DIR}/${1}" "${HOST_DIR}/${1}"
            ;;
        "host2docker")
            doSyncCopy "${HOST_DIR}/${1}" "${DOCKER_DIR}/${1}"
            ;;
        *)
            printenv | grep MINDOC_SYNC
            ;;
    esac
}

export IFS=";"
if ! [ -z "${SYNC_LIST}" ]; then
    for item in $SYNC_LIST; do
        doSync "${item}"
    done
fi