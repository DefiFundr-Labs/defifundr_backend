#!/bin/bash

# Generate mocks for repositories
counterfeiter -o internal/core/ports/mocks/user_repository.go internal/core/ports UserRepository
counterfeiter -o internal/core/ports/mocks/session_repository.go internal/core/ports SessionRepository
counterfeiter -o internal/core/ports/mocks/otp_repository.go internal/core/ports OTPRepository
counterfeiter -o internal/core/ports/mocks/kyc_repository.go internal/core/ports KYCRepository
counterfeiter -o internal/core/ports/mocks/email_repository.go internal/core/ports EmailRepository

# Generate mocks for services
counterfeiter -o internal/core/ports/mocks/auth_service.go internal/core/ports AuthService
counterfeiter -o internal/core/ports/mocks/user_service.go internal/core/ports UserService

# Generate mocks for token maker
counterfeiter -o internal/core/ports/mocks/token_maker.go pkg/token_maker Maker 