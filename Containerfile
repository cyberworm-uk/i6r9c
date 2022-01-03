FROM docker.io/alpine/git:latest AS source
ARG VERSION=main
WORKDIR /go/src
RUN git clone --depth=1 --branch=${VERSION} https://github.com/guest42069/i6r9c/ .

FROM docker.io/library/golang:1.17-alpine AS build
COPY --from=source /go/src /go/src
WORKDIR /go/src/cmd
RUN go mod download
RUN if [[ "`go env | grep "^GOARCH=" | sed 's:GOARCH="\(.*\)":\1:'`" != "arm" ]]; then export PIE="-buildmode=pie"; fi
RUN CGO_ENABLED=0 go build -ldflags '-w -s -buildid=' -trimpath $PIE .

FROM docker.io/library/alpine:latest AS files
RUN apk -U upgrade --no-cache
RUN apk add ca-certificates tzdata --no-cache

FROM scratch
COPY --from=files /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=files /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /go/src/cmd/cmd /i6r9c
USER 1000
ENTRYPOINT [ "/i6r9c" ]
