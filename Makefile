GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run
GOBUILD = env GO111MODULE=on go build

.PHONY: node build clean

build:
	make node 
dev:
	@echo "\n> --- run in development mode --"
	DEBUG=true DATA_DIR=./build go run ./cmd/node/main.go $(cmd)
node:
	mkdir -p $(GOBIN)
	cd ./cmd/node/ && go fmt ./... && $(GOBUILD) -o ./../../$(GOBIN)/block-node
	cp config.example.yml ./build/config.yml
	chmod +x $(GOBIN)/block-node

	@echo "\n> ---"
	@echo "> Build successful. Executable in: \"$(GOBIN)/block-node\" "
	@echo "> ---\n"
clean:
	rm -rf $(GOBIN)