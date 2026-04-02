# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**team-info-server** is a lightweight Go HTTP server that serves team metadata in JSON format. It's designed as a teaching tool for students to learn about environment variables and Kubernetes configuration.

- **Language**: Go 1.26
- **Primary File**: `main.go` (single-file project)
- **Port**: 8090
- **Docker**: Multi-stage build with scratch base image

## Common Commands

### Build
```bash
go build -o bin/server ./main.go
# or with default output
go build
```

### Run
```bash
# Requires TEAM environment variable
TEAM=east-1 go run main.go

# Or run the built binary
TEAM=east-1 ./bin/server
```

### Test
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestTeamInfoHandler
```

### Docker
```bash
# Build image
docker build -t team-info-server .

# Run container
docker run -e TEAM=east-1 -p 8090:8090 team-info-server
```

**Multi-platform builds** are enabled in GitHub Actions CI/CD. The workflow always builds images for `linux/amd64` and `linux/arm64` platforms, but only pushes to Docker Hub when you push a Git tag in semantic versioning format.

**Supported tag formats:**
- `v1.0.0` (release)
- `v1.0.0-alpha` (pre-release)
- `v1.0.0-alpha-1`, `v1.0.0-beta.2`, `v1.0.0-rc.1` (pre-release variants)

**Behavior:**
- On commit push: Builds image locally but does not push
- On version tag (e.g., `v1.2.3` or `v1.2.3-alpha-1`): Builds and pushes to Docker Hub with tags:
  - Version tag (e.g., `1.2.3` or `1.2.3-alpha-1`)
  - `latest`
  - `sha-abc123def` (commit reference)

## Testing

The project includes unit tests in `main_test.go` that verify:
- HTTP handler returns correct JSON response
- Different team values are correctly formatted in git and docker repo URLs
- K8s labels are properly set
- JSON marshaling/unmarshaling works correctly

Tests use table-driven testing pattern to cover multiple team values.

## Architecture

### Core Components

**HTTP Server Structure**:
- Single root handler (`/`) that accepts all HTTP methods
- Returns JSON response with `TeamInfo` struct containing:
  - `team`: The team identifier from TEAM env var
  - `k8s_labels`: Kubernetes labels as a map
  - `git_repo`: Constructed GitHub URL
  - `docker_repo`: Constructed Docker Hub URL

**Startup Flow**:
1. Validates required `TEAM` environment variable (exits if missing - intentional for teaching)
2. Sets up HTTP handler
3. Starts server on port 8090

### Key Design Decisions

- **Lightweight**: Uses only Go standard library, no external dependencies
- **Scratch Image**: Docker image uses `scratch` base for minimal size, not suitable for debugging/shell access
- **Port 8090**: Avoids conflict with Spring PetClinic on port 8080
- **Intentional Error**: Missing TEAM env var causes exit - students must fix this in Kubernetes manifests

## Development Notes

- The project is intentionally simple to focus on Kubernetes/DevOps concepts
- All application logic is in `main.go`
- The HTTP handler logs all requests with method and path
- Response always uses JSON content-type header
- `json.NewEncoder()` is used for streaming JSON encoding

## Go Version Updates

When updating the Go version, update all of these files:
1. `go.mod` — Module Go version directive
2. `Dockerfile` — Builder stage base image (`golang:X.XX-alpine`)
3. `.github/workflows/build.yml` — CI/CD setup-go action version
4. `CLAUDE.md` — This documentation

## GitHub Actions Workflow Notes

**Multi-platform builds:** The workflow uses `docker/setup-buildx-action@v3` to enable multi-platform Docker builds for linux/amd64 and linux/arm64. This step must come before `docker/build-push-action` to support multi-platform builds (the default docker driver does not support this).

**Docker build caching:** Do not add `cache-from` or `cache-to` options to `docker/build-push-action`. The GitHub Actions docker driver does not support cache export and raises "Cache export is not supported for the docker driver" error. Build optimization is not a priority at this stage.
