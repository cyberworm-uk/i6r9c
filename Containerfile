FROM docker.io/alpine/git:latest AS source
ARG VERSION=main
WORKDIR /go/src
RUN git clone --depth=1 --branch=${VERSION} https://github.com/guest42069/i6r9c/ .

FROM docker.io/library/golang:1.17-alpine AS build
COPY --from=source /go/src /go/src
ENV CGO_ENABLED=0
WORKDIR /go/src
RUN go mod download
WORKDIR /go/src/cmd
RUN go build -o irc .

FROM docker.io/library/alpine:latest
RUN apk -U upgrade --no-cache
RUN apk add ca-certificates tzdata --no-cache
COPY --from=build /go/src/cmd/irc /irc
USER 1000
ENTRYPOINT [ "/irc" ]
