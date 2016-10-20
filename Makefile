.PHONY: all fmt lint build install
ALL_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
ALL_PKGS = $(shell glide nv)
PROG := $(shell grep "^[^=]" README.rst |head -1)

all: clean lint build

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

$(PROG): vendor
	go build -o $(PROG) $(ALL_PKGS)

build: $(PROG) test
	./$(PROG)

test: vendor
	go test -cover $(ALL_PKGS)

fmt:
	@echo Formatting Packages...
	go fmt $(ALL_FILES)
