FROM golang:1.14.1-buster

RUN mkdir /work

WORKDIR /work
COPY . .
RUN go build -o app
CMD ["./app"]
