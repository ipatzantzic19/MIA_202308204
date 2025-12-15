package partition

import (
	"Proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

// Values_FDISK obtiene y valida los parámetros específicos para el comando FDISK.
func Values_FDISK(instructions []string) (int32, string, [16]byte, byte, byte, byte, string, int32) {
	var _size int32
	var _diskName string
	var _name [16]byte
	var _unit byte = 'K'
	var _type byte = 'P'
	var _fit byte = 'W'
	var _delete string = "None"
	var _add int32 = 0
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
		} else if strings.HasPrefix(strings.ToLower(valor), "delete") {
			var value = utils.TieneDelete(valor)
			_delete = value
		} else if strings.HasPrefix(strings.ToLower(valor), "add") {
			var value = utils.TieneAdd(valor)
			_add = value
		} else {
			color.Yellow("[FDISK]: Atributo no reconocido")
		}
	}
	return _size, _diskName, _name, _unit, _type, _fit, _delete, _add
}

// FDISK_Create maneja la creación y modificación de particiones en un disco.
func FDISK_Create(_size int32, _diskName string, _name []byte, _unit byte, _type byte, _fit byte, _delete string, _add int32) {
	//fmt.Println(_name)
	path := "VDIC-MIA/Disks/" + _diskName
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] : El disco no existe en la ruta especificada.")
		return
	}

	// Delete
	if _delete != "None" {
		// Borrar particiones
		DeleteP(path, _name, _unit, _type, _fit)
		return
	}

	// Add
	if _add != 0 {
		if _add < 0 {
			RestE(path, _unit, _fit, _add, _name)
			return
		} else if _add > 0 {
			AddE(path, _unit, _fit, _add, _name)
			// fmt.Println("sumando")
			return
		}
	}

	if _type == 'P' {
		//primaria
		ParticionPrimaria(path, _size, _name, _unit, _type, _fit, _delete, _add)
	} else if _type == 'E' {
		//extended
		ParticionExtendida(path, _size, _name, _unit, _type, _fit, _delete, _add)
	} else if _type == 'L' {
		//logic
		ParticionLogica(path, _size, _name, _unit, _type, _fit, _delete, _add)
	} else {
		color.Red("[FDISK]: No reconocido Type")
		return
	}
}
