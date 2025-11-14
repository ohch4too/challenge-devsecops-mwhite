FROM golang:1.25-trixie AS build

WORKDIR /go/src/app
COPY . .

RUN go build -o ./bin/challenge ./cmd/challenge

FROM debian:trixie

RUN apt-get update && apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /go/src/app/bin/challenge /usr/bin/challenge

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f -k https://localhost:10000/v1/users || exit 1

ENTRYPOINT [ "/usr/bin/challenge" ]
