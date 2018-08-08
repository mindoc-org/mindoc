#!/bin/sh
set -e

cd /mindoc/

if [ ! -f "/mindoc/conf/app.conf" ] ; then
    cp /mindoc/conf/app.conf.example /mindoc/conf/app.conf
	sed -i "s#^db_adapter=.*#db_adapter=sqlite3#g" conf/app.conf
	sed -i "s#^db_database.*#db_database=./database/mindoc.db#g" conf/app.conf
fi


/mindoc/mindoc_linux_amd64 install
/mindoc/mindoc_linux_amd64