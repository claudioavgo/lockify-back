FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gymapp-back ./cmd/main.go

EXPOSE 8080

CMD ["./gymapp-back"]