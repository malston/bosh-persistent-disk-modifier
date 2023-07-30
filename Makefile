BINARY_NAME ?= pdm
BINARY_OUTPUT = ./bin/$(BINARY_NAME)
BINARY_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin
GITSHA = $(shell git rev-parse HEAD)
GITDIRTY = $(shell git diff --quiet HEAD || echo "dirty")

.PHONY: all
all: clean install

.PHONY: clean
clean: ## Delete build output
	@rm -rf bin/
	@rm -rf dist/

$(BINARY_OUTPUT): $(BINARY_SOURCES)
	go build -o $(BINARY_OUTPUT) ./main.go

.PHONY: build
build: $(BINARY_OUTPUT) ## Build the main binary

.PHONY: install
install: build ## Copy build to GOPATH/bin
	@cp $(BINARY_OUTPUT) $(GOBIN)
	@echo "[OK] CLI binary installed under $(GOBIN)"

.PHONY: release
release: $(BINARY_SOURCES) ## Cross-compile binary for various operating systems
	@mkdir -p dist
	GOOS=darwin   GOARCH=amd64 go build -o $(BINARY_OUTPUT)     ./main.go && tar -czf dist/$(BINARY_NAME)-darwin-amd64.tgz -C bin . && rm -f $(BINARY_OUTPUT)
	GOOS=linux    GOARCH=amd64 go build -o $(BINARY_OUTPUT)     ./main.go && tar -czf dist/$(BINARY_NAME)-linux-amd64.tgz  -C bin . && rm -f $(BINARY_OUTPUT)
	GOOS=windows  GOARCH=amd64 go build -o $(BINARY_OUTPUT).exe ./main.go && zip -rj  dist/$(BINARY_NAME)-windows-amd64.zip   bin   && rm -f $(BINARY_OUTPUT).exe

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_2-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
