package general

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// MkDisk crea un archivo que simula un disco con el tamaño solicitado.
// Acepta atributos como Size=3000, unit=K|M|G, name=<nombre> y path=<ruta>
func MkDisk(line string) (string, error) {
	attrs := ObtenerComandos(line)

	var size int64 = 0
	unit := "M"
	name := ""
	outPath := DiskPath

	for _, attr := range attrs {
		parts := strings.SplitN(attr, "=", 2)
		key := strings.ToLower(parts[0])
		val := ""
		if len(parts) > 1 {
			val = parts[1]
		}
		switch key {
		case "size":
			v, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				size = v
			}
		case "unit":
			if val != "" {
				unit = strings.ToUpper(val)
			}
		case "name":
			name = val
		case "path":
			if val != "" {
				outPath = val
			}
		}
	}

	if size <= 0 {
		return "", fmt.Errorf("Size inválido o no especificado")
	}

	// calcular bytes
	var mult int64 = 1024 * 1024 // default M
	switch unit {
	case "B":
		mult = 1
	case "K":
		mult = 1024
	case "M":
		mult = 1024 * 1024
	case "G":
		mult = 1024 * 1024 * 1024
	}

	bytes := size * mult

	// asegurar carpeta
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outPath, 0777); err != nil {
			return "", fmt.Errorf("no se pudo crear carpeta destino: %v", err)
		}
	}

	if name == "" {
		name = fmt.Sprintf("disk_%d.dsk", time.Now().Unix())
	}

	filePath := filepath.Join(outPath, name)

	// crear o truncar archivo al tamaño solicitado
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creando archivo: %v", err)
	}
	defer f.Close()

	if err := f.Truncate(bytes); err != nil {
		return "", fmt.Errorf("error asignando tamaño al disco: %v", err)
	}

	return fmt.Sprintf("Disco creado: %s (%d bytes)", filePath, bytes), nil
}

// RmDisk elimina un disco por nombre o ruta. Atributos: name=<nombre> o path=<ruta>
func RmDisk(line string) (string, error) {
	attrs := ObtenerComandos(line)

	name := ""
	targetPath := ""

	for _, attr := range attrs {
		parts := strings.SplitN(attr, "=", 2)
		key := strings.ToLower(parts[0])
		val := ""
		if len(parts) > 1 {
			val = parts[1]
		}
		switch key {
		case "name":
			name = val
		case "path":
			targetPath = val
		}
	}

	var filePath string
	if targetPath != "" {
		filePath = targetPath
	} else if name != "" {
		filePath = filepath.Join(DiskPath, name)
	} else {
		return "", fmt.Errorf("se requiere 'name' o 'path' para rmdisk")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo no encontrado: %s", filePath)
	}

	if err := os.Remove(filePath); err != nil {
		return "", fmt.Errorf("error eliminando disco: %v", err)
	}

	return fmt.Sprintf("Disco eliminado: %s", filePath), nil
}
