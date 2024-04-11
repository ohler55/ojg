
VERSION = $(shell git describe --tags --always --dirty)

LDFLAG = "-s -w -X main.version=$(VERSION)"

build:
	go build -ldflags $(LDFLAG) .

.PHONY: build
