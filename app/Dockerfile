FROM golang:1.14.1-buster

# $GOPATH == /go
RUN mkdir -p /go/src/github.com/RollMan/ac2manager
COPY . /go/src/github.com/RollMan/ac2manager
WORKDIR /go/src/github.com/RollMan/ac2manager
RUN go build -o app
CMD ["./app"]
