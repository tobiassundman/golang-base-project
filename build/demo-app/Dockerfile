# Download go dependencies
FROM golang:1.20-alpine AS downloader

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy
RUN go mod download

COPY ./ ./

# Build the go binary
FROM downloader AS builder

RUN CGO_ENABLED=0 go build -a -o application cmd/demo-app/main.go

# Build runnable docker image
FROM alpine:latest  
WORKDIR /app
COPY --from=builder /build/application ./
CMD ["./application"]
EXPOSE 8080
