package interfaces

// AppendixGetter es una interfaz que define los métodos getter que debe implementar un apéndice
type AppendixGetter interface {
	GetField() string // GetField obtiene el campo del apéndice
	GetLabel() string // GetLabel obtiene la etiqueta del apéndice
	GetValue() string // GetValue obtiene el valor del apéndice
}

// AppendixSetter es una interfaz que define los métodos setter que debe implementar un apéndice
type AppendixSetter interface {
	SetField(field string) error // SetField establece el campo del apéndice
	SetLabel(label string) error // SetLabel establece la etiqueta del apéndice
	SetValue(value string) error // SetValue establece el valor del apéndice
}

// Appendix es una interfaz que combina los getters y setters de Appendix
type Appendix interface {
	AppendixGetter
	AppendixSetter
}
