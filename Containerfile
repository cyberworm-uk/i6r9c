FROM docker.io/alpine/git:latest AS source
ARG VERSION=main
WORKDIR /go/src
RUN git clone --depth=1 --branch=${VERSION} https://github.com/guest42069/i6r9c/ .

FROM docker.io/library/golang:1.17-alpine AS build
COPY --from=source /go/src /go/src
WORKDIR /go/src/cmd
RUN go mod download
RUN if [[ "`go env | grep "^GOARCH=" | sed 's:^GOARCH="\(.*\)"$:\1:'`" != "arm" ]]; then CGO_ENABLED=0 go build -ldflags '-w -s -buildid=' -trimpath -buildmode=pie .;else CGO_ENABLED=0 go build -ldflags '-w -s -buildid=' -trimpath .;fi

FROM docker.io/library/alpine:latest
RUN apk -U upgrade --no-cache
RUN apk add ca-certificates tzdata --no-cache
COPY --from=build /go/src/cmd/cmd /i6r9c
USER 1000
ENTRYPOINT [ "/i6r9c" ]
