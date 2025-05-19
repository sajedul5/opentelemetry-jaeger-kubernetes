# Stage 1: Build the Go app with compatible libraries
FROM golang:1.21-bullseye AS builder

WORKDIR /app
COPY . .

# Build the Go app
RUN go build -o vote-app main.go

# Stage 2: Minimal runtime image
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built binary and templates
COPY --from=builder /app/vote-app .
COPY templates ./templates

EXPOSE 8080

CMD ["./vote-app"]

