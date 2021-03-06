include .version
export

# Go parameters
PROJECT_NAME=PIAAS
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=piaas
TESTMODE?=all
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test build-all

build:
	cd main && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/darwin_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/darwin_amd64/$(BINARY_NAME)
	@echo "        Built darwin-amd64"
test: test$(TESTMODE)
testall:
	$(GOTEST) -v ./...
testdebug:
	$(GOTEST) -v ./... -args -debug
tests: tests$(TESTMODE)
testsall:
	$(GOTEST) -v ./... -count=10
testsdebug:
	$(GOTEST) -v ./... -count=10 -args -debug
clean:
	$(GOCLEAN)
	rm -rf dist
run: build
	./dist/darwin_amd64/$(BINARY_NAME) sync cent -debug
cover:
	go test ./... -cover
deps:
	$(GOGET) github.com/markbates/goth
	$(GOGET) github.com/markbates/pop

# Build on all supported platforms
build-all: build build-linux build-windows
	@echo "\n$(PROJECT_NAME) version: $(VERSION) was built.\n"

# Cross compilation on linux
build-linux:
	cd main && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/linux_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/linux_amd64/$(BINARY_NAME)
	@echo "        Built linux-amd64"

build-windows:
	cd main && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/windows_amd64/$(BINARY_NAME).exe -v && cd ..
	chmod a+x dist/windows_amd64/$(BINARY_NAME).exe
	@echo "        Built windows-amd64"

# Create a new release for publish
release: clean test$(TESTMODE) build-all
	cd dist/darwin_amd64 && zip piaas-darwin-amd64-$(VERSION).zip * && cd ..
	cd dist/linux_amd64 && zip piaas-linux-amd64-$(VERSION).zip * && cd ..
	cd dist/windows_amd64 && zip piaas-windows-amd64-$(VERSION).zip * && cd ..
