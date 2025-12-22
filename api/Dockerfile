# Multi-stage Dockerfile for building and running a Go application
# - Build stage uses the official Golang image to compile a static binary
# - Final stage uses a minimal distroless image for a small runtime surface

FROM golang:1.21-bullseye AS builder

WORKDIR /src

# Cache go modules
# Only copy `go.mod` so the builder can download dependencies and generate
# `go.sum` inside the container. Avoid copying an empty `go.sum` from the host.
COPY go.mod ./
# Ensure go.sum is populated inside the build container
RUN go mod download && go mod tidy

# Copy source and build
COPY . .

# Build a static, optimized binary. Adjust the package path (.) if your main
# package is in a subdirectory like ./cmd/server.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags "-s -w" -o /app .


# Final image: minimal runtime (non-root)
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app /app

# Expose the default port. Your app can still read a different port from
# environment variables if it supports that (e.g. PORT).
EXPOSE 8080

USER nonroot

ENTRYPOINT ["/app"]
