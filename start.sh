#!/bin/bash
set -eux

if [ ! -f "/mindoc/conf/app.conf" ] ; then
    cp /mindoc/conf/app.conf.example /mindoc/conf/app.conf
fi

/mindoc/mindoc_linux_amd64 install

/mindoc/mindoc_linux_amd64