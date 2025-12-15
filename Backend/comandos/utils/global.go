package utils

import (
	"bytes"
	"os"

	"github.com/fatih/color"
)

// ToByte convierte una cadena de texto en un slice de bytes.
func ToByte(str string) []byte {
	result := make([]byte, 1)
	copy(result[:], str)
	return result
}

// ToString convierte un slice de bytes en una cadena de texto,
func ToString(b []byte) string {
	nullIndex := bytes.IndexByte(b, 0)
	if nullIndex == -1 {
		return string(b)
	}
	return string(b[:nullIndex])
}

// ExisteArchivo verifica si un archivo existe en la ruta especificada.
func ExisteArchivo(comando string, archivo string) bool {
	if _, err := os.Stat(archivo); os.IsNotExist(err) {
		color.Red("[" + comando + "]: Archivo No Encontrado")
		return false
	}
	return true
}
