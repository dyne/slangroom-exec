.PHONY: help

DEPS = bun
K := $(foreach exec,$(DEPS),\
        $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

help: ## 🛟 Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m 👉🏻 %-7s\033[0m %s\n", $$1, $$2}'


slangroom-exec: index.ts ## 🛠️ Build slangroom-exec
	bun build ./index.ts --compile --outfile slangroom-exec

clean: ## 🧹 Clean the build
	@rm -f slangroom-exec
	@echo "🧹 Cleaned the build"

tests: slangroom-exec ## 🧪 Run tests
	./test/bats/bin/bats test/test.bats
