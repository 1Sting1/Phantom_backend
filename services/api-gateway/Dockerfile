FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy pkg module first
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY pkg/ ./pkg/

# Initialize pkg module dependencies
WORKDIR /app/pkg
RUN go mod download && go mod tidy

WORKDIR /app/services/api-gateway
COPY services/api-gateway/go.mod services/api-gateway/go.sum* ./

# Download dependencies
RUN go mod tidy

# Copy rest of service files
COPY services/api-gateway/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

