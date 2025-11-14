-include .env
export

DOCKER_COMPOSE_DIR = ./deployments/docker-compose
DOCKER_COMPOSE_CONFIG = docker compose -f ${DOCKER_COMPOSE_DIR}/docker-compose.yml -f ${DOCKER_COMPOSE_DIR}/docker-compose.dev.yml
DOCKER_COMPOSE_RESET = --renew-anon-volumes --force-recreate

build: go-build docker-build

test: unit-tests integration-tests

docker-build:
	${DOCKER_COMPOSE_CONFIG} up --build --no-start

docker-run:
	${DOCKER_COMPOSE_CONFIG} up ${DOCKER_COMPOSE_RESET}
 
docker-up:
	${DOCKER_COMPOSE_CONFIG} up -d ${DOCKER_COMPOSE_RESET}

docker-down:
	${DOCKER_COMPOSE_CONFIG} down

go-build: go-fmt go-get
	go build -o ./bin/challenge ./cmd/challenge

go-fmt:
	go fmt ./...

go-get:
	go get ./...

go-run:
	[ -e test.db ] && rm test.db || true ; \
	TLS_CERT_FILE=./certs/cert.pem \
	TLS_KEY_FILE=./certs/key.pem \
	bin/challenge

unit-tests:
	go test -tags=!integration -v ./cmd/... ./internal/...

integration-tests:
	go test -tags=integration -v ./test
