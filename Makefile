.PHONY: cover

BINS := find xargs cat
BIN_DIR := ./bin
CMD_DIR := ./cmd

all: $(BINS)

$(BINS):
	go build -o $(BIN_DIR)/$@ $(CMD_DIR)/$@

test:
	go test -race -v ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
