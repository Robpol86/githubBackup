.PHONY: all fmt lint build install
ALL_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
ALL_PKGS = $(shell glide nv)
PROG := $(shell basename $(CURDIR))
README_PARSED_FILE := gen_readme_parsed.go


all: clean lint build


clean:
	rm -f $(PROG) $(README_PARSED_FILE)


$(GOPATH)/bin/golint:
	go get -u github.com/golang/lint/golint


$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide


$(README_PARSED_FILE): USAGE = TODO
$(README_PARSED_FILE): VERSION = 0.0.1
$(README_PARSED_FILE):
	@echo "package main\n\nconst (" > $(README_PARSED_FILE)
	@echo "\tusage = \"$(USAGE)\"" >> $(README_PARSED_FILE)
	@echo "\tversion = \"$(VERSION)\"" >> $(README_PARSED_FILE)
	@echo ")" >> $(README_PARSED_FILE)


vendor install: $(GOPATH)/bin/glide
	glide up


$(PROG): vendor $(README_PARSED_FILE)
	go build -o $(PROG) $(ALL_PKGS)


fmt:
	@echo Formatting Packages...
	go fmt $(ALL_FILES)


lint: $(GOPATH)/bin/golint $(README_PARSED_FILE)
	@echo "Running golint"
	golint $(ALL_PKGS)
	@echo "Running go vet"
	go vet $(ALL_PKGS)
	@echo "Checking gofmt"
	gofmt -l $(ALL_FILES)


test: vendor
	go test $(ALL_PKGS)


build: test $(PROG)
	./$(PROG)
	./$(PROG) --help
