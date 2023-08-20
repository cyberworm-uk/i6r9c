FROM --platform=$BUILDPLATFORM docker.io/alpine/git:latest AS source
ARG VERSION=main
WORKDIR /go/src
RUN git clone --depth=1 --branch=${VERSION} https://github.com/guest42069/i6r9c/ .

FROM --platform=$BUILDPLATFORM docker.io/library/golang:alpine AS build
ARG TARGETOS TARGETARCH
COPY --from=source /go/src /go/src
WORKDIR /go/src/cmd
RUN --mount=type=cache,target=/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -buildvcs=false -ldflags '-w -s -buildid=' -trimpath `if [[ "$TARGETARCH" != "arm" ]]; then echo "-buildmode=pie"; fi` .

FROM docker.io/library/alpine:latest
RUN apk -U upgrade --no-cache
RUN apk add ca-certificates tzdata --no-cache
COPY --from=build /go/src/cmd/cmd /i6r9c
USER 1000
ENTRYPOINT [ "/i6r9c" ]
