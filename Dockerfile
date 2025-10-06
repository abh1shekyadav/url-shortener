FROM golang:1.24-bullseye AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=arm64 go build -o url-shortener .
FROM debian:bullseye-slim
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=builder /app/url-shortener .
EXPOSE 8080
CMD ["./url-shortener"]