package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseAppendix mapea los apéndices de un documento
func MapCommonResponseAppendix(request []interfaces.Appendix) []structs.DTEApendice {
	result := make([]structs.DTEApendice, len(request))
	for i, appendix := range request {
		result[i] = structs.DTEApendice{
			Campo:    appendix.GetField(),
			Etiqueta: appendix.GetLabel(),
			Valor:    appendix.GetValue(),
		}
	}

	return result
}
