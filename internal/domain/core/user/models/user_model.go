package models

import "time"

type User struct {
	ID             uint      `json:"id,omitempty"`
	NIT            string    `json:"nit"`
	NRC            string    `json:"nrc"`
	Status         bool      `json:"status"`
	AuthType       string    `json:"auth_type"`
	PasswordPri    string    `json:"password_pri"`
	CommercialName string    `json:"commercial_name"`
	Business       string    `json:"business_name"`
	Email          string    `json:"email"`
	YearInDTE      bool      `json:"year_in_dte"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}
