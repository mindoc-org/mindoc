#!/bin/bash
set -eux

# 默认资源
if [ ! -d "/mindoc/conf" ]; then mkdir -p "/mindoc/conf" ; fi
if [[ -z "$(ls -A -- "/mindoc/conf")" ]] ; then cp -r "/mindoc/__default_assets__/conf" "/mindoc/" ; fi

if [ ! -d "/mindoc/static" ]; then mkdir -p "/mindoc/static" ; fi
if [[ -z "$(ls -A -- "/mindoc/static")" ]] ; then cp -r "/mindoc/__default_assets__/static" "/mindoc/" ; fi

if [ ! -d "/mindoc/views" ]; then mkdir -p "/mindoc/views" ; fi
if [[ -z "$(ls -A -- "/mindoc/views")" ]] ; then cp -r "/mindoc/__default_assets__/views" "/mindoc/" ; fi

if [ ! -d "/mindoc/uploads" ]; then mkdir -p "/mindoc/uploads" ; fi
if [[ -z "$(ls -A -- "/mindoc/uploads")" ]] ; then cp -r "/mindoc/__default_assets__/uploads" "/mindoc/" ; fi

# 如果配置文件不存在就复制
cp --no-clobber /mindoc/conf/app.conf.example /mindoc/conf/app.conf

# 数据库等初始化
/mindoc/mindoc_linux_amd64 install

# 运行
/mindoc/mindoc_linux_amd64

# # Debug Dockerfile
# while [ 1 ]
# do
#     echo "log ..."
#     sleep 5s
# done