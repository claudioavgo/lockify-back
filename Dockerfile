FROM golang:1.24.4-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o lockify ./cmd/main.go

FROM alpine:3.18
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /src/lockify ./

USER appuser
EXPOSE 8080
CMD ["./lockify"]
