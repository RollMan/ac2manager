SRC := $(shell find ./ -name "*.go")
IMAGE_NAME := ec2ctl
CREATION_TIME := $(shell docker inspect $(IMAGE_NAME) -f {{.Created}} | xargs date +%s -d)
LATEST_MOD_TIME := $(shell echo $(SRC) | xargs -n 1 date +%s -r | awk '{if(m<$$1) m=$$1} END{print m}')
GOPATH := /go/src/github.com/RollMan/ac2manager/ec2ctl

TESTARG := -v -timeout 20m

build: $(SRC)
	if [ -z $(CREATION_TIME) ] || [ $(LATEST_MOD_TIME) -ge $(CREATION_TIME) ]; then \
		docker build . -t $(IMAGE_NAME); \
		else \
		echo "Docker container image is ready up to date."; \
	fi

build-vol: $(SRC)
	docker run -v$(PWD):$(GOPATH) --rm -it ec2ctl go build -o ec2ctl

test: build
	docker run --env-file ../.env --rm -it ec2ctl go test $(TESTARG) $(PKG)

test-build:
	docker run -v$(PWD):$(GOPATH) --rm -it ec2ctl go test -c $(TESTARG) $(PKG)

test-vol:
	docker run -v$(PWD):$(GOPATH) --rm -it ec2ctl go test $(TESTARG) $(PKG)

test-dbg:
	docker run -v$(PWD):$(GOPATH) --rm -it ec2ctl dlv test $(PKG)

gomod: build
	docker run --rm -it -v`pwd`:/go/src/github.com/RollMan/ac2manager/ec2ctl ec2ctl go mod tidy
