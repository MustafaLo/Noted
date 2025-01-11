### First Stage

# Build and run in a single stage
### First Stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o noted .

### Second Stage
FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/noted .
COPY example.env /app/example.env

ENTRYPOINT ["./noted"]

