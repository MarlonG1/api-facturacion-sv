package error

import "errors"

var (
	ErrTokenNotFound        = errors.New("token not found in cache")
	ErrBranchOfficeNotFound = errors.New("branch office not found or actually inactive")
	ErrUserNotFound         = errors.New("user not found or actually inactive")
	ErrBranchDoesNotBelong  = errors.New("branch office does not belong to the user")
	ErrDTEDocumentNotFound  = errors.New("dte document not found")
)
