package adminDisk

import (
	"fmt"
	"strings"
)

// DiskCommandProps es el manejador principal para los comandos de disco.
// Identifica y ejecuta el comando de disco apropiado (como MKDISK, FDISK, etc.)
// basado en la entrada del usuario. Devuelve un mensaje de éxito y un error si ocurre alguno.
func DiskCommandProps(command string, instructions []string) (string, error) {

	// Manejo del comando MKDISK
	if strings.ToUpper(command) == "MKDISK" {
		// 1. Llama a la función para obtener y validar los parámetros específicos de MKDISK.
		_size, _fit, _unit, err := Values_MKDISK(instructions)
		if err != nil {
			// Si la validación de parámetros falla, se devuelve el error para ser manejado en un nivel superior.
			return "", err
		}

		// 2. Llama a la función de creación del disco con los parámetros validados.
		msg, err := MKDISK_Create(_size, _fit, _unit)
		if err != nil {
			// Si la creación del disco falla, se devuelve el error.
			return "", err
		}

		// 3. Si todo es exitoso, devuelve el mensaje de éxito.
		return msg, nil
	}

	// Aquí se pueden añadir más casos para otros comandos de disco como FDISK, RMDISK, etc.
	// Por ejemplo:
	// if strings.ToUpper(command) == "FDISK" {
	//     // Lógica para FDISK
	// }

	// Si el comando no es MKDISK (o cualquier otro comando de disco implementado),
	// se devuelve un error indicando que el comando no está implementado.
	return "", fmt.Errorf("[DiskCommandProps]: Comando de disco '%s' no implementado o error interno", command)
}
