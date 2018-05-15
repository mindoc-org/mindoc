FROM golang:1.10-alpine3.7 as builder

## add china aliyun repo  新增了 alpine 3.6 的阿里源
#RUN cp   /etc/apk/repositories /etc/apk/repositories.back && \
#    echo "https://mirrors.aliyun.com/alpine/v3.6/main/" >  /etc/apk/repositories && \
#    echo "https://mirrors.aliyun.com/alpine/v3.6/community/" >> /etc/apk/repositories

RUN apk add --update bash git make gcc g++

ARG MINDOC_BUILD=/go/src/github.com/lifei6671/mindoc
ADD . $MINDOC_BUILD
WORKDIR $MINDOC_BUILD

RUN   go get -u github.com/golang/dep/cmd/dep && dep ensure && \
      CGO_ENABLE=0 go build -v -o mindoc_linux_amd64 -ldflags="-w -X main.VERSION=$TAG -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" && \
      rm -rf commands controllers models modules routers tasks vendor docs search data utils graphics .git Godeps uploads/* .gitignore .travis.yml Dockerfile gide.yaml LICENSE main.go README.md conf/enumerate.go conf/mail.go install.lock

################
# docker build -t trydofor/mindoc:0.10.1 .
################

FROM alpine:3.7

LABEL maintainer="trydofor"

ENV MINDOC_ROOT=/app
ENV HTTP_PORT=8181
ENV ZONEINFO=$MINDOC_ROOT/lib/time/zoneinfo.zip

WORKDIR  $MINDOC_ROOT
COPY --from=builder /go/src/github.com/lifei6671/mindoc .

RUN chmod +x start.sh
RUN mv conf conf.bak

VOLUME $MINDOC_ROOT/conf
VOLUME $MINDOC_ROOT/uploads
EXPOSE $HTTP_PORT

CMD ["./start.sh"]