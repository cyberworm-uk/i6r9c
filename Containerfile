FROM --platform=$BUILDPLATFORM docker.io/library/golang:alpine AS build
ARG TARGETOS TARGETARCH
ENV GOOS="$TARGETOS" GOARCH="$TARGETARCH" GOFLAGS="-buildvcs=false -trimpath" CGO_ENABLED=0
WORKDIR /go/src
COPY . .
RUN --mount=type=cache,target=/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build go build -o /i6r9c -ldflags '-w -s -buildid=' ./cmd/i6r9c

FROM ghcr.io/cyberworm-uk/base:latest
RUN apk add ca-certificates tzdata --no-cache
COPY --from=build /i6r9c /i6r9c
USER 1000
ENTRYPOINT [ "/i6r9c" ]