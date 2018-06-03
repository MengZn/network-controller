# Makefile templete reference
# https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4
BINARY = network-controller
VET_REPORT = vet.report

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=linkernetworks
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: clean vet pb client server

pb:
	protoc ./messages/messages.proto --go_out=plugins=grpc:.

server:
	cd ${BUILD_DIR}/server; \
	go build ${LDFLAGS} -o ${BINARY}-server . ; \
	cd - >/dev/null

client:
	cd ${BUILD_DIR}/client; \
	go build ${LDFLAGS} -o ${BINARY}-client . ; \
	cd - >/dev/null

vet:
	go vet ./... > ${VET_REPORT} 2>&1 ; 

clean:
	-rm -f messages/messages.pb.go
	-rm -f client/${BINARY}-*
	-rm -f server/${BINARY}-*
	-rm -f ${VET_REPORT}

test: pb client server
	go clean -testcache 
	sudo -E env PATH=$$PATH TEST_OVS=1 go test -v ./...

.PHONY: server client vet test clean
