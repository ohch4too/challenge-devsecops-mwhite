# DevSecOps Challenge

*Mike White - mw419@pm.me*

## Project Summary

This project demonstrates progressive development work through four pull requests. Each PR addresses specific concerns: establishing project structure and CI/CD automation, implementing REST API compliance with proper HTTP standards, hardening security through credential management and vulnerability remediation, and refactoring to clean architecture with comprehensive testing. The current state is not a complete product but a snapshot of work in progress. Future development might include authentication mechanisms, authorization controls, API versioning, observability features, and business logic that transforms this from a demo into a functional microservice with defined purpose.

Detailed summaries of each development phase are available in the `WORKSUMMARY/` directory, with individual markdown files corresponding to each pull request.

*Note: Technical implementation and narrative content within this project were developed with assistance from generative AI tools.*

## Configuration Requirements

### Local Development Secrets

Before running the application locally, configure the following secrets:

**TLS Certificates**: Create a `certs/` directory in the project root and place two files inside:
- `cert.pem` - TLS certificate
- `key.pem` - Private key

Generate a self-signed certificate for local development:

```bash
openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -out cert.pem -days 365 -subj "/CN=localhost"
```

**Environment Variables**: Copy `.env.example` to `.env` and edit with your database passwords and other configuration values.

### GitHub CI Secrets

For the GitHub Actions pipeline to run successfully, configure these base64-encoded secrets in your repository settings. See the [GitHub documentation for creating encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets).

Required secret variables:
- `KEY_PEM_BASE64` - Base64-encoded private key
- `CERT_PEM_BASE64` - Base64-encoded certificate
- `ENV_FILE_BASE64` - Base64-encoded environment file

Convert file content to base64 format:

```bash
cat .env | base64 -w0
```

Use this same command pattern for the certificate and key files, then add each encoded value as a secret in your GitHub repository settings.

## Building and Running

The project uses a Makefile to standardize build and execution commands. All commands should be run from the repository root.

### Build Containers

Build both the Go binary and Docker containers:

```bash
make build
```

This runs `go build` to create the binary in `bin/challenge` and builds Docker images defined in the compose configuration.

### Run with Docker Compose

Start the application and PostgreSQL database:

```bash
make docker-up
```

This starts containers in detached mode. The API serves on port 10000 with HTTPS if TLS certificates are configured via environment variables.

Stop containers:

```bash
make docker-down
```

### Run the App Directly

Run the Go binary directly without containers:

```bash
make go-run
```

This requires a local database or uses SQLite. TLS certificate paths are configured through environment variables.

### Testing

Run all tests (unit and integration):

```bash
make test
```

Run only unit tests:

```bash
make unit-tests
```

Run only integration tests (requires Docker containers):

```bash
make integration-tests
```

Unit tests validate individual components using mocks. Integration tests verify the complete stack including database interactions.

## CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/main.yml`) runs on pushes to `workspace` branch and pull requests to `main` branch. The pipeline contains two parallel jobs:

**build-run-test**: Builds the application and containers, starts the Docker environment, executes all tests, then stops containers. This validates functional correctness.

**build-scan**: Performs security scanning with three tools:
- SAST (gosec) scans Go code for security issues
- SCA (govulncheck) checks dependencies for known vulnerabilities  
- Container scanning (Trivy) analyzes Docker images for vulnerabilities

Both jobs require secrets (TLS certificates and environment configuration) injected from GitHub repository secrets. The workflow ensures code quality and security before changes are merged.

