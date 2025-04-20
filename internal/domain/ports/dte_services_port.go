package ports

import "context"

// DTEService es una interfaz para los servicios de DTE
type DTEService interface {
	Create(ctx context.Context, data interface{}, branchID uint) (interface{}, error)
}
