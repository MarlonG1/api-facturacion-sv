package interfaces

// Appendix es una interfaz que define los métodos que debe implementar un apéndice
type Appendix interface {
	GetField() string // GetField obtiene el campo del apéndice
	GetLabel() string // GetLabel obtiene la etiqueta del apéndice
	GetValue() string // GetValue obtiene el valor del apéndice
}
