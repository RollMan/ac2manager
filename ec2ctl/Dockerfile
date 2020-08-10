FROM golang:1.14.1-buster

# Copy files
RUN mkdir -p /go/src/github.com/RollMan/ac2manager/ec2ctl
COPY . /go/src/github.com/RollMan/ac2manager/ec2ctl

# Build application in Golang
WORKDIR /go/src/github.com/RollMan/ac2manager/ec2ctl
RUN go build -o ec2ctl && go get github.com/go-delve/delve/cmd/dlv
CMD ["./ec2ctl"]
