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

	// ValidRetentionDTETypes Es una lista de valores permitidos para el campo DTEType de una retención
	ValidRetentionDTETypes = map[string]bool{
		FacturaElectronica:               true,
		CCFElectronico:                   true,
		FacturaSujetoExcluidoElectronica: true,
	}

	//ValidAdjustmentDTETypes  Es una lista de valores permitidos para el campo DTEType de un ajuste (Nota de crédito o débito)
	ValidAdjustmentDTETypes = map[string]bool{
		CCFElectronico:                  true,
		ComprobanteRetencionElectronico: true,
	}

	// ValidCCFDTETypesRelateDoc Es una lista de valores permitidos para el campo DTEType de un documento relacionado en un CCF
	ValidCCFDTETypesRelateDoc = map[string]bool{
		NotaRemisionElectronica:           true,
		ComprobanteLiquidacionElectronico: true,
		DocContableLiquidacionElectronico: true,
	}

	// ValidInvoiceDTETypesRelateDoc Es una lista de valores permitidos para el campo DTEType de un documento relacionado en una Factura
	ValidInvoiceDTETypesRelateDoc = map[string]bool{
		NotaRemisionElectronica:           true,
		DocContableLiquidacionElectronico: true,
	}

	// ValidDTETypesForContingency Es una lista de valores permitidos que se puede enviar por contingencia
	ValidDTETypesForContingency = map[string]bool{
		FacturaElectronica:               true,
		CCFElectronico:                   true,
		NotaRemisionElectronica:          true,
		NotaCreditoElectronica:           true,
		NotaDebitoElectronica:            true,
		FacturaExportacionElectronica:    true,
		FacturaSujetoExcluidoElectronica: true,
	}
)

func ShowValidRelatedDocTypes(valids map[string]bool) string {
	var result string
	for k := range valids {
		result += k + ", "
	}
	if len(result) > 2 {
		result = result[:len(result)-2]
	}
	return result
}
