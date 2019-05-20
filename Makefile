# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=piaas

all: test build
build:
	mkdir -p dist/darwin-amd64
	cd main && GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ../dist/darwin-amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/darwin-amd64/$(BINARY_NAME)
test:
	$(GOTEST) -v ./... -count=10 -args -logtostderr
clean:
	$(GOCLEAN)
	rm -rf dist
run: build
	./dist/darwin-amd64/$(BINARY_NAME)
deps:
	$(GOGET) github.com/markbates/goth
	$(GOGET) github.com/markbates/pop


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
