FROM golang:1.22-alpine AS builder
WORKDIR /src

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the binary statically
RUN CGO_ENABLED=0 go build -o lockify ./cmd/main.go

# Final image
FROM alpine:3.18
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /src/lockify ./

USER appuser
EXPOSE 8080
CMD ["./lockify"]
