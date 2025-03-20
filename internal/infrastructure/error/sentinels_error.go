package error

import "errors"

var (
	ErrTokenNotFound = errors.New("token not found in cache")
)
