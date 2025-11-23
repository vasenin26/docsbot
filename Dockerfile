# ---- builder ----
FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /src

# Copy modules and download dependencies early for better caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary. When using buildx this will cross-compile for target platforms.
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -ldflags "-s -w" -o /out/docsbot ./cmd/docsbot

# ---- final ----
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/docsbot /docsbot
EXPOSE 9090
USER nonroot
ENTRYPOINT ["/docsbot"]
