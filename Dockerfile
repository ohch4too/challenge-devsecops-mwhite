FROM golang:1.16-buster AS build

WORKDIR /go/src/app
COPY . .

RUN go build -o ./bin/challenge ./cmd/challenge

FROM debian:buster

USER root

COPY --from=build /go/src/app/bin/challenge /usr/bin/challenge

ENTRYPOINT [ "/usr/bin/challenge" ]
