# Go parameters
PROJECT_NAME=PIAAS
VERSION=v0.0.3
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=piaas
TESTMODE?=
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test$(TESTMODE) build-all

build:
	cd main && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/darwin_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/$(VERSION)/darwin_amd64/$(BINARY_NAME)
	@echo "        Built darwin-amd64"
test:
	$(GOTEST) -v ./...
testverbose:
	$(GOTEST) -v ./... -args -logtostderr
tests:
	$(GOTEST) -v ./... -count=10
testsverbose:
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
build-all: build build-linux build-windows
	@echo "\n$(PROJECT_NAME) version: $(VERSION) was built.\n"

# Cross compilation on linux
build-linux:
	cd main && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/linux_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/$(VERSION)/linux_amd64/$(BINARY_NAME)
	@echo "        Built linux-amd64"

build-windows:
	cd main && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/windows_amd64/$(BINARY_NAME).exe -v && cd ..
	chmod a+x dist/$(VERSION)/windows_amd64/$(BINARY_NAME).exe
	@echo "        Built windows-amd64"

# Publish new release
publish: tests$(TESTMODE) build-all
	# Tag the main repo with the release version
	git tag -f -a $(VERSION) -m "$(VERSION)"
	git push origin :refs/tags/$(VERSION)
	git push origin $(VERSION)
	git push

	# Copy the released file to gh-pages
	@rm -rf tmp/piaas-gh-pages
	git clone https://github.com/sohoffice/piaas.git tmp/piaas-gh-pages -b gh-pages --single-branch
	rm -rf tmp/piaas-gh-pages/files/$(VERSION)
	cp -R dist/$(VERSION) tmp/piaas-gh-pages/files
	rm -rf tmp/piaas-gh-pages/files/latest
	cp -R dist/$(VERSION) tmp/piaas-gh-pages/files/latest
	cd tmp/piaas-gh-pages && git add . && git commit -m "Release version $(VERSION)"
	cd tmp/piaas-gh-pages && git push

	@echo "\n$(PROJECT_NAME) version: $(VERSION) was released.\n\n"
