package location

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type Municipality struct {
	Value      string     `json:"value"`
	Department Department `json:"-"`
}

// NewMunicipality crea un nuevo objeto de tipo Municipality con el valor y el departamento especificados
func NewMunicipality(value string, department Department) (*Municipality, error) {
	mun := &Municipality{
		Value:      value,
		Department: department,
	}
	if mun.IsValid() {
		return mun, nil
	}
	return &Municipality{}, dte_errors.NewValidationError("InvalidMunicipality", value, department.GetValue(), getValidExamplesForDepartment(&department))
}

func NewValidatedMunicipality(value string, department string) *Municipality {
	return &Municipality{
		Value:      value,
		Department: *NewValidatedDepartment(department),
	}
}

// IsValid valida que el valor de Municipality sea un número válido para el departamento especificado
func (m *Municipality) IsValid() bool {
	pattern := m.getMunicipalityPattern()
	matched, _ := regexp.MatchString(pattern, m.Value)
	logs.Debug("Municipality VO", map[string]interface{}{
		"pattern":      pattern,
		"municipality": m.Value,
		"matched":      matched,
	})
	return matched
}

// getMunicipalityPattern retorna el patrón de expresión regular para validar el valor de Municipality según el departamento
func (m *Municipality) getMunicipalityPattern() string {
	switch m.Department.Value {
	case "01": // Ahuachapán
		return `^(13|14|15)$`
	case "02": // Santa Ana
		return `^(14|15|16|17)$`
	case "03": // Sonsonate
		return `^(17|18|19|20)$`
	case "04": // Chalatenango
		return `^(34|35|36)$`
	case "05": // La Libertad
		return `^(23|24|25|26|27|28)$`
	case "06": // San Salvador
		return `^(20|21|22|23|24)$`
	case "07": // Cuscatlán
		return `^(17|18)$`
	case "08": // La Paz
		return `^(23|24|25)$`
	case "09": // Cabañas
		return `^(10|11)$`
	case "10": // San Vicente
		return `^(14|15)$`
	case "11": // Usulután
		return `^(24|25|26)$`
	case "12": // San Miguel
		return `^(21|22|23)$`
	case "13": // Morazán
		return `^(27|28)$`
	case "14": // La Unión
		return `^(19|20)$`
	default:
		return `^(00)$` // Para extranjeros según el catálogo
	}
}

func (m *Municipality) Equals(other interfaces.ValueObject[string]) bool {
	return m.GetValue() == other.GetValue()
}

func (m *Municipality) GetValue() string {
	return m.Value
}

func (m *Municipality) ToString() string {
	return m.Value
}

// getValidExamplesForDepartment devuelve ejemplos de códigos de municipio válidos para el departamento
func getValidExamplesForDepartment(department *Department) string {
	switch department.GetValue() {
	case "01": // Ahuachapán
		return "13, 14, 15"
	case "02": // Santa Ana
		return "14, 15, 16, 17"
	case "03": // Sonsonate
		return "17, 18, 19, 20"
	case "04": // Chalatenango
		return "34, 35, 36"
	case "05": // La Libertad
		return "23, 24, 25, 26"
	case "06": // San Salvador
		return "20, 21, 22, 23, 24"
	case "07": // Cuscatlán
		return "17, 18"
	case "08": // La Paz
		return "23, 24, 25"
	case "09": // Cabañas
		return "10, 11"
	case "10": // San Vicente
		return "14, 15"
	case "11": // Usulután
		return "24, 25, 26"
	case "12": // San Miguel
		return "21, 22, 23"
	case "13": // Morazán
		return "27, 28"
	case "14": // La Unión
		return "19, 20"
	default:
		return "00"
	}
}
