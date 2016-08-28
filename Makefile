.PHONY: all fmt lint build install
ALL_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
ALL_PKGS = $(shell glide nv)
PROG := $(shell basename $(CURDIR))


all: clean lint build


clean:
	rm -f $(PROG)


$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide


vendor install: $(GOPATH)/bin/glide
	glide up


$(PROG): vendor
	go build -o $(PROG) $(PROG).go


fmt:
	@echo Formatting Packages...
	go fmt $(ALL_FILES)


lint: vendor
	@echo "Running golint"
	golint $(ALL_PKGS)
	@echo "Running go vet"
	go vet $(ALL_PKGS)
	@echo "Checking gofmt"
	gofmt -l $(ALL_FILES)


test: vendor
	go test $(ALL_PKGS)


build: test $(PROG)
