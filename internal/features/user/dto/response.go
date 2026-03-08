package userdto

import "time"

// SuccessResponse is a generic success response.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
}

// ErrorResponse is a generic error response.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
}

// UserResponse represents user data in a response.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
