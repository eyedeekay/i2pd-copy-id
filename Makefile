
GOPATH=$(shell pwd)/.go

echo:
	@echo "$(GOPATH)"
	find . -path ./.go -prune -o -name "*.go" -exec gofmt -w {} \;
	find . -path ./.go -prune -o -name "*.i2pkeys" -exec rm {} \;

deps:
	go get -u github.com/eyedeekay/sam-forwarder/daemon
	go get -u github.com/eyedeekay/sam3

build:
	go build -a -tags netgo \
		-ldflags '-w -extldflags "-static"'
