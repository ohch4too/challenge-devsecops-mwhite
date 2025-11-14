## PR#1: Project Foundation and Structure

This pull request establishes a foundation for the application by implementing Go best practices and DevOps automation. The restructuring addresses critical organizational and build issues that hindered development and deployment. By adopting standard Go project layout, the codebase is maintainable and accessible to new developers. The addition of CI/CD automation ensures code quality through continuous testing, while the docker-compose improvements enable reliable local development and testing environments. These changes transform the project from a demo into a properly structured application with automated quality gates, setting the stage for future security enhancements and feature development.

### Standard Go Project Layout

Reorganized code into 'challenge' for application entry point and internal for private packages. This structure improves maintainability and follows Go community standards. The internal directory enforces package privacy, preventing external imports. This separation enables better code organization and easier navigation for developers familiar with Go conventions.

```go
// cmd/challenge/main.go
import (
    "challenge/internal/api"
    "challenge/internal/db"
    "challenge/internal/dummy"
)
```

### Dockerfile Build Path Correction

Fixed the build command to reference the new challenge location and optimized the final image using debian base instead of golang. This change was necessary after restructuring and reduces image size by removing unnecessary Go toolchain from runtime. Smaller images improve security by minimizing attack surface and reduce deployment time.

```dockerfile
RUN go build -o ./bin/challenge ./cmd/challenge

FROM debian:buster
```

### GitHub Actions CI/CD Workflow

Added an automated pipeline triggered on pull requests to validate builds and run tests. This automation catches build failures and test regressions early in development, ensuring consistent build environment and preventing integration issues. The workflow provides confidence that code changes work correctly before merge.

```yaml
jobs:
  build-and-run:
    runs-on: ubuntu-latest
    steps:
      - name: Build App and Container
        run: make build
      - name: Run Tests
        run: make test
```

### Enhanced Makefile

Consolidated build and test commands with proper configuration paths and environment management flags. The Makefile provides standardized commands that work across development and CI environments. The `DOCKER_COMPOSE_RESET` flags ensure clean test runs by recreating containers, preventing state pollution between executions and eliminating environment-specific variations.

```makefile
DOCKER_COMPOSE_CONFIG = docker compose -f ${DOCKER_COMPOSE_DIR}/docker-compose.yml -f ${DOCKER_COMPOSE_DIR}/docker-compose.dev.yml

integration-tests:
    go test -tags=integration -v ./test
```

### Docker Compose Configuration Split

Separate base and development configurations with health checks and dependency management were added to prevent startup race conditions. The split enables different configurations for 'production' and local development while sharing common settings. Health checks ensure postgres is ready before application starts, eliminating connection failures and improving deployment reliability.

```yaml
depends_on:
  postgres:
    condition: service_healthy
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U challenge"]
```

### Integration Test Framework

Created integration tests validating service availability and endpoint response with retry logic to handle container startup timing. The Makefile targets separats integration tests from unit tests for selective execution. This test validates the complete application stack including database connectivity, catching issues that unit tests miss.

```go
func waitForAPI(url string, timeout time.Duration) error {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        resp, err := http.Get(url)
        if err == nil {
            return nil
        }
        time.Sleep(retryDelay)
    }
    return fmt.Errorf("API not available after %v", timeout)
}
```

### API Response Status Fix

Corrected `RespondJSON` function to an return actual status code instead of hardcoded 200. This bug caused all responses to report success even during errors, breaking REST semantics. The fix ensures error conditions return appropriate status codes, enabling correct client-side error handling and improving API debugging.

```go
func RespondJSON(w *gin.Context, status int, payload interface{}) {
    res.Status = status
    w.JSON(status, res)  // Changed from w.JSON(200, res)
}
```