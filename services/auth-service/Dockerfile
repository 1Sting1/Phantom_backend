FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy pkg module first
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY pkg/ ./pkg/

# Initialize pkg module dependencies
WORKDIR /app/pkg
RUN go mod download && go mod tidy

WORKDIR /app/services/auth-service
COPY services/auth-service/go.mod services/auth-service/go.sum* ./

RUN go mod tidy

# Copy rest of service files
COPY services/auth-service/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8001

CMD ["./main"]

