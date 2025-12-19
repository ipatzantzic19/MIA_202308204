package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

func freeBlock(block_num int32, file *os.File) {
	if block_num == -1 {
		return
	}
	bm_block_pos := Sb_System.S_bm_block_start + block_num
	var cero byte = '0'
	if _, err := file.Seek(int64(bm_block_pos), 0); err != nil {
		color.Red("[DeleteFile]: Error al mover el puntero al bitmap de bloques")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[DeleteFile]: Error al escribir en el bitmap de bloques")
		return
	}
	Sb_System.S_free_blocks_count++
}

func freeIndirectBlocks(file *os.File, level int, block_pointer int32) {
	if level <= 0 {
		freeBlock((block_pointer-Sb_System.S_block_start)/size.SizeBloqueArchivo(), file)
		return
	}

	var apuntador structures.BloqueApuntador
	if _, err := file.Seek(int64(block_pointer), 0); err != nil {
		color.Red("[DeleteFile]: Error al mover el puntero al bloque de apuntadores")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[DeleteFile]: Error al leer el bloque de apuntadores")
		return
	}

	for i := 0; i < 16; i++ {
		if apuntador.B_pointers[i] != -1 {
			freeIndirectBlocks(file, level-1, apuntador.B_pointers[i])
		}
	}
	freeBlock((block_pointer-Sb_System.S_block_start)/size.SizeBloqueApuntador(), file)
}

func removeEntryFromDir(file *os.File, parent_inodo_num int32, fileName string) {
	var parentInodo structures.TablaInodo
	if _, err := file.Seek(int64(parent_inodo_num), 0); err != nil {
		color.Red("[DeleteFile]: Error al mover el puntero al inodo padre")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &parentInodo); err != nil {
		color.Red("[DeleteFile]: Error al leer el inodo padre")
		return
	}

	for i := 0; i < 12; i++ {
		if parentInodo.I_block[i] != -1 {
			var carpeta structures.BloqueCarpeta
			if _, err := file.Seek(int64(parentInodo.I_block[i]), 0); err != nil {
				color.Red("[DeleteFile]: Error al mover el puntero al bloque de carpeta")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
				color.Red("[DeleteFile]: Error al leer el bloque de carpeta")
				return
			}
			for j := 0; j < 4; j++ {
				if ToString(carpeta.B_content[j].B_name[:]) == fileName {
					carpeta.B_content[j].B_inodo = -1
					for k := 0; k < 12; k++ {
						carpeta.B_content[j].B_name[k] = 0
					}

					if _, err := file.Seek(int64(parentInodo.I_block[i]), 0); err != nil {
						color.Red("[DeleteFile]: Error al mover el puntero al bloque de carpeta para escribir")
						return
					}
					if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
						color.Red("[DeleteFile]: Error al escribir en el bloque de carpeta")
						return
					}
					return
				}
			}
		}
	}
}

func DeleteFile(inodo_num int32, path string, parent_inodo_num int32, fileName string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[DeleteFile]: Error al abrir el archivo")
		return
	}
	defer file.Close()

	// Leer el inodo
	var inodo structures.TablaInodo
	if _, err := file.Seek(int64(inodo_num), 0); err != nil {
		color.Red("[DeleteFile]: Error al mover el puntero al inodo")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[DeleteFile]: Error al leer el inodo")
		return
	}

	// Liberar bloques directos
	for i := 0; i < 12; i++ {
		if inodo.I_block[i] != -1 {
			freeBlock((inodo.I_block[i]-Sb_System.S_block_start)/size.SizeBloqueArchivo(), file)
		}
	}

	// Liberar bloques indirectos
	if inodo.I_block[12] != -1 {
		freeIndirectBlocks(file, 1, inodo.I_block[12])
	}
	if inodo.I_block[13] != -1 {
		freeIndirectBlocks(file, 2, inodo.I_block[13])
	}
	if inodo.I_block[14] != -1 {
		freeIndirectBlocks(file, 3, inodo.I_block[14])
	}

	// Liberar el inodo
	bm_inodo_pos := Sb_System.S_bm_inode_start + (inodo_num-Sb_System.S_inode_start)/size.SizeTablaInodo()
	var cero byte = '0'
	if _, err := file.Seek(int64(bm_inodo_pos), 0); err != nil {
		color.Red("[DeleteFile]: Error al mover el puntero al bitmap de inodos")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[DeleteFile]: Error al escribir en el bitmap de inodos")
		return
	}
	Sb_System.S_free_inodes_count++

	// Actualizar superbloque
	if _, err := file.Seek(int64(Sb_System.S_bm_inode_start-size.SizeSuperBloque()), 0); err != nil {
		color.Red("[DeleteFile]: Error al mover el puntero al superbloque")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &Sb_System); err != nil {
		color.Red("[DeleteFile]: Error al escribir el superbloque")
		return
	}

	removeEntryFromDir(file, parent_inodo_num, fileName)
}
