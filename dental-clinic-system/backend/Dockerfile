# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o i-dentist-api .

# Final stage
FROM alpine:1.22

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/i-dentist-api .

# Copy configuration files if needed
COPY --from=builder /app/resources ./resources

EXPOSE 8080

CMD ["./i-dentist-api"]