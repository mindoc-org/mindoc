#!/bin/bash
set -eux

# 数据库等初始化
/mindoc/mindoc_linux_amd64 install

# 导出同步检查
mkdir -p /mindoc-sync-host
if ! [ -f "/mindoc-sync-host/sync.sh" ]; then
    # 同步方向: docker->HOST 或 HOST -> docker
    # echo "export MINDOC_SYNC=" >> /mindoc-sync-host/sync.sh # 不同步
    echo "export MINDOC_SYNC=docker2host" >> /mindoc-sync-host/sync.sh # 默认 docker->HOST
    
    # 同步内容
    # conf: 配置
    # database: sqlite方式数据库
    # runtime: 运行时数据(日志等)
    # static: 静态文件
    # uploads: 上传文件
    # views: 页面视图
    # echo "export SYNC_LIST='conf;database;runtime;static;uploads;views'" >> /mindoc-sync-host/sync.sh # 同步所有内容
    # echo "export SYNC_LIST=" >> /mindoc-sync-host/sync.sh # 不同步任何内容
    echo "export SYNC_LIST='conf;database;uploads'" >> /mindoc-sync-host/sync.sh # 同步conf、database、uploads

    # 同步操作(sync/copy/sync --dry-run 等，具体参考rclone文档，host2docker务必谨慎操作)
    # echo "export SYNC_ACTION=sync --dry-run" >> /mindoc-sync-host/sync.sh # 无操作且仅显示同步文件信息(--dry-run)
    echo "export SYNC_ACTION=sync" >> /mindoc-sync-host/sync.sh # 默认同步
    
    # 同步脚本
    echo "source /mindoc/sync_host.sh" >> /mindoc-sync-host/sync.sh
fi
# 同步操作
source /mindoc-sync-host/sync.sh

# 运行
/mindoc/mindoc_linux_amd64