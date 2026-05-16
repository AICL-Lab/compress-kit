FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y --no-install-recommends build-essential clang-format python3 python3-pip curl ca-certificates xz-utils \
    && rm -rf /var/lib/apt/lists/*

ARG NODE_VERSION=20.20.2
RUN arch="$(dpkg --print-architecture)" \
    && case "${arch}" in \
        amd64) node_arch='x64' ;; \
        arm64) node_arch='arm64' ;; \
        *) echo "Unsupported architecture: ${arch}" >&2; exit 1 ;; \
    esac \
    && curl -fsSL "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${node_arch}.tar.xz" \
    | tar -xJ --strip-components=1 -C /usr/local

ARG GO_VERSION=1.21.13
RUN arch="$(dpkg --print-architecture)" \
    && case "${arch}" in \
        amd64) go_arch='amd64' ;; \
        arm64) go_arch='arm64' ;; \
        *) echo "Unsupported architecture: ${arch}" >&2; exit 1 ;; \
    esac \
    && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${go_arch}.tar.gz" | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"

USER vscode
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y --profile minimal --component clippy rustfmt
ENV PATH="/home/vscode/.cargo/bin:${PATH}"
