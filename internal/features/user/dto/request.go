package userdto

import (
	"errors"
	"strings"
)

// UpdateProfileRequest is the request body for updating a user profile.
type UpdateProfileRequest struct {
	FirstName          string `json:"first_name" binding:"required"`
	LastName           string `json:"last_name" binding:"required"`
	Nationality        string `json:"nationality" binding:"required"`
	Gender             string `json:"gender"`
	ResidentialCountry string `json:"residential_country"`
	JobRole            string `json:"job_role"`
	EmploymentType     string `json:"employment_type"`
}

func (r *UpdateProfileRequest) Validate() error {
	if strings.TrimSpace(r.FirstName) == "" || strings.TrimSpace(r.LastName) == "" {
		return errors.New("first name and last name cannot be empty")
	}
	if strings.TrimSpace(r.Nationality) == "" {
		return errors.New("nationality cannot be empty")
	}
	return nil
}

// UpdateUserPasswordRequest is the request body for changing a password.
type UpdateUserPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (r *UpdateUserPasswordRequest) Validate() error {
	if len(r.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}
	return nil
}
