FROM golang:alpine

RUN mkdir -p /go/src/app
COPY . /go/src/app

WORKDIR /go/src/app
RUN go get -d -v
RUN go install -v

ENTRYPOINT ["app"]
