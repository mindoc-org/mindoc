FROM golang:1.8.3-alpine3.6

# add china aliyun repo  新增了 alpine 3.6 的阿里源
RUN cp   /etc/apk/repositories /etc/apk/repositories.back && \
    echo "https://mirrors.aliyun.com/alpine/v3.6/main/" >  /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.6/community/" >> /etc/apk/repositories

RUN apk add --update bash git make gcc g++

ADD . /go/src/github.com/lifei6671/mindoc

WORKDIR /go/src/github.com/lifei6671/mindoc

RUN chmod +x start.sh

#RUN  go get -d ./... && \  # 这行 docker build 不过去,所以注释掉了
RUN go get github.com/mitchellh/gox && \
    gox -os "windows linux darwin" -arch amd64
CMD ["./start.sh"]
