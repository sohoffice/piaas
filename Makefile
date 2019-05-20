# Go parameters
PROJECT_NAME=PIAAS
VERSION=v0.0.1
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=piaas
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test build-all

build:
	mkdir -p dist/darwin-amd64
	cd main && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/darwin_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/$(VERSION)/darwin_amd64/$(BINARY_NAME)
	@echo "        Built darwin-amd64"
test:
	$(GOTEST) -v ./... -args -logtostderr
tests:
	$(GOTEST) -v ./... -count=10 -args -logtostderr
clean:
	$(GOCLEAN)
	rm -rf dist
run: build
	./dist/darwin-amd64/$(BINARY_NAME)
deps:
	$(GOGET) github.com/markbates/goth
	$(GOGET) github.com/markbates/pop

# Build on all supported platforms
build-all: build build-linux
	@echo "\n$(PROJECT_NAME) version: $(VERSION)\n"

# Cross compilation on linux
build-linux:
	cd main && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/linux_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/$(VERSION)/linux_amd64/$(BINARY_NAME)
	@echo "        Built linux-amd64"
