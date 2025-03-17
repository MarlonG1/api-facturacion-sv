package utils

import (
	"time"
)

var timezone *time.Location

func TimeInit() error {
	var err error
	timezone, err = time.LoadLocation("America/El_Salvador")
	if err != nil {
		return err
	}

	return nil
}

// TimeNow retorna el tiempo actual en la zona horaria configurada
func TimeNow() time.Time {
	return time.Now().In(timezone)
}
