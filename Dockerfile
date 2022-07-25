FROM amd64/golang:1.18.1 AS build

ARG TAG=0.0.1

# 编译-环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1
ENV GOARCH=amd64
ENV GOOS=linux

# 工作目录
ADD . /go/src/github.com/mindoc-org/mindoc
WORKDIR /go/src/github.com/mindoc-org/mindoc

# 编译
RUN go env
RUN go mod tidy -v
RUN go build -v -o mindoc_linux_amd64 -ldflags "-w -s -X 'main.VERSION=$TAG' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'"
RUN cp conf/app.conf.example conf/app.conf
# 清理不需要的文件
RUN rm appveyor.yml docker-compose.yml Dockerfile .travis.yml .gitattributes .gitignore go.mod go.sum main.go README.md simsun.ttc start.sh conf/*.go
RUN rm -rf cache commands controllers converter .git .github graphics mail models routers utils

# 测试编译的mindoc是否ok
RUN ./mindoc_linux_amd64 version

# 必要的文件复制
ADD simsun.ttc /usr/share/fonts/win/
ADD start.sh /go/src/github.com/mindoc-org/mindoc


# Ubuntu 20.04
FROM ubuntu:focal

# 切换默认shell为bash
SHELL ["/bin/bash", "-c"]

WORKDIR /mindoc

# 文件复制
COPY --from=build /usr/share/fonts/win/simsun.ttc /usr/share/fonts/win/
COPY --from=build /go/src/github.com/mindoc-org/mindoc/mindoc_linux_amd64 /mindoc/
COPY --from=build /go/src/github.com/mindoc-org/mindoc/start.sh /mindoc/
COPY --from=build /go/src/github.com/mindoc-org/mindoc/LICENSE.md /mindoc/
# 文件夹复制
COPY --from=build /go/src/github.com/mindoc-org/mindoc/lib /mindoc/lib
COPY --from=build /go/src/github.com/mindoc-org/mindoc/conf /mindoc/__default_assets__/conf
COPY --from=build /go/src/github.com/mindoc-org/mindoc/static /mindoc/__default_assets__/static
COPY --from=build /go/src/github.com/mindoc-org/mindoc/views /mindoc/__default_assets__/views
COPY --from=build /go/src/github.com/mindoc-org/mindoc/uploads /mindoc/__default_assets__/uploads

RUN chmod a+r /usr/share/fonts/win/simsun.ttc

# 备份原有源
RUN mv /etc/apt/sources.list /etc/apt/sources.list-backup
# 最小化源，缩短apt update时间(ca-certificates必须先安装才支持换aliyun源)
RUN echo 'deb http://archive.ubuntu.com/ubuntu/ focal main restricted' > /etc/apt/sources.list
RUN apt-get update
RUN apt install -y ca-certificates
# 更换aliyun源(echo多行内容不能以#开头，会被docker误判为注释行，所以采用\n#开头)
RUN echo $'\
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse\
\n# deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse\
\n# deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse\
\n# deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse\
\n# deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse\
\n# deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse'\
> /etc/apt/sources.list

# 更新软件包信息
RUN apt-get update
# 安装必要的系统工具
RUN apt install -y apt-transport-https ca-certificates curl wget xz-utils

# 时区设置(如果不设置, calibre依赖的tzdata在安装过程中会要求选择时区)
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# tzdata的前端类型默认为readline（Shell情况下）或dialog（支持GUI的情况下）
ARG DEBIAN_FRONTEND=noninteractive
# 安装时区信息
RUN apt install -y --no-install-recommends tzdata
# 重新配置tzdata软件包，使得时区设置生效
RUN dpkg-reconfigure --frontend noninteractive tzdata

# 安装 calibre 依赖的包
RUN apt install -y libgl-dev libnss3-dev libxcomposite-dev libxrandr-dev libxi-dev libxdamage-dev
# 安装文泉驿字体
RUN apt install -y fonts-wqy-microhei fonts-wqy-zenhei
# 安装中文语言包
RUN apt-get install -y locales language-pack-zh-hans language-pack-zh-hans-base
# 设置默认编码
RUN locale-gen "zh_CN.UTF-8"
RUN update-locale LANG=zh_CN.UTF-8
ENV LANG=zh_CN.UTF-8
ENV LANGUAGE=zh_CN:en
ENV LC_ALL=zh_CN.UTF-8

# 安装-calibre
# RUN apt-get install -y calibre # 此种方式安装省事，但会安装很多额外不需要的软件包，导致体积过大
RUN mkdir -p /tmp/calibre-cache
# 获取最新版本号
RUN curl -s http://code.calibre-ebook.com/latest>/tmp/calibre-cache/version
# 下载最新版本
# RUN wget -O /tmp/calibre-cache/calibre-x86_64.txz -c https://download.calibre-ebook.com/`cat /tmp/calibre-cache/version`/calibre-`cat /tmp/calibre-cache/version`-x86_64.txz
# 使用 ghproxy.com 替换 github 实现加速
# RUN wget -O /tmp/calibre-cache/calibre-x86_64.txz -c https://ghproxy.com/https://github.com/kovidgoyal/calibre/releases/download/v`cat /tmp/calibre-cache/version`/calibre-`cat /tmp/calibre-cache/version`-x86_64.txz
RUN wget -O /tmp/calibre-cache/calibre-x86_64.txz -c https://github.com/kovidgoyal/calibre/releases/download/v`cat /tmp/calibre-cache/version`/calibre-`cat /tmp/calibre-cache/version`-x86_64.txz
# 注: 调试阶段，下载alibre-5.22.1-x86_64.txz到本地(使用 python -m http.server)，加速构建
# RUN wget -O /tmp/calibre-cache/calibre-x86_64.txz -c http://10.96.8.252:8000/calibre-5.22.1-x86_64.txz
# 解压
RUN mkdir -p /opt/calibre
# RUN tar --extract --file=/tmp/calibre-cache/calibre-x86_64.txz --directory /opt/calibre
RUN tar xJof /tmp/calibre-cache/calibre-x86_64.txz -C /opt/calibre
ENV PATH=$PATH:/opt/calibre
# 设置calibre相关环境变量
ENV QTWEBENGINE_CHROMIUM_FLAGS="--no-sandbox"
ENV QT_QPA_PLATFORM='offscreen'
# 测试 calibre 可正常使用
RUN ebook-convert --version
# 清理calibre缓存
RUN rm -rf /tmp/calibre-cache

# refer: https://docs.docker.com/engine/reference/builder/#volume
VOLUME ["/mindoc/conf","/mindoc/static","/mindoc/views","/mindoc/uploads","/mindoc/runtime","/mindoc/database"]

# refer: https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8181/tcp

ENV ZONEINFO=/mindoc/lib/time/zoneinfo.zip
RUN chmod +x /mindoc/start.sh

ENTRYPOINT ["/bin/bash", "/mindoc/start.sh"]

# https://docs.docker.com/engine/reference/commandline/build/#options
# docker build --progress plain --rm --build-arg TAG=2.1 --tag gsw945/mindoc:2.1 .
# https://docs.docker.com/engine/reference/commandline/run/#options
# set MINDOC=//d/mindoc # windows
# export MINDOC=/home/ubuntu/mindoc-docker # linux
# docker run --rm -it  -p 8181:8181 -v "%MINDOC%":"/mindoc-sync-host" --name mindoc -e MINDOC_ENABLE_EXPORT=true -d gsw945/mindoc:2.1
