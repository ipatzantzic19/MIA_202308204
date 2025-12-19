package utils

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

func CrearBloqueCarpetaInicial(posActual int32, posPadre int32) int32 {
	posBloque := BuscarPosicionNewBloque()
	var newCarpeta structures.BloqueCarpeta
	newCarpeta.B_content[0].B_name = NameCarpeta12(".")
	newCarpeta.B_content[0].B_inodo = posActual
	newCarpeta.B_content[1].B_name = NameCarpeta12("..")
	newCarpeta.B_content[1].B_inodo = posPadre
	newCarpeta.B_content[2].B_name = NameCarpeta12("")
	newCarpeta.B_content[2].B_inodo = -1
	newCarpeta.B_content[3].B_name = NameCarpeta12("")
	newCarpeta.B_content[3].B_inodo = -1
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBloque), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Write(file, binary.LittleEndian, &newCarpeta); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return -1
	}
	return posBloque
}

func CrearBloqueCarpetaOtra(hijo int32, nombreH string) int32 {
	posBloque := BuscarPosicionNewBloque()
	var newCarpeta structures.BloqueCarpeta
	newCarpeta.B_content[0].B_name = NameCarpeta12(nombreH)
	newCarpeta.B_content[0].B_inodo = hijo
	for i := 1; i < 4; i++ {
		newCarpeta.B_content[i].B_name = NameCarpeta12("")
		newCarpeta.B_content[i].B_inodo = -1
	}
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBloque), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Write(file, binary.LittleEndian, &newCarpeta); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return -1
	}
	return posBloque
}

func CrearBloqueApuntador(hijo int32) int32 {
	posBloque := BuscarPosicionNewBloque()
	var newApuntador structures.BloqueApuntador
	newApuntador.B_pointers[0] = hijo
	for i := 1; i < 16; i++ {
		newApuntador.B_pointers[i] = -1
	}
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBloque), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Write(file, binary.LittleEndian, &newApuntador); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return -1
	}
	return posBloque
}

func CrearInodoCarpeta(pos int32, bloque int32) {
	var newInodo structures.TablaInodo
	newInodo.I_uid = global.UsuarioLogeado.UID
	newInodo.I_gid = global.UsuarioLogeado.GID
	newInodo.I_s = 0
	newInodo.I_atime = ObFechaInt()
	newInodo.I_ctime = ObFechaInt()
	newInodo.I_mtime = ObFechaInt()
	newInodo.I_type = 0
	newInodo.I_perm = 664
	newInodo.I_block[0] = bloque
	for i := 1; i < 15; i++ {
		newInodo.I_block[i] = -1
	}
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &newInodo); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return
	}
}

func CrearInodoArchivo(pos int32) {
	var newInodo structures.TablaInodo
	newInodo.I_uid = global.UsuarioLogeado.UID
	newInodo.I_gid = global.UsuarioLogeado.GID
	newInodo.I_s = 0
	newInodo.I_atime = ObFechaInt()
	newInodo.I_ctime = ObFechaInt()
	newInodo.I_mtime = ObFechaInt()
	newInodo.I_type = 1
	newInodo.I_perm = 664
	for j := 0; j < 15; j++ {
		newInodo.I_block[j] = -1
	}

	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &newInodo); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return
	}
}
