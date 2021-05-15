FROM golang:latest

RUN mkdir /go/src/work \
    && go get go.mongodb.org/mongo-driver/mongo

WORKDIR /go/src/work

ADD . /go/src/work
