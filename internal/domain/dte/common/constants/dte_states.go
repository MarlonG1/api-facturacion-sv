package constants

const (
	DocumentReceived = "RECEIVED"
	DocumentRejected = "REJECTED"
	DocumentInvalid  = "INVALIDATED"
	DocumentPending  = "PENDING"
)

const (
	TransmissionContingency = "CONTINGENCY"
	TransmissionNormal      = "NORMAL"
)

const (
	PhysicalDocument   = 1
	ElectronicDocument = 2
)

var (
	// ValidReceiverDocumentStates contiene los estados válidos para un documento tributario electrónico
	ValidReceiverDocumentStates = map[string]bool{
		DocumentReceived: true,
		DocumentPending:  true,
		DocumentRejected: true,
		DocumentInvalid:  true,
	}

	// ValidTransmissionTypes contiene los tipos de transmisión válidos para un documento tributario electrónico
	ValidTransmissionTypes = map[string]bool{
		TransmissionContingency: true,
		TransmissionNormal:      true,
	}
)
