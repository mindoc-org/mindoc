#!/bin/bash

#在开发编译容器里
cd ..
git pull
if [ -f "mindoc" ];then
  rm -rf mindoc
fi
go build
if [ ! -f "mindoc" ];then
  echo "编译生成mindoc失败"
  exit
fi
if [ -f "myrelease/centos6-mindoc/files/mindoc.tar.gz" ];then
   rm -rf myrelease/centos6-mindoc/files/mindoc.tar.gz
fi
tar -zcf centos6-mindoc/files/mindoc.tar.gz . --exclude myrelease --exclude .git --exclude .github --exclude logs --exclude  database --exclude conf/app.conf --exclude *.go --exclude *.tmp --exclude appveyor.yml --exclude glide.yaml --exclude Dockerfile  --exclude start.h 

cd myrelease
if [ -f "centos6-mindoc-release.tar.gz" ];then
  rm -rf centos6-mindoc-release.tar.gz
fi
tar -zcf centos6-mindoc-release.tar.gz centos6-mindoc

# 以下在容器管理机上执行
# 先复制 centos6-mindoc-release.tar.gz 到管理机
#tar -zxf centos6-mindoc-release.tar.gz .
#cd centos6-mindoc-release
#docker build -t mindoc:1.0.0 .


