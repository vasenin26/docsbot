# docsbot - metrics server

This repository contains a minimal metrics and health HTTP server for the docsbot project.
The server exposes Prometheus metrics and health endpoints only — no business logic.

Usage

Build locally:

    make build

Run locally:

    make run

The server listens on port 9090 by default. You can override the port via METRICS_PORT environment variable.

Endpoints

- GET /healthz - returns 200 OK with {"status":"ok"}
- GET /readyz - same as /healthz (readiness)
- GET /metrics - Prometheus metrics endpoint

Docker

Build image (host architecture):

    make docker-build

Build multi-arch image (requires docker buildx and a registry):

    make docker-buildx PLATFORMS=linux/amd64,linux/arm64 TAG=yourrepo/docsbot:tag

Testing

Run unit tests:

    make test

Notes

- The project uses github.com/prometheus/client_golang for metrics.
- Dockerfile is multi-stage; final image is based on distroless static nonroot.
- If you need a debuggable image, change the final image to alpine in the Dockerfile.
