.PHONY: help

PLATFORMS = linux-x64 linux-arm64 linux-aarch64 windows-x64 darwin-x64 darwin-arm64
SOURCES = $(shell find src -type f -name '*.ts')
LIBS = node_modules

DEPS = bun
K := $(foreach exec,$(DEPS),\
        $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

help: ## 🛟  Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m 👉 %-14s\033[0m %s\n", $$1, $$2}'


all: $(addprefix slangroom-exec-, $(PLATFORMS)) ## 🛠️  Build all platforms

slangroom-exec: $(SOURCES) $(LIBS) ## 🚀 Build slangroom-exec for the current platform
	bun build ./src/index.ts --compile --minify --outfile slangroom-exec

slangroom-exec-%: $(SOURCES) $(LIBS)
	bun build ./src/index.ts --compile --minify --target=bun-$*-modern --outfile slangroom-exec-$*

clean: ## 🧹 Clean the build
	@rm -f $(addprefix slangroom-exec-, $(PLATFORMS))
	@rm -f slangroom-exec
	@make -C bindings/go clean
	@echo "🧹 Cleaned the build"

tests: slangroom-exec ## 🧪 Run tests
ifeq ($(OS),Windows_NT)
	./test/bats/bin/bats test/*.bats
else
	./test/bats/bin/bats -j 15 test/*.bats
endif
	bun test --coverage

$(LIBS): package.json
	bun i

video:
	PATH=docs:$$PATH
	cd docs && vhs slangroom-exec.tape
