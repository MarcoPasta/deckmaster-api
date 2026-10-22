# Multi-stage Dockerfile for building and running a Go application
# - Build stage uses the official Golang image to compile a static binary
# - Final stage uses a minimal distroless image for a small runtime surface

FROM golang:1.27.3-trixie AS builder

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       ca-certificates \
       bash \
    && rm -rf /var/lib/apt/lists/*

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
    go build -ldflags "-s -w" -o /app ./cmd/deckmaster

# Ensure a predictable non-root user exists in the final image. Create a
# minimal passwd entry and a home directory in the builder and copy them
# into the final (distroless) image so tools that expect a named user work.
RUN mkdir -p /home/nonroot \
    && chown 1000:1000 /home/nonroot \
    && printf 'nonroot:x:1000:1000:nonroot:/home/nonroot:/sbin/nologin\n' > /tmp/passwd


# Final image: minimal runtime with shell access
FROM alpine:3.24

RUN adduser -D -u 1000 nonroot

WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app /app

# Copy data files and migrations from the build context (rooted at /src)
COPY --from=builder /src/migrations /migrations

EXPOSE 8080

USER nonroot

ENTRYPOINT ["/app"]
