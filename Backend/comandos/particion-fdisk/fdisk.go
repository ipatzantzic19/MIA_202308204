package partition

import (
	"Proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

// Values_FDISK obtiene y valida los parámetros específicos para el comando FDISK.
func Values_FDISK(instructions []string) (int32, string, [16]byte, byte, byte, byte) {
	var _size int32
	var _diskName string
	var _name [16]byte
	var _unit byte = 'K'
	var _type byte = 'P'
	var _fit byte = 'W'

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = utils.TieneSize("FDISK", valor)
			_size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "diskname") {
			var value = utils.TieneDiskName("FDISK", valor)
			_diskName = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = utils.TieneNombre("FDISK", valor)
			if len(value) > 16 {
				color.Red("[FDISK]: El nombre no puede ser mayor a 16 caracteres")
				break
			} else {
				_name = utils.DevolverNombreByte(value)
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			var value = utils.TieneUnit("FDISK", valor)
			_unit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "type") {
			var value = utils.TieneTypeFDISK(valor)
			_type = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			var value = utils.TieneFit("FDISK", valor)
			_fit = value
		} else {
			color.Yellow("[FDISK]: Atributo no reconocido")
		}
	}
	return _size, _diskName, _name, _unit, _type, _fit
}

// FDISK_Create maneja la creación y modificación de particiones en un disco.
func FDISK_Create(_size int32, _diskName string, _name []byte, _unit byte, _type byte, _fit byte) {
	//fmt.Println(_name)
	path := "VDIC-MIA/Disks/" + _diskName
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] : El disco no existe en la ruta especificada.")
		return
	}

	if _type == 'P' {
		//primaria
		ParticionPrimaria(path, _size, _name, _unit, _type, _fit)
	} else if _type == 'E' {
		//extended
		ParticionExtendida(path, _size, _name, _unit, _type, _fit)
	} else if _type == 'L' {
		//logic
		ParticionLogica(path, _size, _name, _unit, _type, _fit)
	} else {
		color.Red("[FDISK]: No reconocido Type")
		return
	}
}
