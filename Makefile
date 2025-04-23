# SPDX-FileCopyrightText: 2024-2025 Dyne.org foundation
#
# SPDX-License-Identifier: AGPL-3.0-or-later

.PHONY: help

PLATFORMS = linux-x64 linux-arm64 linux-aarch64 windows-x64 darwin-x64 darwin-arm64
BINARIES_NO_EXT := $(addprefix slangroom-exec-, $(PLATFORMS))
ARCHIVES := $(addsuffix .tar.gz, $(BINARIES_NO_EXT))
BINARIES := $(foreach binary, $(BINARIES_NO_EXT), ${binary}$(if $(filter slangroom-exec-windows-%,$(binary)),.exe))
SOURCES = $(shell find src -type f -name '*.ts')
LIBS = node_modules

DEPS = bun
K := $(foreach exec,$(DEPS),\
        $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

help: ## 🛟  Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m 👉 %-14s\033[0m %s\n", $$1, $$2}'


all: $(BINARIES_NO_EXT) ## 🛠️  Build all platforms

slangroom-exec: $(SOURCES) $(LIBS) ## 🚀 Build slangroom-exec for the current platform
	bun build ./src/index.ts --compile --minify --outfile slangroom-exec

update: ## 🌍 Update slangroom to latest release
	bun update @slangroom/core@latest

slangroom-exec-%: $(SOURCES) $(LIBS)
	bun build ./src/index.ts --compile --minify --target=bun-$*-modern --outfile slangroom-exec-$*

clean: ## 🧹 Clean the build
	@rm -f $(BINARIES)
	@rm -f $(ARCHIVES)
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

archives: $(BINARIES_NO_EXT) $(ARCHIVES) ## 📦 Create archives containing slangroom-exec and slexfe for all platforms

slangroom-exec-%.tar.gz:
	tar -czf $@ \
		--transform='s|slangroom-exec-$*$(if $(filter windows-%,$*),\.exe)|slangroom-exec|' \
		slangroom-exec-$*$(if $(filter windows-%,$*),.exe) \
		-C src slexfe
