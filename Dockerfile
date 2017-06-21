FROM golang:1.8.1-alpine


RUN apk add --update bash git make gcc

ADD . /go/src/github.com/lifei6671/mindoc


WORKDIR /go/src/github.com/lifei6671/mindoc

RUN chmod +x start.sh

RUN go get -d ./... && \
    go get github.com/mitchellh/gox
RUN brew install glide
RUN glide install
RUN go build -ldflags "-w"
    
CMD ["./start.sh"]
