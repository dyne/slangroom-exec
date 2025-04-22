FROM debian:12-slim

COPY src/slexfe /usr/local/bin/
COPY Makefile Makefile
COPY src /src
COPY bun.lockb bun.lockb
COPY .tool-versions .tool-versions
COPY package.json package.json
COPY tsconfig.json tsconfig.json

# mise
RUN apt-get update \
    && apt-get -y --no-install-recommends install \
		sudo curl git ca-certificates build-essential \
	&& rm -rf /var/lib/apt/lists/*

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="/mise/shims:$PATH"
RUN curl https://mise.run | sh
RUN mise trust
RUN --mount=type=cache,target=/mise/cache mise i

RUN make slangroom-exec
RUN cp slangroom-exec /usr/local/bin/slangroom-exec

ENTRYPOINT ["/usr/local/bin/slangroom-exec"]
