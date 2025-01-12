### First Stage
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o noted .

### Second Stage
FROM alpine:latest

# Use /workdir instead of /app to make it clearer this is for mounted files
WORKDIR /workdir
COPY --from=builder /app/noted /usr/local/bin/noted

# Use full path in ENTRYPOINT
ENTRYPOINT ["/usr/local/bin/noted"]

