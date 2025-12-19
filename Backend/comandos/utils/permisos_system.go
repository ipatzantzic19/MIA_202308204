package utils

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func ValidarPermisoWSystem(posI int32, path string) bool {
	var inodo structures.TablaInodo
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return false
	}

	permiso := strconv.Itoa(int(inodo.I_perm))
	if len(permiso) == 1 {
		permiso = "00" + permiso
	} else if len(permiso) == 2 {
		permiso = "0" + permiso
	}

	if global.UsuarioLogeado.UID == 1 && global.UsuarioLogeado.GID == 1 {
		return true
	} else if (inodo.I_uid == global.UsuarioLogeado.UID) && (inodo.I_gid == global.UsuarioLogeado.GID) {
		if (permiso[0] == '2') || (permiso[0] == '3') || (permiso[0] == '6') || (permiso[0] == '7') {
			return true
		}
	} else if inodo.I_gid == global.UsuarioLogeado.GID {
		if (permiso[1] == '2') || (permiso[1] == '3') || (permiso[1] == '6') || (permiso[1] == '7') {
			return true
		}
	} else {
		if (permiso[2] == '2') || (permiso[2] == '3') || (permiso[2] == '6') || (permiso[2] == '7') {
			return true
		}
	}
	return false
}

func ValidarPermisoRSystem(posI int32, path string) bool {
	var inodo structures.TablaInodo
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return false
	}

	permiso := strconv.Itoa(int(inodo.I_perm))
	if len(permiso) == 1 {
		permiso = "00" + permiso
	} else if len(permiso) == 2 {
		permiso = "0" + permiso
	}

	if global.UsuarioLogeado.UID == 1 && global.UsuarioLogeado.GID == 1 {
		return true
	} else if (inodo.I_uid == global.UsuarioLogeado.UID) && (inodo.I_gid == global.UsuarioLogeado.GID) {
		if (permiso[0] == '4') || (permiso[0] == '5') || (permiso[0] == '6') || (permiso[0] == '7') {
			return true
		}
	} else if inodo.I_gid == global.UsuarioLogeado.GID {
		if (permiso[1] == '4') || (permiso[1] == '5') || (permiso[1] == '6') || (permiso[1] == '7') {
			return true
		}
	} else {
		if (permiso[2] == '4') || (permiso[2] == '5') || (permiso[2] == '6') || (permiso[2] == '7') {
			return true
		}
	}
	return false
}

func ValidarPermisoXSystem(posI int32, path string) bool {
	var inodo structures.TablaInodo
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return false
	}

	permiso := strconv.Itoa(int(inodo.I_perm))
	if len(permiso) == 1 {
		permiso = "00" + permiso
	} else if len(permiso) == 2 {
		permiso = "0" + permiso
	}

	if global.UsuarioLogeado.UID == 1 && global.UsuarioLogeado.GID == 1 {
		return true
	} else if (inodo.I_uid == global.UsuarioLogeado.UID) && (inodo.I_gid == global.UsuarioLogeado.GID) {
		if (permiso[0] == '1') || (permiso[0] == '3') || (permiso[0] == '5') || (permiso[0] == '7') {
			return true
		}
	} else if inodo.I_gid == global.UsuarioLogeado.GID {
		if (permiso[1] == '1') || (permiso[1] == '3') || (permiso[1] == '5') || (permiso[1] == '7') {
			return true
		}
	} else {
		if (permiso[2] == '1') || (permiso[2] == '3') || (permiso[2] == '5') || (permiso[2] == '7') {
			return true
		}
	}
	return false
}
