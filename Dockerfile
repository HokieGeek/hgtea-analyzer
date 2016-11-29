FROM golang:alpine

RUN mkdir -p /go/src/github.com/hokiegeek/hgtealib
COPY . /go/src/github.com/hokiegeek/hgtealib

WORKDIR /go/src/github.com/hokiegeek/hgtealib
COPY ./sample_hgteas.json $HOME/.hgteas.json

RUN apk add --update git

RUN go get -d -v
RUN go install -v ./...

RUN apk del git

ENTRYPOINT ["teas"]
