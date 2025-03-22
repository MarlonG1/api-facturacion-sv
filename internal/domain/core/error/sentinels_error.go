package error

import "errors"

var (
	ErrBranchMatrixNotFound       = errors.New("branch matrix not found")
	ErrAtLeastOneBranch           = errors.New("at least one branch is required")
	ErrDontHaveBranchMatrix       = errors.New("don't have branch matrix, is required that the user have a branch matrix")
	ErrMoreThanOneBranchMatrix    = errors.New("more than one branch matrix, is required that the user have only one branch matrix")
	ErrBranchMatrixWithoutAddress = errors.New("for the branch matrix is required that have an address associated")
)
