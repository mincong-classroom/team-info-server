FROM golang:1.26-alpine AS builder
WORKDIR /workspace
COPY go.mod ./
COPY . ./
RUN go build -o bin/server ./main.go

FROM scratch
# OCI image labels link the published image back to its source and docs. Docker Hub
# and other registries surface `org.opencontainers.image.source` as the linked
# repository. In CI these are also set dynamically by docker/metadata-action, but
# declaring them here keeps local builds self-describing too.
LABEL org.opencontainers.image.title="team-info-server" \
      org.opencontainers.image.description="Serves team metadata (team ID, members, k8s labels, repo URLs) as JSON for the mincong-classroom Kubernetes exercises." \
      org.opencontainers.image.source="https://github.com/mincong-classroom/team-info-server" \
      org.opencontainers.image.documentation="https://github.com/mincong-classroom/team-info-server#readme" \
      org.opencontainers.image.url="https://hub.docker.com/r/mincongclassroom/team-info-server" \
      org.opencontainers.image.licenses="GPL-3.0-only"
COPY --from=builder /workspace/bin/server /server
EXPOSE 8090
ENTRYPOINT ["/server"]
