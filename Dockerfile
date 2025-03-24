FROM golang:bookworm AS build

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


# upgrade to the latest
FROM ubuntu:latest

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

RUN sed -i "s/archive.ubuntu.com/mirrors.aliyun.com/g" /etc/apt/sources.list /etc/apt/sources.list.d/*


# 更新软件包信息
RUN apt-get update

# 时区设置(如果不设置, calibre依赖的tzdata在安装过程中会要求选择时区)
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# tzdata的前端类型默认为readline（Shell情况下）或dialog（支持GUI的情况下）
ARG DEBIAN_FRONTEND=noninteractive
# 安装时区信息
RUN apt install -y --no-install-recommends tzdata
# 重新配置tzdata软件包，使得时区设置生效
RUN dpkg-reconfigure --frontend noninteractive tzdata

# 安装文泉驿字体
# 安装中文语言包
RUN apt install -y fonts-wqy-microhei fonts-wqy-zenhei locales language-pack-zh-hans-base
# 设置默认编码
RUN locale-gen "zh_CN.UTF-8"
RUN update-locale LANG=zh_CN.UTF-8
ENV LANG=zh_CN.UTF-8
ENV LANGUAGE=zh_CN:en
ENV LC_ALL=zh_CN.UTF-8

# 安装必要依赖、下载、解压 calibre 并清理缓存
RUN apt-get install -y --no-install-recommends \
      libgl-dev libnss3-dev libxcomposite-dev libxrandr-dev libxi-dev libxdamage-dev \
      wget xz-utils && \
    mkdir -p /tmp/calibre-cache /opt/calibre && \
    wget -O /tmp/calibre-cache/calibre-x86_64.txz -c https://download.calibre-ebook.com/7.26.0/calibre-7.26.0-x86_64.txz && \
    tar xJof /tmp/calibre-cache/calibre-x86_64.txz -C /opt/calibre && \
    rm -rf /tmp/calibre-cache && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 设置环境变量
ENV PATH="/opt/calibre:$PATH" \
    QTWEBENGINE_CHROMIUM_FLAGS="--no-sandbox" \
    QT_QPA_PLATFORM="offscreen"

# 测试 calibre 是否可正常使用
RUN ebook-convert --version

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
# docker run -d --name=mindoc --restart=always -v /www/mindoc/uploads:/mindoc/uploads -v /www/mindoc/database:/mindoc/database  -v /www/mindoc/conf:/mindoc/conf  -e MINDOC_DB_ADAPTER=sqlite3 -e MINDOC_DB_DATABASE=./database/mindoc.db -e MINDOC_CACHE=true -e MINDOC_CACHE_PROVIDER=file -p 8181:8181 mindoc-org/mindoc:v2.1
