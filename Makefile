.PHONY: help

SOURCES = $(shell find src -type f -name '*.ts')
LIBS = node_modules

DEPS = bun
K := $(foreach exec,$(DEPS),\
        $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

help: ## 🛟  Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m 👉 %-14s\033[0m %s\n", $$1, $$2}'


slangroom-exec: $(SOURCES) $(LIBS) ## 🛠️  Build slangroom-exec
	bun build ./src/index.ts --compile --outfile slangroom-exec

clean: ## 🧹 Clean the build
	@rm -f slangroom-exec
	@echo "🧹 Cleaned the build"

tests: slangroom-exec ## 🧪 Run tests
	./test/bats/bin/bats -j 15 test/*.bats
	bun test --coverage

$(LIBS): package.json
	bun i

video:
	PATH=docs:$$PATH
	cd docs && vhs slangroom-exec.tape
