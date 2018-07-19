#!/bin/sh
set -e

cd /mindoc/

if [ ! -f "/mindoc/conf/app.conf" ] ; then
    cp /mindoc/conf/app.conf.example /mindoc/conf/app.conf
	sed -i "s#^db_adapter=.*#db_adapter=sqlite3#g" conf/app.conf
	sed -i "s#^db_database.*#db_database=./database/mindoc.db#g" conf/app.conf
fi

if [ ! -z $DB_ADAPTER ]; then
	sed -i "s#^db_adapter=.*#db_adapter=${DB_ADAPTER}#g" conf/app.conf
fi

if [ ! -z $MYSQL_PORT_3306_TCP_ADDR ] ; then
    sed -i 's/^db_host.*/db_host='$MYSQL_PORT_3306_TCP_ADDR'/g' conf/app.conf
fi

if [ ! -z $MYSQL_PORT_3306_TCP_PORT ] ; then
    sed -i 's/^db_port.*/db_port='$MYSQL_PORT_3306_TCP_PORT'/g' conf/app.conf
fi

if [ ! -z $MYSQL_INSTANCE_NAME ] ; then
    sed -i "s#^db_database.*#db_database=${MYSQL_INSTANCE_NAME}#g" conf/app.conf
fi

if [ ! -z $MYSQL_USERNAME ] ; then
    sed -i 's/^db_username.*/db_username='$MYSQL_USERNAME'/g' conf/app.conf
fi

if [ ! -z $MYSQL_PASSWORD ] ; then
    sed -i 's/^db_password.*/db_password='$MYSQL_PASSWORD'/g' conf/app.conf
fi

if [ ! -z $HTTP_PORT ] ; then
    sed -i 's/^httpport.*/httpport='$HTTP_PORT'/g' conf/app.conf
fi

if [ ! -z $CDNJS ]; then
    sed -i "s#^cdnjs=.*#cdnjs=$CDNJS#g" conf/app.conf
fi
if [ ! -z $CDNIMG ]; then
    sed -i "s#^cdnimg=.*#cdnimg=$CDNIMG#g" conf/app.conf
fi

if [ ! -z $CDNCSS ]; then
    sed -i "s#^cdncss=.*#cdncss=$CDNCSS#g" conf/app.conf
fi

if [ ! -z $CDN ]; then
    sed -i "s#^cdn=.*#cdn=$CDN#g" conf/app.conf
fi

if [ ! -z $CACHE ]; then
    sed -i "s#cache=.*#cache=$CACHE#g" conf/app.conf
fi

if [ ! -z $CACHE_PROVIDER ]; then
    sed -i "s#cache_provider=.*#cache_provider=$CACHE_PROVIDER#g" conf/app.conf
fi

if [ ! -z $CACHE_MEMCACHE_HOST ]; then
    sed -i "s#cache_memcache_host=.*#cache_memcache_host=$CACHE_MEMCACHE_HOST#g" conf/app.conf
fi

if [ ! -z $CACHE_REDIS_HOST ]; then
    sed -i "s#cache_redis_host=.*#cache_redis_host=$CACHE_REDIS_HOST#g" conf/app.conf
fi

if [ ! -z $CACHE_REDIS_DB ]; then
    sed -i "s#cache_redis_db=.*#cache_redis_db=$CACHE_REDIS_DB#g" conf/app.conf
fi

if [ ! -z $CACHE_REDIS_PASSWROD ]; then
    sed -i "s#cache_redis_password=.*#cache_redis_password=$CACHE_REDIS_PASSWROD#g" conf/app.conf
fi

if [ ! -z $BASEURL ]; then
    sed -i "s#baseurl=.*#baseurl=$BASEURL#g" conf/app.conf
fi

if [ ! -z $ENABLE_EXPORT ]; then
    sed -i "s#enable_export=.*#enable_export=$ENABLE_EXPORT#g" conf/app.conf
fi


sed -i 's/^runmode.*/runmode=prod/g' conf/app.conf

/mindoc/mindoc_linux_amd64 install
/mindoc/mindoc_linux_amd64