package ports

import "encoding/json"

type DTERepository interface {
	GetDTEByID(id string) (json.RawMessage, error)
}
