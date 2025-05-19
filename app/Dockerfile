# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app

# Run stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /go-app /go-app
EXPOSE 8080
ENTRYPOINT ["/go-app"]