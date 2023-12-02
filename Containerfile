FROM --platform=$BUILDPLATFORM docker.io/alpine/git:latest AS source
ARG VERSION=main
WORKDIR /go/src
RUN git clone --depth=1 --branch=${VERSION} https://github.com/cyberworm-uk/i6r9c/ .

FROM --platform=$BUILDPLATFORM docker.io/library/golang:alpine AS build
ARG TARGETOS TARGETARCH
ENV GOOS="$TARGETOS" GOARCH="$TARGETARCH" GOFLAGS="-buildvcs=false -trimpath" CGO_ENABLED=0 
COPY --from=source /go/src /go/src
WORKDIR /go/src/cmd
RUN --mount=type=cache,target=/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build go build -ldflags '-w -s -buildid=' .

FROM ghcr.io/cyberworm-uk/base:latest
RUN apk add ca-certificates tzdata --no-cache
COPY --from=build /go/src/cmd/cmd /i6r9c
USER 1000
ENTRYPOINT [ "/i6r9c" ]
