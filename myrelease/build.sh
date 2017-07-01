#!/bin/bash

# clean
rm -rf centos6-mindoc-source.tar.gz
rm -rf centos6-mindoc/files/mindoc.tar.gz

cd ..
tar -zcf centos6-mindoc/files/mindoc.tar.gz . --exclude myrelease --exclude .git --exclude .github --exclude logs --exclude  database --exclude conf/app.conf --exclude *.go --exclude *.tmp --exclude appveyor.yml --exclude glide.yaml --exclude Dockerfile  --exclude start.h 
cd myrelease

tar -zcf centos6-mindoc-source.tar.gz centos6-mindoc

#copy centos6-mindoc-source.tar.gz to docker manger and tar
#cd centos6-mindoc-source
#docker build -t mindoc:1.0.0 .


