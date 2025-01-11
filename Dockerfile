FROM golang:alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy the Go modules and dependencies first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o noted

FROM alpine:latest
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/noted .

# Add an example.env file for the user
COPY example.env /root/example.env

# Command to run the CLI tool
ENTRYPOINT ["./noted"]