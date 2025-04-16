package structs

type RetentionDTEResponse struct {
	Identificacion  *DTEIdentification  `json:"identificacion"`
	Resumen         *RetentionSummary   `json:"resumen"`
	Emisor          RetentionIssuer     `json:"emisor"`
	Receptor        DTEReceiver         `json:"receptor"`
	CuerpoDocumento []RetentionItem     `json:"cuerpoDocumento"`
	Extension       *RetentionExtension `json:"extension"`
	Apendice        []DTEApendice       `json:"apendice"`
}

type RetentionSummary struct {
	TotalSujRetencion      float64 `json:"totalSujetoRetencion"`
	TotalIvaRetenido       float64 `json:"totalIVAretenido"`
	TotalIvaRetenidoLetras string  `json:"totalIVAretenidoLetras"`
}

type RetentionItem struct {
	NumItem            int     `json:"numItem"`
	TipoDTE            string  `json:"tipoDte"`
	TipoDoc            int     `json:"tipoDoc"`
	NumDoc             string  `json:"numDocumento"`
	FechaEmision       string  `json:"fechaEmision"`
	MontoSujetoGravado float64 `json:"montoSujetoGrav"`
	CodigoRetencionMH  string  `json:"codigoRetencionMH"`
	IvaRetenido        float64 `json:"ivaRetenido"`
	Descripcion        string  `json:"descripcion"`
}

type RetentionIssuer struct {
	NIT                 string     `json:"nit"`
	NRC                 string     `json:"nrc"`
	Nombre              string     `json:"nombre"`
	CodActividad        string     `json:"codActividad"`
	DescActividad       string     `json:"descActividad"`
	TipoEstablecimiento string     `json:"tipoEstablecimiento"`
	Direccion           DTEAddress `json:"direccion"`
	Telefono            string     `json:"telefono"`
	Correo              string     `json:"correo"`
	NombreComercial     *string    `json:"nombreComercial"`
	CodEstableMH        *string    `json:"codigoMH"`
	CodEstable          *string    `json:"codigo"`
	CodPuntoVentaMH     *string    `json:"puntoVentaMH"`
	CodPuntoVenta       *string    `json:"puntoVenta"`
}

type RetentionExtension struct {
	NombreEntrega    string  `json:"nombEntrega"`
	DocumentoEntrega string  `json:"docuEntrega"`
	NombreRecibe     string  `json:"nombRecibe"`
	DocumentoRecibe  string  `json:"docuRecibe"`
	Observacion      *string `json:"observaciones"`
}
