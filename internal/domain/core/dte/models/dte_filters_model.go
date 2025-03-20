package models

import "time"

type DTEFilters struct {
	StartDate    *time.Time `query:"startDate,omitempty"`
	EndDate      *time.Time `query:"endDate,omitempty"`
	Status       string     `query:"status,omitempty"`
	Transmission string     `query:"transmission,omitempty"`
	BranchID     string     `query:"-"`
	DTEType      string     `query:"-"`

	// Paginación
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
