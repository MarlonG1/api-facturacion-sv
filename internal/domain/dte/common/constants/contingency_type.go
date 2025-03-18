package constants

const (
	FallaConexionSistema  = iota + 1 // Falla en conexiones del sistema del emisor
	FallaServicioInternet            // Falla en el suministro del servicio de Internet
	FallaEnergiaElectrica            // Falla en el suministro del servicio de energía eléctrica
	NoDisponibilidadMH               // No disponibilidad de sistema del MH
	OtroMotivo                       // Otro motivo
)

var (
	// AllowedContingencyTypes contiene los tipos de contingencias permitidos, usado para validaciones
	AllowedContingencyTypes = []int{
		FallaConexionSistema,
		FallaServicioInternet,
		FallaEnergiaElectrica,
		NoDisponibilidadMH,
		OtroMotivo,
	}

	// ContingencyReasons contiene los motivos de contingencia permitidos, usado para validaciones
	ContingencyReasons = map[int]string{
		FallaConexionSistema:  "Error de conexión con sistemas internos",
		FallaServicioInternet: "Falla en el servicio de internet",
		FallaEnergiaElectrica: "Interrupción del servicio eléctrico",
		NoDisponibilidadMH:    "Servicio del Ministerio de Hacienda no disponible",
		OtroMotivo:            "Error en el proceso de emisión del documento",
	}
)

func GetContingencyReason(contingencyType int) string {
	if reason, exists := ContingencyReasons[contingencyType]; exists {
		return reason
	}
	return ContingencyReasons[5]
}
