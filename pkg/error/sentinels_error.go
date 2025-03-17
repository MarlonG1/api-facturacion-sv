package error

import "fmt"

var (
	ErrFailedToCreateLogFiles = fmt.Errorf("failed to create log files")
	ErrLogDirectoryNotFound   = fmt.Errorf("log directory not found")
)
