package adminDisk

import (
	"Proyecto/comandos/utils"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Values_RMDISK analiza y valida los parámetros del comando RMDISK.
func Values_RMDISK(instructions []string) (string, bool) {
	var _diskName string = ""
	for _, valor := range instructions {
		// Identifica y asigna el valor del parámetro DISKNAME
		if strings.HasPrefix(strings.ToLower(valor), "diskname") {
			// Llama a la función auxiliar para validar y extraer el nombre del disco
			_diskName = utils.TieneDiskName("RMDISK", valor)
			break
		} else {
			color.Yellow("[RMDISK]: Atributo no reconocido")
			_diskName = ""
			break
		}
	}
	if _diskName == "" {
		return "", false
	} else {
		return _diskName, true
	}
}

func RMDISK_EXECUTE(diskName string) {
	PATH := "VDIC-MIA/Disks/" + diskName
	if _, err := os.Stat(PATH); os.IsNotExist(err) {
		color.Red("[RMDISK]: No existe el disco")
		return
	}
	err := os.Remove(PATH)
	if err != nil {
		color.Red("[RMDISK]: Error al borrar el disco")
		return
	}
	color.Green("[RMDISK]: Disco '" + diskName + "' Borrado")
}
