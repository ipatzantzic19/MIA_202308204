package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/comandos/global"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

func BuscarPosicionNewInodo() int32 {
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	var bitI = 0
	var bit byte
	var one byte = '1'
	startI := Sb_System.S_bm_inode_start
	end := startI + Sb_System.S_inodes_count

	//
	bitI = 0
	for j := startI; j <= end; j++ {
		if _, err := file.Seek(int64(j), 0); err != nil {
			color.Red("[/]: Error en mover puntero")
			return -1
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[/]: Error en la lectura del archivo")
			return -1
		}
		if bit == '0' {
			if _, err := file.Seek(int64(j), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return -1
			}
			if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
				color.Red("[-]: Error en la escritura del archivo")
				return -1
			}
			break
		}
		bitI++
	}
	Sb_System.S_free_inodes_count -= 1
	posInodo := Sb_System.S_inode_start + (size.SizeTablaInodo() * int32(bitI))
	BuscarPrimerInodoVacio()
	return posInodo
}

func BuscarPrimerInodoVacio() {
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return
	}
	defer file.Close()
	var bitI = 0
	var bit byte
	startI := Sb_System.S_bm_inode_start
	end := startI + Sb_System.S_inodes_count
	bitI = 0
	for j := startI; j <= end; j++ {
		if _, err := file.Seek(int64(j), 0); err != nil {
			color.Red("[/]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[/]: Error en la lectura del archivo")
			return
		}
		if bit == '0' {
			bitI++
			break
		}
		bitI++
	}
	Sb_System.S_first_ino = int32(bitI)
}

func BuscarPosicionInodoBM(posInodo int32) int32 {
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	var bitI = 0
	var bit byte
	startI := Sb_System.S_bm_inode_start
	end := startI + Sb_System.S_inodes_count
	for j := startI; j <= end; j++ {
		if _, err := file.Seek(int64(j), 0); err != nil {
			color.Red("[/]: Error en mover puntero")
			return -1
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[/]: Error en la lectura del archivo")
			return -1
		}
		if (bit == '1') && (posInodo == (Sb_System.S_inode_start + (size.SizeTablaInodo() * int32(bitI)))) {
			return int32(bitI)
		}
		bitI++
	}
	return -1
}

func BuscarPosicionNewBloque() int32 {
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	var bit2 int32 = 0
	var bit byte
	var one byte = '1'
	start := Sb_System.S_bm_block_start
	end := start + Sb_System.S_blocks_count
	var posBloque int32 = -1
	for i := start; i <= end; i++ {
		if _, err := file.Seek(int64(i), 0); err != nil {
			color.Red("[/]: Error en mover puntero")
			return -1
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[/]: Error en la lectura del archivo")
			return -1
		}
		if bit == '0' {
			if _, err := file.Seek(int64(i), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return -1
			}
			if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
				color.Red("[-]: Error en la escritura del archivo")
				return -1
			}
			break
		}
		bit2++
	}
	Sb_System.S_free_blocks_count -= 1
	posBloque = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())
	BuscarPrimerBloqueVacio()
	return posBloque
}

func BuscarPrimerBloqueVacio() {
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return
	}
	defer file.Close()
	var bitI = 0
	var bit byte
	startI := Sb_System.S_bm_block_start
	end := startI + Sb_System.S_blocks_count
	bitI = 0
	for j := startI; j <= end; j++ {
		if _, err := file.Seek(int64(j), 0); err != nil {
			color.Red("[/]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[/]: Error en la lectura del archivo")
			return
		}
		if bit == '0' {
			bitI++
			break
		}
		bitI++
	}
	Sb_System.S_first_blo = int32(bitI)
}

func BuscarPosicionBloqueBM(posBlock int32) int32 {
	var bit2 int32 = 0
	var bit byte
	// var one byte = '1'
	start := Sb_System.S_bm_block_start
	end := start + Sb_System.S_blocks_count

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()

	for i := start; i < end; i++ {
		if _, err := file.Seek(int64(i), 0); err != nil {
			color.Red("[/]: Error en mover puntero")
			return -1
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[/]: Error en la lectura del SuperBloque")
			return -1
		}
		if (bit == '1') && (posBlock == (Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador()))) {
			return bit2
		}
		bit2++
	}

	return -1
}
