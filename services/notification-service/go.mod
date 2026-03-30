module Phantom_backend/services/notification-service

go 1.21

replace Phantom_backend/pkg => ../../pkg

require (
	Phantom_backend/pkg v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/gorilla/mux v1.8.1
	go.uber.org/zap v1.26.0
)

require go.uber.org/multierr v1.11.0 // indirect
