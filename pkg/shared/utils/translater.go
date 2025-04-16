package utils

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"strings"
)

func TranslateMessage(key, code string, params ...interface{}) string {
	code = strings.ToLower(fmt.Sprintf("%s.%s", key, code))
	return i18n.Translate(code, params...)
}

func TranslateHealthUp(key string) string {
	return TranslateMessage("health.up", key)
}

func TranslateHealthDown(key string) string {
	return TranslateMessage("health.down", key)
}

func TranslateHealthError(key string, params ...interface{}) string {
	return TranslateMessage("health.error", key, params...)
}
