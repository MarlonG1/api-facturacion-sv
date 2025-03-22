package structs

type InvalidationResponse struct {
	Identificacion InvalidationIdentification `json:"identificacion"`
	Emisor         InvalidationIssuer         `json:"emisor"`
	Documento      DocumentResponse           `json:"documento"`
	Motivo         ReasonResponse             `json:"motivo"`
}

type DocumentResponse struct {
	TipoDte           string  `json:"tipoDte"`
	CodigoGeneracion  string  `json:"codigoGeneracion"`
	SelloRecibido     string  `json:"selloRecibido"`
	NumeroControl     string  `json:"numeroControl"`
	FecEmi            string  `json:"fecEmi"`
	MontoIva          float64 `json:"montoIva"`
	CodigoGeneracionR *string `json:"codigoGeneracionR"`
	Nombre            *string `json:"nombre"`
	TipoDocumento     *string `json:"tipoDocumento"`
	NumDocumento      *string `json:"numDocumento"`
	Telefono          *string `json:"telefono"`
	Correo            *string `json:"correo"`
}

type ReasonResponse struct {
	TipoAnulacion     int     `json:"tipoAnulacion"`
	MotivoAnulacion   *string `json:"motivoAnulacion"`
	NombreResponsable string  `json:"nombreResponsable"`
	TipDocResponsable string  `json:"tipDocResponsable"`
	NumDocResponsable string  `json:"numDocResponsable"`
	NombreSolicita    string  `json:"nombreSolicita"`
	TipDocSolicita    string  `json:"tipDocSolicita"`
	NumDocSolicita    string  `json:"numDocSolicita"`
}

type InvalidationIdentification struct {
	Version          int    `json:"version"`
	Ambiente         string `json:"ambiente"`
	CodigoGeneracion string `json:"codigoGeneracion"`
	FecAnula         string `json:"fecAnula"`
	HorAnula         string `json:"horAnula"`
}

type InvalidationIssuer struct {
	NIT                   string  `json:"nit"`
	Nombre                string  `json:"nombre"`
	TipoEstablecimiento   string  `json:"tipoEstablecimiento"`
	Telefono              string  `json:"telefono"`
	Correo                string  `json:"correo"`
	CodigoEstablecimiento *string `json:"codEstable"`
	POSCodigo             *string `json:"codPuntoVenta"`
	NombreComercial       *string `json:"nomEstablecimiento"`
}
