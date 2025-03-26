package constants

const (
	// Tipos de documentos electrónicos validos para emitir
	FacturaElectronica                = "01" // Factura Electrónica
	CCFElectronico                    = "03" // Comprobante de Crédito Fiscal Electrónico
	NotaRemisionElectronica           = "04" // Nota de Remisión Electrónica
	NotaCreditoElectronica            = "05" // Nota de Crédito Electrónica
	NotaDebitoElectronica             = "06" // Nota de Débito Electrónica
	ComprobanteRetencionElectronico   = "07" // Comprobante de Retención Electrónico
	ComprobanteLiquidacionElectronico = "08" // Comprobante de Liquidación Electrónico
	DocContableLiquidacionElectronico = "09" // Documento Contable de Liquidación Electrónico
	FacturaExportacionElectronica     = "11" // Factura de Exportación Electrónica
	FacturaSujetoExcluidoElectronica  = "14" // Factura Sujeto Excluido Electrónica
	ComprobanteDonacionElectronico    = "15" // Comprobante de Donación Electrónico

	// Tipos de documentos electrónicos validos para recibir
	NIT             = "36"
	DUI             = "13"
	CarnetResidente = "02"
	Pasaporte       = "03"
	OtroDocumento   = "37"
)

var (
	// ValidDTETypes Es una lista de valores permitidos para el campo DTEType
	ValidDTETypes = map[string]bool{
		FacturaElectronica:                true,
		CCFElectronico:                    true,
		NotaRemisionElectronica:           true,
		NotaCreditoElectronica:            true,
		NotaDebitoElectronica:             true,
		ComprobanteRetencionElectronico:   true,
		ComprobanteLiquidacionElectronico: true,
		DocContableLiquidacionElectronico: true,
		FacturaExportacionElectronica:     true,
		FacturaSujetoExcluidoElectronica:  true,
		ComprobanteDonacionElectronico:    true,
	}

	// ValidReceiverDTETypes Es una lista de valores permitidos para el campo DTEType de un receptor
	ValidReceiverDTETypes = []string{
		NIT,
		DUI,
		CarnetResidente,
		Pasaporte,
		OtroDocumento,
	}
)
