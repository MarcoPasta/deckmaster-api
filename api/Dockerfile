# Multi-stage Dockerfile for building and running a Go application
# - Build stage uses the official Golang image to compile a static binary
# - Final stage uses a minimal distroless image for a small runtime surface

FROM golang:1.25-bookworm AS builder

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

# Ensure a predictable non-root user exists in the final image. Create a
# minimal passwd entry and a home directory in the builder and copy them
# into the final (distroless) image so tools that expect a named user work.
RUN mkdir -p /home/nonroot \
    && chown 1000:1000 /home/nonroot \
    && printf 'nonroot:x:1000:1000:nonroot:/home/nonroot:/sbin/nologin\n' > /tmp/passwd


# Final image: minimal runtime (non-root)
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app /app

# Copy the passwd entry and home directory so the `nonroot` user exists by
# name inside the final image (some tools / extensions lookup users by name).
COPY --from=builder /tmp/passwd /etc/passwd
COPY --from=builder /home/nonroot /home/nonroot

# Expose the default port. Your app can still read a different port from
# environment variables if it supports that (e.g. PORT).
EXPOSE 8080

USER nonroot

ENTRYPOINT ["/app"]
