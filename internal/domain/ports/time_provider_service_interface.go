package ports

import "time"

// TimeProvider es una interfaz que define los métodos que debe implementar un proveedor de tiempo
type TimeProvider interface {
	Now() time.Time        // Now retorna la fecha y hora actual
	Sleep(d time.Duration) // Sleep pausa la ejecución por un tiempo determinado
}
