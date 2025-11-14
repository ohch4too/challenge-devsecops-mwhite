FROM golang:1.25-trixie AS build

WORKDIR /go/src/app
COPY . .

RUN go build -o ./bin/challenge ./cmd/challenge

FROM debian:trixie

COPY --from=build /go/src/app/bin/challenge /usr/bin/challenge

ENTRYPOINT [ "/usr/bin/challenge" ]
