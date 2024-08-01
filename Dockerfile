# Stage 1: Build
FROM golang:1.20-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make bash

# Create and set working directory
WORKDIR /go/src/sekretariat

# Copy source code and build
COPY . .
RUN make build

# Stage 2: Final image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Expose the port
EXPOSE 8080

# Copy the built binary and configuration files
COPY --from=builder /go/src/sekretariat/bin/sekretariat /sekretariat
COPY --from=builder /go/src/sekretariat/files/etc/sekretariat/sekretariat.development.yaml /sekretariat.development.yaml

# Set the entrypoint
ENTRYPOINT ["/sekretariat"]
