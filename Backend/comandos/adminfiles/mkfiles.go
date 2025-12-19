package adminfiles

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func Values_MKFILE(instructions []string) (string, bool, int32, string, bool, bool) {
	var path, cont string
	var r, overwrite bool
	var size int32 = 0

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("MKFILE", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "cont") {
			var value = utils.TieneContFile("MKFILE", valor)
			cont = value
		} else if strings.HasPrefix(strings.ToLower(valor), "r") {
			r = true
		} else if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = utils.TieneSizeV2("MKFILE", valor)
			size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "overwrite") {
			overwrite = true
		} else {
			color.Yellow("[MKFILE]: Atributo no reconocido")
			return "", false, -1, "", false, false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[MKFILE]: No hay path")
		return "", false, -1, "", false, false
	}
	if size < 0 {
		color.Red("[MKFILE]: Size no valido")
		return "", false, -1, "", false, false
	}
	return path, r, size, cont, true, overwrite
}

func MKFILE_EXECUTE(path string, r bool, _size int32, cont string, overwrite bool) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[MKDIR]: Usuario no logeado")
		return
	}
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[MKDIR]: Error al abrir archivo")
		return
	}
	defer file.Close()

	if cont == "" {
		if _size < 0 {
			color.Red("[MKFILE]: El valor size debe ser mayor o igual a 0")
		}
	} else {
		_size = 1
	}

	// var inicioSB int32
	if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[MKFILE]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		}
		inicioSB := nodo.Particion_P.Part_start
		if _, err := file.Seek(int64(inicioSB), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
		// sb := structures.SuperBloque{}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
			color.Red("[MKFILE]: Error en la lectura del SuperBloque")
			return
		}
	} else if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[MKFILE]: El disco (secundario) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		}
		inicioSB := nodo.Particion_L.Part_start + size.SizeEBR()
		if _, err := file.Seek(int64(inicioSB), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
		// sb := structures.SuperBloque{}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
			color.Red("[MKFILE]: Error en la lectura del SuperBloque")
			return
		}
	}

	rutaS := utils.SplitRuta(path)
	if len(rutaS) == 0 {
		color.Red("[MKFILE]: Ruta invalida")
		return
	}

	exist := utils.GetInodoFSystem(rutaS, 0, int32(len(rutaS)-1), utils.Sb_System.S_inode_start, nodo.Path)
	if exist != -1 {
		if overwrite {
			parentPath := ""
			if len(rutaS) > 1 {
				parentPath = strings.Join(rutaS[:len(rutaS)-1], "/")
			}
			parentInodoNum := utils.GetInodoFSystem(utils.SplitRuta(parentPath), 0, int32(len(rutaS)-2), utils.Sb_System.S_inode_start, nodo.Path)
			utils.DeleteFile(exist, nodo.Path, parentInodoNum, rutaS[len(rutaS)-1])
		} else {
			color.Red("[MKFILE]: Ya existe el archivo -> " + path + ". ¿Desea sobrescribirlo? (use -overwrite)")
			return
		}
	}

	posInodoI := utils.Sb_System.S_inode_start
	var existP = true
	var inodo structures.TablaInodo
	if len(rutaS) > 1 {
		for i := 0; i < (len(rutaS) - 1); i++ {
			if existP {
				aux := posInodoI
				posInodoI = utils.ExistPathSystem(rutaS, int32(i), posInodoI, nodo.Path)
				if posInodoI == aux {
					existP = false
				}
			}
			if !existP {
				if r {
					posInodoI = utils.CrearCarpetaSystem(rutaS, int32(i), posInodoI, nodo.Path)
					if nodo.Es_Particion_P {
						if _, err := file.Seek(int64(nodo.Particion_P.Part_start), 0); err != nil {
							color.Red("[MKFILE]: Error en mover puntero")
							return
						}
					} else if nodo.Es_Particion_L {
						if _, err := file.Seek(int64(nodo.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
							color.Red("[MKFILE]: Error en mover puntero")
							return
						}
					}
					if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
						color.Red("[MKFILE]: Error en la escritura del archivo")
						return
					}
					if posInodoI == -1 {
						return
					}
				} else {
					color.Red("[MKFILE]: No se puede crear el archivo")
					return
				}
			}
		}
	}

	if posInodoI == -1 {
		color.Red("[MKFILE]: Algo salio mal ")
		return
	}

	if !utils.ValidarPermisoWSystem(posInodoI, nodo.Path) {
		color.Red("[MKFILE]: No tiene permisos para escribir en el archivo -> " + rutaS[len(rutaS)-1])
		return
	}

	texto := ""
	if cont != "" {
		contenidoarchivo, eCa := os.ReadFile(cont)
		if eCa != nil {
			color.Red("[MKFILE]: Error al abrir el archivo")
			return
		}
		texto = utils.ToString(contenidoarchivo[:])
	} else {
		if _size > 0 {
			conta := 0
			texto = ""
			for i := 0; i < int(_size); i++ {
				texto += fmt.Sprint(conta)
				conta++
				if conta == 10 {
					conta = 0
				}
			}
		} else {
			texto = ""
		}
	}

	contenido := utils.SplitContent(texto)
	if len(contenido) >= 4380 {
		color.Red("[MKFILE]: El archivo es demasiado grande para la cantidad de bloques disponibles")
		return
	}

	if utils.Sb_System.S_free_blocks_count < int32(len(contenido)) {
		color.Red("[MKFILE]: No hay suficientes bloques para crear el archivo")
		return
	}

	if _, err := file.Seek(int64(posInodoI), 0); err != nil {
		color.Red("[MKFILE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[MKFILE]: Error en la lectura del archivo")
		return
	}

	var freeBlocks []int32
	var bit byte
	start := utils.Sb_System.S_bm_block_start
	end := start + utils.Sb_System.S_blocks_count
	bloquesNecesarios := int32(len(contenido))

	for i := start; i < end && len(freeBlocks) < int(bloquesNecesarios); i++ {
		if _, err := file.Seek(int64(i), 0); err != nil {
			color.Red("[MKFILE]: Error al mover el puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[MKFILE]: Error al leer el bitmap de bloques")
			return
		}
		if bit == '0' {
			freeBlocks = append(freeBlocks, i-start)
			// MARCAR BLOQUE COMO OCUPADO
			var uno byte = '1'
			if _, err := file.Seek(int64(i), 0); err != nil {
				color.Red("[MKFILE]: Error al mover el puntero para escribir en el bitmap")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &uno); err != nil {
				color.Red("[MKFILE]: Error al escribir en el bitmap de bloques")
				return
			}
		}
	}

	if int32(len(freeBlocks)) < bloquesNecesarios {
		color.Red("[MKFILE]: No hay suficientes bloques libres para crear el archivo")
		// DEBERÍA LIBERAR LOS BLOQUES QUE YA OCUPÉ
		return
	}

	utils.Sb_System.S_free_blocks_count -= int32(len(contenido))

	var bit2 int32 = 0
	for k := start; k < end; k++ {
		if _, err := file.Seek(int64(k), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[MKFILE]: Error en la lectura del archivo")
			return
		}
		if bit == '0' {
			break
		}
		bit2++
	}

	utils.Sb_System.S_first_blo = bit2
	if nodo.Es_Particion_P {
		if _, err := file.Seek(int64(nodo.Particion_P.Part_start), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
	} else if nodo.Es_Particion_L {
		if _, err := file.Seek(int64(nodo.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
	}

	if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[MKFILE]: Error en la escritura del archivo")
		return
	}

	var newInodoA structures.TablaInodo
	posNewI := utils.BuscarPosicionNewInodo()
	utils.CrearInodoArchivo(posNewI)
	if _, err := file.Seek(int64(posNewI), 0); err != nil {
		color.Red("[MKFILE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &newInodoA); err != nil {
		color.Red("[MKFILE]: Error en la lectura del archivo")
		return
	}

	//agregar carpeta
	atc := utils.AgregarCarpetaSystem(posNewI, posInodoI, rutaS[len(rutaS)-1])
	if atc == -1 {
		return
	}

	var tamanio int32 = 0
	for tm := range contenido {
		tamanio += int32(len(contenido[tm]))
	}

	newInodoA.I_s = tamanio
	newInodoA.I_atime = utils.ObFechaInt()
	newInodoA.I_mtime = utils.ObFechaInt()
	newInodoA.I_uid = global.UsuarioLogeado.UID
	newInodoA.I_gid = global.UsuarioLogeado.GID
	newInodoA.I_perm = 664

	for i, blockIndex := range freeBlocks {
		newInodoA = utils.AgregarArchivoSystem(contenido[i], newInodoA, int32(i), blockIndex)
	}

	if _, err := file.Seek(int64(posNewI), 0); err != nil {
		color.Red("[MKFILE]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &newInodoA); err != nil {
		color.Red("[MKFILE]: Error en la escritura del archivo")
		return
	}

	if nodo.Es_Particion_P {
		if _, err := file.Seek(int64(nodo.Particion_P.Part_start), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
			color.Red("[MKFILE]: Error en la escritura del archivo")
			return
		}
	} else if nodo.Es_Particion_L {
		if _, err := file.Seek(int64(nodo.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
			color.Red("[MKFILE]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
			color.Red("[MKFILE]: Error en la escritura del archivo")
			return
		}
	}

	color.Green("[MKFILE]: Se creo archivo -> " + path)
}
