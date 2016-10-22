.PHONY: all build clean fmt install lint test
ALL_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
ALL_PKGS := $(shell glide nv)
PROG := $(shell grep "^[^=]" README.rst |head -1)
VERSION := $(shell grep -oP '^\d+\.\d+\.\d+(?= - \d{4}-\d{2}-\d{2}$$)' README.rst |head -1)
LDFLAGS := -X main.version=$(VERSION)

all: clean lint test build

clean:
	rm -f $(PROG)

$(GOPATH)/bin/golint:
	go get -u github.com/golang/lint/golint

lint: $(GOPATH)/bin/golint
	@echo "Running golint"
	golint $(ALL_PKGS)
	@echo "Running go vet"
	go vet $(ALL_PKGS)
	@echo "Checking gofmt"
	gofmt -l $(ALL_FILES) |(! grep '.')

$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide

install vendor: $(GOPATH)/bin/glide
	glide up

test: vendor
	go test -coverprofile cover.out -cover $(ALL_PKGS)

$(PROG): vendor
	go build -ldflags "$(LDFLAGS)" -o $(PROG) $(ALL_PKGS)

build: $(PROG)
	./$(PROG)

fmt:
	@echo Formatting Packages...
	go fmt $(ALL_FILES)
