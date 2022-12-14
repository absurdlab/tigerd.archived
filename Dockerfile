FROM bufbuild/buf:1.10.0 AS buf

FROM golang:1.18.1 AS buildStage

ENV HOME /build
ENV CGO_ENABLED 0
ENV GOOS linux

ARG GO_PROXY=https://goproxy.io,direct

COPY --from=buf /usr/local/bin/buf /usr/local/bin/buf

WORKDIR /build

COPY go.mod go.sum ./
RUN go env -w GOPROXY=$GO_PROXY && \
    go mod download -x

COPY . .
RUN go generate && \
    go build \
    -a \
    -ldflags "-X github.com/absurdlab/tigerd/buildinfo.GitHash=`git rev-parse HEAD | head -c 8`" \
    -installsuffix cgo \
    -o tigerd \
    .

FROM alpine:3.15

LABEL org.opencontainers.image.title="tigerd"
LABEL org.opencontainers.image.source="https://github.com/absurdlab/tigerd"
LABEL org.opencontainers.image.authors="Weinan Qiu"

COPY --from=buildStage /build/tigerd /usr/bin/tigerd

ENTRYPOINT ["/usr/bin/tigerd"]
