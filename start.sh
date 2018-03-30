#!/bin/sh
set -e

cd $MINDOC_ROOT
appconf=conf/app.conf

if [ ! -f "$appconf" ] ; then
    cp conf.bak/app.conf.example $appconf
    sed -i "s#^db_adapter=.*#db_adapter=sqlite3#g" $appconf
    sed -i "s#^db_database.*#db_database=./database/mindoc.db#g" $appconf
fi

if [ ! -z $DB_TYPE ]; then
    sed -i "s#^db_adapter=.*#db_adapter=${DB_TYPE}#g" $appconf
fi

if [ ! -z $DB_HOST ] ; then
    sed -i 's/^db_host.*/db_host='$DB_HOST'/g' $appconf
fi

if [ ! -z $DB_PORT ] ; then
    sed -i 's/^db_port.*/db_port='$DB_PORT'/g' $appconf
fi

if [ ! -z $DB_NAME ] ; then
    sed -i "s#^db_database.*#db_database=${DB_NAME}#g" $appconf
fi

if [ ! -z $DB_USER ] ; then
    sed -i 's/^db_username.*/db_username='$DB_USER'/g' $appconf
fi

if [ ! -z $DB_PASS ] ; then
    sed -i 's/^db_password.*/db_password='$DB_PASS'/g' $appconf
fi

if [ ! -z $HTTP_PORT ] ; then
    sed -i 's/^httpport.*/httpport='$HTTP_PORT'/g' $appconf
fi

if [ ! -z $CDNJS ]; then
    sed -i "s#^cdnjs=.*#cdnjs=$CDNJS#g" $appconf
fi
if [ ! -z $CDNIMG ]; then
    sed -i "s#^cdnimg=.*#cdnimg=$CDNIMG#g" $appconf
fi

if [ ! -z $CDNCSS ]; then
    sed -i "s#^cdncss=.*#cdncss=$CDNCSS#g" $appconf
fi

if [ ! -z $CDN ]; then
    sed -i "s#^cdn=.*#cdn=$CDN#g" $appconf
fi

if [ ! -z $CACHE ]; then
    sed -i "s#cache=.*#cache=$CACHE#g" $appconf
fi

if [ ! -z $CACHE_PROVIDER ]; then
    sed -i "s#cache_provider=.*#cache_provider=$CACHE_PROVIDER#g" $appconf
fi

if [ ! -z $CACHE_MEMCACHE_HOST ]; then
    sed -i "s#cache_memcache_host=.*#cache_memcache_host=$CACHE_MEMCACHE_HOST#g" $appconf
fi

if [ ! -z $CACHE_REDIS_HOST ]; then
    sed -i "s#cache_redis_host=.*#cache_redis_host=$CACHE_REDIS_HOST#g" $appconf
fi

if [ ! -z $CACHE_REDIS_DB ]; then
    sed -i "s#cache_redis_db=.*#cache_redis_db=$CACHE_REDIS_DB#g" $appconf
fi

if [ ! -z $CACHE_REDIS_PASSWROD ]; then
    sed -i "s#cache_redis_password=.*#cache_redis_password=$CACHE_REDIS_PASSWROD#g" $appconf
fi

if [ ! -z $BASEURL ]; then
    sed -i "s#baseurl=.*#baseurl=$BASEURL#g" $appconf
fi

sed -i 's/^runmode.*/runmode=prod/g' $appconf

./mindoc_linux_amd64 install
./mindoc_linux_amd64