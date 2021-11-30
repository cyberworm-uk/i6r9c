FROM golang:1.17-alpine AS build
ENV CGO_ENABLED=0
WORKDIR /go/src
COPY . .
RUN go mod download
WORKDIR /go/src/cmd
RUN go build -o irc .
FROM gcr.io/distroless/static
COPY --from=build /go/src/cmd/irc /irc
ENTRYPOINT [ "/irc" ]