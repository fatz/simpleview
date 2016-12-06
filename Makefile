# Go parameters
GOCMD=godep go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

TARGET=$(PWD)/target

build: bindatafs
	go get github.com/tools/godep
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGET)/simpleview.linux.amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(TARGET)/simpleview.darwin.amd64

clean:
	rm $(TARGET)/*

run:
	go-bindata-assetfs -debug public
	$(GOCMD) run bindata_assetfs.go main.go

test:
	$(GOCMD) test

bindatafs:
	go-bindata-assetfs public

all: test build
