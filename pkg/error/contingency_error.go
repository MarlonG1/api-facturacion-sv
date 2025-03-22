package error

import "fmt"

type DocumentError struct {
	DocumentID string
	Err        error
}

func (e *DocumentError) Error() string {
	return fmt.Sprintf("error with document %s: %v", e.DocumentID, e.Err)
}

func (e *DocumentError) Unwrap() error {
	return e.Err
}

type BatchError struct {
	BatchID string
	Err     error
}

func (e *BatchError) Error() string {
	return fmt.Sprintf("error with batch %s: %v", e.BatchID, e.Err)
}

func (e *BatchError) Unwrap() error {
	return e.Err
}
