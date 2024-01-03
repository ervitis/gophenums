.PHONY: tests build format

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

tests: ## Execute tests
	rm -f `pwd`/_tests/simple/test1_gen.go && \
	go test -race -v ./... && \
	sleep 2 && \
	rm -f `pwd`/_tests/simple/test1_gen.go

build: ## Build the application
	go build -ldflags "-s -w" -o ./tmp/gophenums ./cmd/main.go

format: ## Format files
	go fmt ./...