FROM amd64/golang:1.13 AS build

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
RUN go build -o mindoc_linux_amd64 -ldflags "-w -s -X main.VERSION=$TAG -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'"
RUN cp conf/app.conf.example conf/app.conf
# 清理不需要的文件
RUN rm appveyor.yml docker-compose.yml Dockerfile .travis.yml .gitattributes .gitignore go.mod go.sum main.go README.md simsun.ttc start.sh
RUN rm -rf cache commands controllers converter .git .github graphics mail models routers utils

# 测试编译的mindoc是否ok
# RUN ./mindoc_linux_amd64 version

# 必要的文件复制
ADD simsun.ttc /usr/share/fonts/win/
ADD start.sh /go/src/github.com/mindoc-org/mindoc


# Ubuntu 20.04
FROM ubuntu:focal

# 时区设置(如果不设置, calibre依赖的tzdata在安装过程中会要求选择时区)
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /usr/share/fonts/win/simsun.ttc /usr/share/fonts/win/
COPY --from=build /go/src/github.com/mindoc-org/mindoc /mindoc
WORKDIR /mindoc
RUN chmod a+r /usr/share/fonts/win/simsun.ttc

# 更新软件包信息
RUN apt-get update
# 安装必要的系统工具
RUN apt install -y apt-transport-https ca-certificates curl wget xz-utils
# 安装文泉驿字体
RUN apt install -y fonts-wqy-microhei fonts-wqy-zenhei
# 安装中文语言包
RUN apt-get install language-pack-zh-hans language-pack-zh-hans-base language-pack-gnome-zh-hans language-pack-gnome-zh-hans-base
# 设置默认编码
RUN localectl set-locale LANG=zh_CN.UTF-8
# 安装-calibre
# RUN apt-get install -y calibre
RUN mkdir -p /tmp/calibre-cache
# 获取最新版本号
RUN curl -s http://code.calibre-ebook.com/latest>/tmp/calibre-cache/version
# 下载最新版本
RUN wget -O /tmp/calibre-cache/calibre-x86_64.txz -c https://download.calibre-ebook.com/`cat /tmp/calibre-cache/version`/calibre-`cat /tmp/calibre-cache/version`-x86_64.txz
# 解压
RUN mkdir -p /opt/calibre
# RUN tar --extract --file=/tmp/calibre-cache/calibre-x86_64.txz --directory /opt/calibre
RUN tar xJof /tmp/calibre-cache/calibre-x86_64.txz -C /opt/calibre
ENV PATH=$PATH:/opt/calibre/bin
# 测试 calibre 可正常使用
RUN ebook-convert --version
ENV QTWEBENGINE_CHROMIUM_FLAGS="--no-sandbox"
ENV ZONEINFO=/mindoc/lib/time/zoneinfo.zip
RUN chmod +x start.sh

CMD ["./start.sh"]