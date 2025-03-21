package ports

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// DTERepository es una interfaz que define los métodos que debe implementar un repositorio de DTEs
type DTERepository interface {
	// GetDTEByID obtiene un DTE por su ID
	GetDTEByID(string) (*dte.DTEDetails, error)
	// CreateDTE crea un DTE
	CreateDTE(*dte.DTEDetails, uint) error
	// GetAllDTEsByBranchOfficeID obtiene todos los DTEs de una sucursal
	GetAllDTEsByBranchOfficeID(uint, *dte.DTEFilters) ([]dte.DTEResponse, error)
	// UpdateStatus actualiza el estado de un DTE
	UpdateStatus(string, string) error
	// UpdateReceptionStamp actualiza el sello de recepción de un DTE
	UpdateReceptionStamp(string, string) error
	// UpdateTransmissionType actualiza el tipo de transmisión de un DTE (Normal, Contingencia)
	UpdateTransmissionType(string, string) error
	// ValidateStatus valida el estado de un DTE
	ValidateStatus(string, string) error
}
