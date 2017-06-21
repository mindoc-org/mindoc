FROM golang:1.9


RUN apk add --update bash git make gcc

ADD . /go/src/github.com/lifei6671/mindoc


WORKDIR /go/src/github.com/lifei6671/mindoc

RUN chmod +x start.sh
#RUN go get -d ./...
#RUN go get github.com/mitchellh/gox
#RUN go gox -os "windows linux darwin" -arch amd64
RUN go get -d ./... && \
    go get github.com/mitchellh/gox && \
    gox -os "linux" -arch amd64


CMD ["./start.sh"]
