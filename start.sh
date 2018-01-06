#!/bin/sh
set -e

cd /go/src/github.com/lifei6671/mindoc/

if [ ! -f "conf/app.conf" ] ; then
    cp conf/app.conf.example conf/app.conf
fi

if [ ! -z $MYSQL_PORT_3306_TCP_ADDR ] ; then
    sed -i 's/^db_host.*/db_host='$MYSQL_PORT_3306_TCP_ADDR'/g' conf/app.conf
fi

if [ ! -z $MYSQL_PORT_3306_TCP_PORT ] ; then
    sed -i 's/^db_port.*/db_port='$MYSQL_PORT_3306_TCP_PORT'/g' conf/app.conf
fi

if [ ! -z $MYSQL_INSTANCE_NAME ] ; then
    sed -i 's/^db_database.*/db_database='$MYSQL_INSTANCE_NAME'/g' conf/app.conf
fi

if [ ! -z $MYSQL_USERNAME ] ; then
    sed -i 's/^db_username.*/db_username='$MYSQL_USERNAME'/g' conf/app.conf
fi

if [ ! -z $MYSQL_PASSWORD ] ; then
    sed -i 's/^db_password.*/db_password='$MYSQL_PASSWORD'/g' conf/app.conf
fi

if [ ! -z $HTTP_PORT ] ; then
    sed -i "s/^httpport.*/httpport=${HTTP_PORT}/g" conf/app.conf
fi

if [ ! -z $CDNJS ]; then
    sed -i "s#^cdnjs.*#cdnjs=$CDNJS#g" conf/app.conf
fi

if [ ! -z $CDNCSS ]; then
    sed -i "s#^cdncss.*#cdncss=$CDNCSS#g" conf/app.conf
fi

if [ ! -z $CDN ]; then
    sed -i "s#^cdn=.*#cdn=$CDN#g" conf/app.conf
fi

sed -i 's/^runmode.*/runmode=prod/g' conf/app.conf

/go/src/github.com/lifei6671/mindoc/mindoc_linux_amd64 install

/go/src/github.com/lifei6671/mindoc/mindoc_linux_amd64