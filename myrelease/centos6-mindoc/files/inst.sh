#!/bash/bin

#install mindoc
mkdir -p /usr/local/mindoc
tar -zxf mindoc.tar.gz -C /usr/local/mindoc/
mv app.conf /usr/local/mindoc/conf/app.conf
mv start.sh /usr/local/mindoc/start.sh

#install wkhtmlto
yum install -y libXrender libXext fontconfig mkfontscale xz*
xz -d wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
rm -rf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
tar -xf wkhtmltox-0.12.4_linux-generic-amd64.tar -C /usr/local/
mv msyh.ttc /usr/share/fonts/msyh.ttc
mv simsun.ttc /usr/share/fonts/simsun.ttc
cd /usr/share/fonts
mkfontscale
mkfontdir
fc-cache
