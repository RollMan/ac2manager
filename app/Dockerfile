FROM golang:1.14.1-buster

# Install node and npm
RUN mkdir /node
WORKDIR /node
RUN curl -L https://raw.githubusercontent.com/tj/n/master/bin/n -o n && \
    bash n 14.7.0

# Copy files
RUN mkdir -p /go/src/github.com/RollMan/ac2manager/app
COPY . /go/src/github.com/RollMan/ac2manager/app

# Build javascripts by webpack
WORKDIR /go/src/github.com/RollMan/ac2manager/app/static
RUN npm  install -D webpack webpack-cli && npm run build

# Build application in Golang
WORKDIR /go/src/github.com/RollMan/ac2manager/app
RUN go build -o app
CMD ["./app"]
