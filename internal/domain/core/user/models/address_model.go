package models

// Address representa la dirección de una sucursal o casa matriz
type Address struct {
	ID           uint   `json:"id,omitempty"`
	BranchID     uint   `json:"branch_id"`
	Municipality string `json:"municipality"`
	Department   string `json:"department"`
	Complement   string `json:"complement"`
}
