# Simple usage with a mounted data directory:
# > docker build -t tgerret/bococoin-node .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.bocod:/root/.bococoin -v ~/.bococli:/root/.bococoin/cli bococoin bocod init --chain-id bcc
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.bocod:/root/.bococoin -v ~/.bococli:/root/.bococoin/cli bococoin bocod start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Set working directory for the build
WORKDIR /go/src/github.com/Bococoin/core

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make linux

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/bocod /usr/bin/bocod
COPY --from=build-env /go/bin/bococli /usr/bin/bococli

# Run gaiad by default, omit entrypoint to ease using container with gaiacli
# CMD ["bocod"]
