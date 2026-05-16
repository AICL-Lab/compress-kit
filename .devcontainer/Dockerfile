FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y --no-install-recommends build-essential clang-format python3 python3-pip curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ARG GO_VERSION=1.21.13
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"

USER vscode
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y --profile minimal --component clippy rustfmt
ENV PATH="/home/vscode/.cargo/bin:${PATH}"
