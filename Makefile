GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run
GOBUILD = env GO111MODULE=on go build

.PHONY: cli node build clean

build:
	make node cli
dev:
	@echo "\n> --- run in development mode --"
	DEBUG=true DATA_DIR=./build go run ./cmd/node/main.go
dev-cli:
	@echo "\n> --- run in development mode --"
	DEBUG=true DATA_DIR=./build go run ./cmd/cli/main.go	
node:
	mkdir -p $(GOBIN)
	cd ./cmd/node/ && go fmt ./... && $(GOBUILD) -o ./../../$(GOBIN)/block-node
	cp config.example.yml ./build/config.yml
	chmod +x $(GOBIN)/block-node

	@echo "\n> ---"
	@echo "> Build successful. Executable in: \"$(GOBIN)/block-node\" "
	@echo "> ---\n"
cli:
	mkdir -p $(GOBIN)
	cd ./cmd/cli/ && go fmt ./... && $(GOBUILD) -o ./../../$(GOBIN)/block-cli
	cp config.example.yml ./build/config.yml
	chmod +x $(GOBIN)/block-cli

	@echo "\n> ---"
	@echo "> Build successful. Executable in: \"$(GOBIN)/block-cli\" "
	@echo "> ---\n"
clean:
	rm -rf $(GOBIN)