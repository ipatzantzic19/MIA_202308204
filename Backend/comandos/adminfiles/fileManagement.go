package adminfiles

import (
	"fmt"
	"strings"
)

// FilesCommandProps es el manejador principal para los comandos relacionados con archivos,
// como 'mkfile' (crear archivo) y 'mkdir' (crear directorio).
// Recibe el comando específico y una lista de sus propiedades (parámetros).
func FilesCommandProps(command string, props []string) (string, error) {
	// Se convierte el comando a minúsculas para un manejo de casos uniforme y no sensible a mayúsculas.
	switch strings.ToLower(command) {
	case "mkfile":
		// Lógica para el comando 'mkfile'.
		// En esta fase, es una simulación que confirma la recepción del comando y sus propiedades.
		// TODO: Implementar la lógica real para crear un archivo dentro del sistema de archivos del disco virtual.
		fmt.Println(">> Recibido comando 'mkfile'")
		fmt.Println(">> Propiedades recibidas: ", props)
		return "Comando 'mkfile' procesado exitosamente (simulación)", nil

	case "mkdir":
		// Lógica para el comando 'mkdir'.
		// Similar a 'mkfile', actualmente es una simulación.
		// TODO: Implementar la lógica real para crear un directorio dentro del sistema de archivos del disco virtual.
		fmt.Println(">> Recibido comando 'mkdir'")
		fmt.Println(">> Propiedades recibidas: ", props)
		return "Comando 'mkdir' procesado exitosamente (simulación)", nil

	default:
		// Si el comando no es 'mkfile' ni 'mkdir', se retorna un error para indicar que no es un comando de archivos válido.
		return "", fmt.Errorf("comando de archivos '%s' no reconocido", command)
	}
}
