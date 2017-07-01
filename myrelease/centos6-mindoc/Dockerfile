FROM wolcen/centos6-wkhtmltox:1.0.0
MAINTAINER wolcen@msn.com

ADD files /tmp
RUN cd /tmp;sh /tmp/inst.sh

EXPOSE 8181
CMD ["/usr/local/mindoc/start.sh"]
