package error

import (
	"errors"
	"fmt"
)

var (
	ErrTokenNotFound           = errors.New("token not found in cache")
	ErrBranchOfficeNotFound    = errors.New("branch office not found or actually inactive")
	ErrUserNotFound            = errors.New("user not found or actually inactive")
	ErrBranchDoesNotBelong     = errors.New("branch office does not belong to the user")
	ErrExpiredToken            = fmt.Errorf("token has expired, please login again")
	ErrDTEDocumentNotFound     = errors.New("dte document not found")
	ErrHaciendaTokenGeneration = fmt.Errorf("failed to generate Hacienda token")
)
