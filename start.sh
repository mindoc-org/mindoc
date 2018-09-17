#!/bin/sh
set -e

cd /mindoc/

if [ ! -f "/mindoc/conf/app.conf" ] ; then
    cp /mindoc/conf/app.conf.example /mindoc/conf/app.conf
fi

export ZONEINFO=/mindoc/lib/time/zoneinfo.zip
/mindoc/mindoc_linux_amd64 install
/mindoc/mindoc_linux_amd64