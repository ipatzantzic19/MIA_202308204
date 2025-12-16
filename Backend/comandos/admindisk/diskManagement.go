package adminDisk

import (
	partition "Proyecto/comandos/particion-fdisk"
	"fmt"
	"strings"
)

// DiskCommandProps es el manejador principal para los comandos de disco.
// Identifica y ejecuta el comando de disco apropiado (como MKDISK, FDISK, etc.)
// basado en la entrada del usuario. Devuelve un mensaje de éxito y un error si ocurre alguno.
func DiskCommandProps(command string, instructions []string) (string, error) {

	// Este parámetro crea un disco
	if strings.ToUpper(command) == "MKDISK" {
		// 1. Llama a la función para obtener y validar los parámetros específicos de MKDISK.
		_size, _fit, _unit, err := Values_MKDISK(instructions)
		if err != nil {
			return "", err
		} // Si la validación de parámetros falla, se devuelve el error

		// 2. Llama a la función de creación del disco con los parámetros validados.
		msg, err := MKDISK_Create(_size, _fit, _unit)

		if err != nil {
			return "", err
		}

		// 3. Si todo es exitoso, devuelve el mensaje de éxito.
		return msg, nil
	}
	// Este parámetro elimina un disco
	if strings.ToUpper(command) == "RMDISK" {
		_diskName, valid := Values_RMDISK(instructions)
		if !valid {
			return "", fmt.Errorf("[DiskCommandProps]: Error al validar los parámetros de RMDISK")
		}
		RMDISK_EXECUTE(_diskName) // Ejecuta la eliminación del disco
		return "[RMDISK]: Disco eliminado exitosamente", nil
	}

	if strings.ToUpper(command) == "FDISK" {
		// 1. Llama a la función para obtener los parámetros de FDISK.
		_size, _diskName, _name, _unit, _type, _fit := partition.Values_FDISK(instructions)

		// 2. Llama a la función de creación/modificación de la partición.
		partition.FDISK_Create(_size, _diskName, _name[:], _unit, _type, _fit)

		// 3. Devuelve un mensaje de éxito genérico.
		return "[FDISK]: Comando ejecutado, revise la consola para más detalles.", nil
	}
	if strings.ToUpper(command) == "MOUNT" {
		// 1. Llama a la función para obtener los parámetros de MOUNT.
		_diskName, _name, _error := Values_Mount(instructions)
		if _error {
			return "", fmt.Errorf("[DiskCommandProps]: Error al validar los parámetros de MOUNT")
		}

		// 2. Llama a la función de montaje de la partición.
		MOUNT_EXECUTE(_diskName, _name[:])

		// 3. Devuelve un mensaje de éxito genérico.
		return "[MOUNT]: Comando ejecutado, revise la consola para más detalles.", nil
	}
	if strings.ToUpper(command) == "MOUNTED" {
		MOUNTED_EXECUTE()
		return "[MOUNTED]: Comando ejecutado, revise la consola para más detalles.", nil
	}

	// Si el comando no es MKDISK (o cualquier otro comando de disco implementado),
	// se devuelve un error indicando que el comando no está implementado.
	return "", fmt.Errorf("[DiskCommandProps]: Comando de disco '%s' no implementado o error interno", command)
}
