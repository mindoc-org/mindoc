FROM golang:1.10-alpine3.7

# add china aliyun repo  新增了 alpine 3.6 的阿里源
RUN cp   /etc/apk/repositories /etc/apk/repositories.back && \
    echo "https://mirrors.aliyun.com/alpine/v3.6/main/" >  /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.6/community/" >> /etc/apk/repositories

RUN apk add --update bash git make gcc g++

ADD . /go/src/github.com/lifei6671/mindoc

WORKDIR /go/src/github.com/lifei6671/mindoc

RUN chmod +x start.sh

RUN   go get -u github.com/golang/dep/cmd/dep && dep ensure && \
    CGO_ENABLE=1 go build -v -o mindoc_linux_amd64 -ldflags="-w -X main.VERSION=$TAG -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" && \
        rm -rf commands controllers models modules routers tasks vendor docs search data utils graphics .git Godeps uploads/* .gitignore .travis.yml Dockerfile gide.yaml LICENSE main.go README.md conf/enumerate.go conf/mail.go install.lock


CMD ["./start.sh"]