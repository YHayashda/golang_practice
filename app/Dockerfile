FROM golang:latest

RUN mkdir /go/src/work

WORKDIR /go/src/work

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ADD . /go/src/work
