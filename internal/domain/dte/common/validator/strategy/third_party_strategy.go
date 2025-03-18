package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ThirdPartyStrategy struct {
	Document interfaces.DTEDocument
}

// Validate valida que los items del DTE sean consistentes con la venta a tercer o que no haya venta a terceros
func (s *ThirdPartyStrategy) Validate() *dte_errors.DTEError {
	// Si no hay venta a terceros, es válido
	if s.Document.GetThirdPartySale() == nil {
		return nil
	}

	// Si hay venta a terceros, todos los items deben estar relacionados
	// No se pueden mezclar ventas propias con ventas a terceros
	for _, item := range s.Document.GetItems() {
		if err := s.validateThirdPartyItem(item); err != nil {
			return err
		}
	}

	// Validar que el DTE solo tenga un tercero
	if err := s.validateSingleThirdParty(); err != nil {
		return err
	}

	return nil
}

// validateThirdPartyItem valida cada item del DTE para asegurar que todos estén relacionados a la venta de tercero
func (s *ThirdPartyStrategy) validateThirdPartyItem(item interfaces.Item) *dte_errors.DTEError {
	// Validar que todos los items correspondan a la venta de tercero
	if s.Document.GetThirdPartySale() != nil {
		// No se permiten ventas mezcladas (propias y de terceros)
		if item.GetRelatedDoc() == nil {
			return dte_errors.NewDTEErrorSimple("MixedSalesNotAllowed")
		}
	}
	return nil
}

// validateSingleThirdParty valida que solo haya un tercero por DTE
func (s *ThirdPartyStrategy) validateSingleThirdParty() *dte_errors.DTEError {
	// Validar que solo haya un tercero por DTE
	// Esta validación es para asegurar que no se mezclen ventas de diferentes terceros
	if s.Document.GetThirdPartySale() == nil {
		return nil
	}

	// El nombre del tercero es requerido
	if s.Document.GetThirdPartySale().GetName() == "" {
		return dte_errors.NewDTEErrorSimple("RequiredField", "ThirdPartyName")
	}

	return nil
}
