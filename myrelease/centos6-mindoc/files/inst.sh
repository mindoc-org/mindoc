#!/bash/bin

#install mindoc
mkdir -p /usr/local/mindoc
tar -zxf mindoc.tar.gz -C /usr/local/mindoc/
mv app.conf /usr/local/mindoc/conf/app.conf
mv start.sh /usr/local/mindoc/start.sh
chmod 777 /usr/local/mindoc/start.sh
