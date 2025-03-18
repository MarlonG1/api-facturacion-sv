package main

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

func main() {
	err := env.InitEnvConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = logs.InitLogger()
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}
	logs.Info("Logger initialized successfully")
}
