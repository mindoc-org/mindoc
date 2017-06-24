#FROM golang:1.8.1-alpine
FROM mindoc:0.2.0.1


#RUN apk add --update bash git make gcc

ADD . /go/src/github.com/lifei6671/mindoc


WORKDIR /go/src/github.com/lifei6671/mindoc

RUN chmod +x start.sh

#RUN  go get -d ./... && \
#    go get github.com/mitchellh/gox && \
#    gox -os "windows linux darwin" -arch amd64
RUN  gox -os "windows linux darwin" -arch amd64


CMD ["./start.sh"]