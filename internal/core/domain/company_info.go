package domain

import "github.com/google/uuid"

type CompanyInfo struct {
	UserID              uuid.UUID `json:"user_id"`
	CompanyName         *string   `json:"company_name,omitempty"`
	CompanySize         *string   `json:"company_size,omitempty"`
	CompanyIndustry     *string   `json:"company_industry,omitempty"`
	CompanyDescription  *string   `json:"company_description,omitempty"`
	CompanyHeadquarters *string   `json:"company_headquarters,omitempty"`
	CompanyWebsite      *string   `json:"company_website,omitempty"`
	AccountType         string  `json:"account_type,omitempty"`
}
