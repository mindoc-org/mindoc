#!/bin/sh
set -e

cd /go/src/github.com/lifei6671/godoc/

goFile="godoc"


chmod +x $goFile

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

if [ ! -z $HTTP_PORT] ; then
    sed -i 's/^httpport.*/httpport='$HTTP_PORT'/g' conf/app.conf
fi

for GOOS in darwin linux windows; do
   for GOARCH in 386 amd64; do
     go build -v -o godoc-$GOOS-$GOARCH
   done
done

/go/src/github.com/lifei6671/godoc/godoc-linux-amd64