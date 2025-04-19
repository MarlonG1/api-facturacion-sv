package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/invalidation_models"
)

type InvalidationBasicStrategy struct {
	Document *invalidation_models.InvalidationDocument
}

func (s *InvalidationBasicStrategy) Validate() *dte_errors.DTEError {
	// 1. Validar campos requeridos de identificación
	if err := s.validateIdentification(); err != nil {
		return err
	}

	// 2. Validar campos requeridos del emisor
	if err := s.validateIssuer(); err != nil {
		return err
	}

	return nil
}

func (s *InvalidationBasicStrategy) validateIdentification() *dte_errors.DTEError {
	id := s.Document.Identification
	if id == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Identification")
	}

	//Validar campos requeridos
	if id.Version.GetValue() == 0 || id.Ambient.GetValue() == "" || id.GenerationCode.GetValue() == "" ||
		id.GetEmissionDate().IsZero() || id.GetEmissionTime().IsZero() {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Identification fields")
	}

	return nil
}

func (s *InvalidationBasicStrategy) validateIssuer() *dte_errors.DTEError {
	issuer := s.Document.Issuer
	if issuer == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Issuer")
	}

	// Validar campos requeridos del emisor según schema
	if issuer.NIT.GetValue() == "" || issuer.Name == "" || issuer.EstablishmentType.GetValue() == "" {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Issuer required fields")
	}

	// Validar correo y teléfono
	if issuer.Email.GetValue() == "" {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Issuer email")
	}

	return nil
}
