FROM golang:alpine

RUN mkdir -p /go/src/github.com/hokiegeek/hgtealib
COPY . /go/src/github.com/hokiegeek/hgtealib

WORKDIR /go/src/github.com/hokiegeek/hgtealib
RUN go get -d -v
RUN go install -v ./...

ENTRYPOINT ["hgtea"]
