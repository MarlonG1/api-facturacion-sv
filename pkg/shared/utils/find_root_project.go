package utils

import (
	"os"
	"path/filepath"
)

// FindProjectRoot busca hacia arriba hasta encontrar el directorio raíz del proyecto y regresa su ruta
func FindProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Busca hacia arriba hasta encontrar el directorio raíz del proyecto
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Si llegamos a la raíz del sistema de archivos, detenemos la búsqueda
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Si no se encuentra, regresa el directorio de trabajo actual
	currentDir, _ := os.Getwd()
	return currentDir
}
