FROM golang:1.8.1-alpine


RUN apk add --update bash git make gcc

ADD . /go/src/github.com/lifei6671/godoc


WORKDIR /go/src/github.com/lifei6671/godoc

RUN chmod +x start.sh

RUN  go get -d ./... && \
    go build  -v -o godoc-linux-amd64 -ldflags "-w" 


CMD ["./start.sh"]