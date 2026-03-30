module Phantom_backend/services/landing-service

go 1.21

replace Phantom_backend/pkg => ../../pkg

require (
	Phantom_backend/pkg v0.0.0
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	go.uber.org/zap v1.26.0
)

require (
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)
