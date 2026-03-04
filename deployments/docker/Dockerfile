# ---------- Build stage ----------
FROM golang:1.25.7-alpine AS builder

# Install development tools
RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@latest

WORKDIR /app

# Install git (for go modules)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod ./
# Install dependencies (if any)
# RUN go mod init github.com/james-wukong/online-orders && go mod tidy

RUN if [ -f go.mod ]; then go mod download && go mod tidy && go mod vendor; fi

# Copy source code
COPY . .

# Build the binary
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
#     go build -o app ./cmd/api

# # ---------- Runtime stage ----------
# FROM alpine:3.19

# WORKDIR /app

# # Create non-root user
# RUN adduser -D -g '' appuser
# USER appuser

# # Copy binary from builder
# COPY --from=builder /app/app .

# Expose HTTP port
EXPOSE 8092

# Run
# CMD ["./app"]
# air will watch for file changes and rebuild/restart
CMD ["air", "-c", ".air.toml"]
