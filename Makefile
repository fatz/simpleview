# Go parameters
GOCMD=godep go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

TARGET=$(PWD)/target
FILES="bindata_assetfs.go main.go"

req:
	go get github.com/tools/godep
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	godep restore

build: req bindatafs
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(FILES) -o $(TARGET)/simpleview.linux.amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(FILES) -o $(TARGET)/simpleview.darwin.amd64

clean:
	rm $(TARGET)/*

run:
	go-bindata-assetfs -debug public
	$(GOCMD) run bindata_assetfs.go main.go

test: req bindatafs
	$(GOCMD) test *.go

bindatafs:
	go-bindata-assetfs public

all: req bindatafs test build
