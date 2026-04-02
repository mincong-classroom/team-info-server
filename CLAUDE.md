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

### Docker
```bash
# Build image
docker build -t team-info-server .

# Run container
docker run -e TEAM=east-1 -p 8090:8090 team-info-server
```

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
