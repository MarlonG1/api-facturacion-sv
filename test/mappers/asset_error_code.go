package mappers

import (
	"errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func assertErrorCode(t *testing.T, err error, expectedCode string) {
	expectedCode = strings.ToLower(expectedCode)
	var validationErr *dte_errors.ValidationError
	var serviceErr *shared_error.ServiceError
	var actualErr string

	if errors.As(err, &validationErr) {
		actualErr = strings.ToLower(validationErr.GetType())
	} else if errors.As(err, &serviceErr) {
		actualErr = strings.ToLower(serviceErr.GetCode())
	} else {
		t.Errorf("Error type not recognized, expected validation or service error with code %s", expectedCode)
	}

	assert.Equal(t, expectedCode, actualErr, "Error code does not match")
}
