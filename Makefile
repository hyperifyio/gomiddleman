.PHONY: build clean tidy

GOMIDDLEMAN_SOURCES := ./cmd/gomiddleman/main.go

all: build

build: gomiddleman

tidy:
	go mod tidy

gomiddleman: $(GOMIDDLEMAN_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o gomiddleman ./cmd/gomiddleman

clean:
	rm -f gomiddleman
