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
# Requires TEAM_ID environment variable
TEAM_ID=red go run main.go

# Optionally set TEAM_MEMBERS (comma-separated)
TEAM_ID=red TEAM_MEMBERS="Alice Doe, Bob Smith" go run main.go

# Or run the built binary
TEAM_ID=red ./bin/server
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
docker run -e TEAM_ID=red -p 8090:8090 team-info-server
```

**Multi-platform builds** are enabled in GitHub Actions CI/CD. The workflow always builds images for `linux/amd64` and `linux/arm64` platforms, but only pushes to Docker Hub when you push a Git tag in the year.release format.

**Version tag format:**
- `v{year}.{release}` — Official releases (e.g., `v2026.0`, `v2026.1`)
  - `{year}`: Current year (e.g., `2026`)
  - `{release}`: Release sequence in that year (0 = first release, 1 = second release, etc.)

**Release candidate format:**
- `v{year}.{release}-rc{i}` — Release candidates (e.g., `v2026.0-rc1`, `v2026.0-rc2`)
  - `{i}`: Release candidate number

**Behavior:**
- On commit push: Builds image locally but does not push
- On release tag (e.g., `v2026.0`): Builds and pushes to Docker Hub with tags:
  - Version tag (e.g., `2026.0`)
  - `latest`
  - `sha-abc123def` (commit reference)
- On release candidate tag (e.g., `v2026.0-rc1`): Builds and pushes to Docker Hub with tags:
  - Version tag (e.g., `2026.0-rc1`)
  - `sha-abc123def` (commit reference)
  - **Note:** `latest` tag is NOT updated for release candidates

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
  - `team`: The team identifier from the TEAM_ID env var
  - `members`: Team member names from the TEAM_MEMBERS env var (array, may be empty)
  - `k8s_labels`: Kubernetes labels as a map
  - `git_repo`: Constructed GitHub URL
  - `docker_repos`: Constructed Docker Hub URLs (array)

Note: the env var keys use the `TEAM_` prefix (`TEAM_ID`, `TEAM_MEMBERS`), but the JSON response keys remain `team` and `members` (consumed by the downstream ESIGELEC About page).

**Startup Flow**:
1. Validates required `TEAM_ID` environment variable (exits if missing - intentional for teaching)
2. Parses optional `TEAM_MEMBERS` environment variable (comma-separated)
3. Sets up HTTP handler
4. Starts server on port 8090

**Environment Variables**:
- `TEAM_ID` (required): Team identifier, validated against the format below. Missing/invalid value exits the process.
- `TEAM_MEMBERS` (optional): Comma-separated team member names (e.g. `"Alice Doe, Bob Smith"`). Whitespace around each name is trimmed and empty entries are dropped. When unset, `members` is an empty array (`[]`), never `null`, so the field is always safe to treat as an array. Parsed by `parseMembers` in `main.go`.

**Valid TEAM_ID Values**:
- A single lowercase word — letters `a`–`z` only, no digits, hyphens, underscores, or spaces (e.g., `red`, `black`, `green`). Validated by `teamPattern` (`^[a-z]+$`) in `main.go`.
- `teacher` is used for instructor-only access. It is no longer special-cased in code (it passes the general rule like any other word) and is intentionally not documented in the README.
- The previous 2025 `{region}-{digit}` format (e.g., `east-1`) is no longer accepted.

### Key Design Decisions

- **Lightweight**: Uses only Go standard library, no external dependencies
- **Scratch Image**: Docker image uses `scratch` base for minimal size, not suitable for debugging/shell access
- **Port 8090**: Avoids conflict with Spring PetClinic on port 8080
- **Intentional Error**: Missing TEAM_ID env var causes exit - students must fix this in Kubernetes manifests

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

**Docker Hub description sync:** On a published tag, the workflow runs `peter-evans/dockerhub-description` to upload `README.md` as the Docker Hub full description for `mincongclassroom/team-info-server` (and set a short description). This keeps the registry page and the GitHub README in sync so students find the same configuration docs on either side. It reuses the existing `DOCKER_USERNAME` / `DOCKER_PASSWORD` secrets and only runs when `steps.publish.outputs.should_push == 'true'`.

**OCI image labels:** The `Dockerfile` declares `org.opencontainers.image.*` labels (source, documentation, url, description, licenses) so the published image points back to this repository. Registries surface `org.opencontainers.image.source` as the linked source repo. In CI, `docker/metadata-action` also injects dynamic labels; declaring the static ones in the `Dockerfile` keeps local `docker build` output self-describing too.
