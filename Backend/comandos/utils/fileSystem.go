package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var Sb_System structures.SuperBloque
var CambioContSystem = false

func TieneContFile(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "cont=") {
		color.Red("[" + comando + "]: No tiene cont o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene cont Valido")
		return ""
	}
	return value[1]
}

func TieneDestinoFile(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "destino=") {
		color.Red("[" + comando + "]: No tiene destino o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene destino Valido")
		return ""
	}
	return value[1]
}

func TieneCat(comando string, valor string) {
	re := regexp.MustCompile(`file\d+=`)
	if !re.MatchString(valor) {
		color.Red("[" + comando + "]: No tiene fileN o tiene un valor no valido")
		return
	}
}
func TieneSizeV2(comando string, valor string) int32 {
	if !strings.HasPrefix(strings.ToLower(valor), "size=") {
		color.Red("[" + comando + "]: No tiene size o tiene un valor no valido")
		return -1
	}
	value := strings.Split(valor, "=")
	i, err := strconv.Atoi(value[1])
	if err != nil {
		fmt.Println("Error conversion", err)
		return -1
	}
	if i < 0 {
		color.Red("[" + comando + "]: No tiene size Valido")
		return -1
	}
	return int32(i)
}

func TieneNameRename(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "name=") {
		color.Red("[" + comando + "]: No tiene name o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene un name valido")
		return ""
	}

	return value[1]
}

func TieneFile(valor string) string {
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[CAT]: No tiene file valido")
		return ""
	}
	return value[1]
}

func EscribirJournalSystem(tipoOp string, tipo byte, nombre string, contenido string, nodo global.ParticionesMontadas) {
	var actJour, newJour structures.Journal
	file, err := os.OpenFile(global.UsuarioLogeado.Mounted.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[-]: Error al abrir archivo")
		return
	}
	defer file.Close()
	// mount := global.UsuarioLogeado.Mounted
	if nodo.Es_Particion_P {
		if _, err := file.Seek(int64(nodo.Particion_P.Part_start+size.SizeSuperBloque()), 0); err != nil {
			color.Red("[-]: Error en mover puntero")
			return
		}
	} else if nodo.Es_Particion_L {
		if _, err := file.Seek(int64(nodo.Particion_L.Part_start+size.SizeEBR()+size.SizeSuperBloque()), 0); err != nil {
			color.Red("[-]: Error en mover puntero")
			return
		}
	}
	if err := binary.Read(file, binary.LittleEndian, &actJour); err != nil {
		color.Red("[-]: Error en la lectura del archivo")
		return
	}
	for actJour.J_Sig != -1 {
		if _, err := file.Seek(int64(actJour.J_Start+size.SizeJournal()), 0); err != nil {
			color.Red("[-]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &actJour); err != nil {
			color.Red("[-]: Error en la lectura del archivo")
			return
		}
	}

	if (actJour.J_Start + size.SizeJournal()) > Sb_System.S_bm_inode_start {
		return
	}

	actJour.J_Sig = actJour.J_Start + size.SizeJournal()
	newJour.J_Start = actJour.J_Sig
	newJour.J_Tipo_Operacion = NameArchivosByte(tipoOp)
	newJour.J_Tipo = tipo
	newJour.J_Path = ObJournalData(nombre)
	newJour.J_Contenido = ObJournalData(contenido)
	newJour.J_Fecha = ObFechaInt()
	newJour.J_Sig = -1

	if _, err := file.Seek(int64(actJour.J_Start), 0); err != nil {
		color.Red("[-]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &actJour); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return
	}

	if _, err := file.Seek(int64(newJour.J_Start), 0); err != nil {
		color.Red("[-]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &newJour); err != nil {
		color.Red("[-]: Error en la escritura del archivo")
		return
	}

}

// ---------------------------------------------------------------------------------------
func AgregarArchivoSystem(cadena string, inodo structures.TablaInodo, j int32, aux int32) structures.TablaInodo {
	var pointer1, pointer2, pointer3, newPointer1, newPointer2, newPointer3 structures.BloqueApuntador
	var in structures.TablaInodo
	in.I_type = -1
	i := 0

	file, err := os.OpenFile(global.UsuarioLogeado.Mounted.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return structures.TablaInodo{}
	}
	defer file.Close()

	for i = 0; i < 15; i++ {
		if (inodo.I_block[i] != -1) && (i < 12) && (i == int(j)) {
			var archivo structures.BloqueArchivo
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return structures.TablaInodo{}
			}
			if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return structures.TablaInodo{}
			}
			archivo.B_content = DevolverContenidoArchivo(cadena)
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return structures.TablaInodo{}
			}
			if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
				color.Red("[/]: Error en la escritura del archivo")
				return structures.TablaInodo{}
			}
			return inodo
		} else if (inodo.I_block[i] == -1) && (i < 12) && (aux != -1) {
			var archivo structures.BloqueArchivo
			seek := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
			archivo.B_content = DevolverContenidoArchivo(cadena)
			if _, err := file.Seek(int64(seek), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return structures.TablaInodo{}
			}
			if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return structures.TablaInodo{}
			}
			inodo.I_block[i] = seek
			CambioCont = true
			return inodo
		} else if (i == 12) && (inodo.I_block[i] == -1) && (aux != -1) {
			if Sb_System.S_free_blocks_count > 0 {
				var bit2 int32 = 0
				var bit byte
				var one byte = '1'
				start := Sb_System.S_bm_block_start
				end := start + Sb_System.S_blocks_count
				for i := start; i < end; i++ {
					if _, err := file.Seek(int64(i), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						if _, err := file.Seek(int64(i), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						break
					}
					bit2++
				}

				inodo.I_block[i] = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())

				var archivo structures.BloqueArchivo
				posbloque := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
				archivo.B_content = DevolverContenidoArchivo(cadena)
				if _, err := file.Seek(int64(posbloque), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}
				newPointer1.B_pointers[0] = posbloque
				CambioCont = true
				for j := 1; j < 16; j++ {
					newPointer1.B_pointers[j] = -1
				}
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &newPointer1); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}

				bit2 = 0
				for bmi := start; bmi < end; bmi++ {
					if _, err := file.Seek(int64(bmi), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						break
					}
					bit2++
				}
				Sb_System.S_free_blocks_count -= 1
				Sb_System.S_first_blo = bit2
				return inodo
			} else {
				color.Red("[/]: No hay espacio disponible")
				return in
			}
		} else if (i == 12) && (inodo.I_block[i] != -1) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return structures.TablaInodo{}
			}
			if err := binary.Read(file, binary.LittleEndian, &pointer1); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return structures.TablaInodo{}
			}
			for p1 := 0; p1 < 16; p1++ {
				if (pointer1.B_pointers[p1] == -1) && (aux != -1) {
					var archivo structures.BloqueArchivo
					var posBloque = Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
					archivo.B_content = DevolverContenidoArchivo(cadena)
					if _, err := file.Seek(int64(posBloque), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
						color.Red("[/]: Error en la escritura del archivo")
						return structures.TablaInodo{}
					}
					CambioCont = true
					pointer1.B_pointers[p1] = posBloque
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
						color.Red("[/]: Error en la escritura del archivo")
						return structures.TablaInodo{}
					}
					return inodo
				} else if (pointer1.B_pointers[p1] != -1) && (j == int32(12+p1)) {
					var archivo structures.BloqueArchivo
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					archivo.B_content = DevolverContenidoArchivo(cadena)
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
						color.Red("[/]: Error en la escritura del archivo")
						return structures.TablaInodo{}
					}
					return inodo
				}
			}
		} else if (i == 13) && (inodo.I_block[i] == -1) && (aux != -1) {
			if Sb_System.S_free_blocks_count > 1 {
				//apuntador 1
				var bit2 int32 = 0
				var bit byte
				var one byte = '1'
				start := Sb_System.S_bm_block_start
				end := start + Sb_System.S_blocks_count
				for i := start; i < end; i++ {
					if _, err := file.Seek(int64(i), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						if _, err := file.Seek(int64(i), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						break
					}
					bit2++
				}

				inodo.I_block[i] = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())

				//apuntador 2
				bit2 = 0
				for i := start; i < end; i++ {
					if _, err := file.Seek(int64(i), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						if _, err := file.Seek(int64(i), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						break
					}
					bit2++
				}

				newPointer1.B_pointers[0] = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())

				var archivo structures.BloqueArchivo
				posbloque := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
				archivo.B_content = DevolverContenidoArchivo(cadena)
				if _, err := file.Seek(int64(posbloque), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}
				CambioCont = true
				newPointer2.B_pointers[0] = posbloque

				for j := 1; j < 16; j++ {
					newPointer1.B_pointers[j] = -1
					newPointer2.B_pointers[j] = -1
				}
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &newPointer1); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}
				if _, err := file.Seek(int64(newPointer1.B_pointers[0]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &newPointer2); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}

				bit2 = 0
				for bmi := start; bmi < end; bmi++ {
					if _, err := file.Seek(int64(bmi), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						break
					}
					bit2++
				}
				Sb_System.S_free_blocks_count -= 2
				Sb_System.S_first_blo = bit2
				return inodo
			} else {
				color.Red("[/]: No hay espacio disponible")
				return in
			}
		} else if (i == 13) && (inodo.I_block[i] != -1) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return structures.TablaInodo{}
			}
			if err := binary.Read(file, binary.LittleEndian, &pointer1); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return structures.TablaInodo{}
			}
			for p1 := 0; p1 < 16; p1++ {
				if pointer1.B_pointers[p1] == -1 {
					//****************/----------------
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &pointer2); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					for p2 := 0; p2 < 16; p2++ {
						if (pointer2.B_pointers[p2] == -1) && (aux != -1) {
							var archivo structures.BloqueArchivo
							var posBloque = Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
							archivo.B_content = DevolverContenidoArchivo(cadena)
							if _, err := file.Seek(int64(posBloque), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
								color.Red("[/]: Error en la escritura del archivo")
								return structures.TablaInodo{}
							}
							CambioCont = true
							pointer2.B_pointers[p2] = posBloque
							if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Write(file, binary.LittleEndian, &pointer2); err != nil {
								color.Red("[/]: Error en la escritura del archivo")
								return structures.TablaInodo{}
							}
							return inodo
						} else if (pointer2.B_pointers[p2] != -1) && (j == int32(28+p2+(16*p1))) {
							var archivo structures.BloqueArchivo
							if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return structures.TablaInodo{}
							}
							archivo.B_content = DevolverContenidoArchivo(cadena)
							if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
								color.Red("[/]: Error en la escritura del archivo")
								return structures.TablaInodo{}
							}
							return inodo
						}
					}
					//****************/----------------
				} else if (pointer1.B_pointers[p1] == -1) && (aux != -1) {
					if Sb_System.S_free_blocks_count > 0 {
						pointer1.B_pointers[p1] = Sb_System.S_block_start + (Sb_System.S_first_blo * size.SizeBloqueApuntador())

						var archivo structures.BloqueArchivo
						posBloque := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
						archivo.B_content = DevolverContenidoArchivo(cadena)
						if _, err := file.Seek(int64(posBloque), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						CambioCont = true
						newPointer2.B_pointers[0] = posBloque
						for inicializer := 1; inicializer < 16; inicializer++ {
							newPointer2.B_pointers[inicializer] = -1
						}

						if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &newPointer2); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}

						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}

						start := Sb_System.S_bm_block_start
						end := start + Sb_System.S_blocks_count
						var bit2 int32 = 0
						var bit byte
						var bandera = false
						var one byte = '1'
						//actualización sb
						for bmi := start; bmi < end; bmi++ {
							if _, err := file.Seek(int64(bmi), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return structures.TablaInodo{}
							}
							if bit == '0' && bandera {
								break
							}
							if bit == '0' && !bandera {
								if _, err := file.Seek(int64(bmi), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return structures.TablaInodo{}
								}
								if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return structures.TablaInodo{}
								}
								bandera = true
							}
							bit2++
						}
						Sb_System.S_free_blocks_count -= 1
						Sb_System.S_first_blo = bit2

						return inodo
					} else {
						color.Red("[/]: No hay espacio disponible")
						return in
					}
				}
			}
		} else if (i == 14) && (inodo.I_block[i] == -1) && (aux != -1) {
			if Sb_System.S_free_blocks_count > 2 {
				// apuntador 1
				var bit2 int32 = 0
				var bit byte
				var one byte = '1'
				start := Sb_System.S_bm_block_start
				end := start + Sb_System.S_blocks_count
				for i := start; i < end; i++ {
					if _, err := file.Seek(int64(i), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						if _, err := file.Seek(int64(i), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						break
					}
					bit2++
				}

				inodo.I_block[i] = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())

				// apuntador 2
				bit2 = 0
				for i := start; i < end; i++ {
					if _, err := file.Seek(int64(i), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						if _, err := file.Seek(int64(i), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						break
					}
					bit2++
				}

				newPointer1.B_pointers[0] = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())

				//Tercer apuntador
				bit2 = 0
				for i := start; i < end; i++ {
					if _, err := file.Seek(int64(i), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						if _, err := file.Seek(int64(i), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						break
					}
					bit2++
				}

				newPointer2.B_pointers[0] = Sb_System.S_block_start + (bit2 * size.SizeBloqueApuntador())

				var archivo structures.BloqueArchivo
				posbloque := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
				archivo.B_content = DevolverContenidoArchivo(cadena)
				if _, err := file.Seek(int64(posbloque), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}
				CambioCont = true
				newPointer3.B_pointers[0] = posbloque
				for j := 1; j < 16; j++ {
					newPointer1.B_pointers[j] = -1
					newPointer2.B_pointers[j] = -1
					newPointer3.B_pointers[j] = -1
				}
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &newPointer1); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}
				if _, err := file.Seek(int64(newPointer1.B_pointers[0]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &newPointer2); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}
				if _, err := file.Seek(int64(newPointer2.B_pointers[0]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return structures.TablaInodo{}
				}
				if err := binary.Write(file, binary.LittleEndian, &newPointer3); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return structures.TablaInodo{}
				}

				bit2 = 0
				for bmi := start; bmi < end; bmi++ {
					if _, err := file.Seek(int64(bmi), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					if bit == '0' {
						break
					}
					bit2++
				}
				Sb_System.S_free_blocks_count -= 3
				Sb_System.S_first_blo = bit2

				return inodo
			} else {
				color.Red("[/]: No hay espacio disponible")
				return in
			}
		} else if (i == 14) && (inodo.I_block[i] != -1) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return structures.TablaInodo{}
			}
			if err := binary.Read(file, binary.LittleEndian, &pointer1); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return structures.TablaInodo{}
			}
			for p1 := 0; p1 < 16; p1++ {
				if pointer1.B_pointers[p1] == -1 {
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return structures.TablaInodo{}
					}
					if err := binary.Read(file, binary.LittleEndian, &pointer2); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return structures.TablaInodo{}
					}
					for p2 := 0; p2 < 16; p2++ {
						if pointer2.B_pointers[p2] == -1 {
							if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Read(file, binary.LittleEndian, &pointer3); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return structures.TablaInodo{}
							}
							for p3 := 0; p3 < 16; p3++ {
								if (pointer3.B_pointers[p3] == -1) && (aux != -1) {
									var archivo structures.BloqueArchivo
									var posBloque = Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
									archivo.B_content = DevolverContenidoArchivo(cadena)
									if _, err := file.Seek(int64(posBloque), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return structures.TablaInodo{}
									}
									if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
										color.Red("[/]: Error en la escritura del archivo")
										return structures.TablaInodo{}
									}
									CambioCont = true
									pointer3.B_pointers[p3] = posBloque
									if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return structures.TablaInodo{}
									}
									if err := binary.Write(file, binary.LittleEndian, &pointer3); err != nil {
										color.Red("[/]: Error en la escritura del archivo")
										return structures.TablaInodo{}
									}
									return inodo
								} else if (pointer3.B_pointers[p3] != -1) && (j == int32(284+p3+(16*p2)+(256*p1))) {
									var archivo structures.BloqueArchivo
									if _, err := file.Seek(int64(pointer3.B_pointers[p3]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return structures.TablaInodo{}
									}
									if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
										color.Red("[/]: Error en la lectura del archivo")
										return structures.TablaInodo{}
									}
									archivo.B_content = DevolverContenidoArchivo(cadena)
									if _, err := file.Seek(int64(pointer3.B_pointers[p3]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return structures.TablaInodo{}
									}
									if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
										color.Red("[/]: Error en la escritura del archivo")
										return structures.TablaInodo{}
									}
									return inodo
								}
							}
						} else if (pointer2.B_pointers[p2] == -1) && (aux != -1) {
							if Sb_System.S_free_blocks_count > 0 {
								pointer2.B_pointers[p2] = Sb_System.S_block_start + (Sb_System.S_first_blo * size.SizeBloqueApuntador())

								var archivo structures.BloqueArchivo
								posBloque := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
								archivo.B_content = DevolverContenidoArchivo(cadena)
								if _, err := file.Seek(int64(posBloque), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return structures.TablaInodo{}
								}
								if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return structures.TablaInodo{}
								}
								CambioCont = true
								newPointer3.B_pointers[0] = posBloque
								for inicializer := 1; inicializer < 16; inicializer++ {
									newPointer3.B_pointers[inicializer] = -1
								}

								if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return structures.TablaInodo{}
								}
								if err := binary.Write(file, binary.LittleEndian, &newPointer3); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return structures.TablaInodo{}
								}

								if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return structures.TablaInodo{}
								}
								if err := binary.Write(file, binary.LittleEndian, &pointer2); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return structures.TablaInodo{}
								}

								start := Sb_System.S_bm_block_start
								end := start + Sb_System.S_blocks_count
								var bit2 int32 = 0
								var bit byte
								var bandera = false
								var one byte = '1'
								//actualización sb
								for bmi := start; bmi < end; bmi++ {
									if _, err := file.Seek(int64(bmi), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return structures.TablaInodo{}
									}
									if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
										color.Red("[/]: Error en la lectura del archivo")
										return structures.TablaInodo{}
									}
									if bit == '0' && bandera {
										break
									}
									if bit == '0' && !bandera {
										if _, err := file.Seek(int64(bmi), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return structures.TablaInodo{}
										}
										if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
											color.Red("[/]: Error en la escritura del archivo")
											return structures.TablaInodo{}
										}
										bandera = true
									}
									bit2++
								}
								Sb_System.S_free_blocks_count -= 1
								Sb_System.S_first_blo = bit2

								return inodo
							} else {
								color.Red("[/]: No hay espacio disponible")
								return in
							}
						}
					}
				} else if (pointer1.B_pointers[p1] == -1) && (aux != -1) {
					if Sb_System.S_free_blocks_count > 1 {
						pointer1.B_pointers[p1] = Sb_System.S_block_start + (Sb_System.S_first_blo * size.SizeBloqueApuntador())

						start := Sb_System.S_bm_block_start
						end := start + Sb_System.S_blocks_count
						var bit2 int32 = 0
						var bit byte
						var bandera = false
						var one byte = '1'
						//actualización sb
						for bmi := start; bmi < end; bmi++ {
							if _, err := file.Seek(int64(bmi), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return structures.TablaInodo{}
							}
							if bit == '0' && bandera {
								break
							}
							if bit == '0' && !bandera {
								if _, err := file.Seek(int64(bmi), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return structures.TablaInodo{}
								}
								if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return structures.TablaInodo{}
								}
								bandera = true
							}
							bit2++
						}
						Sb_System.S_free_blocks_count -= 1
						Sb_System.S_first_blo = bit2

						newPointer2.B_pointers[0] = Sb_System.S_block_start + (Sb_System.S_first_blo * size.SizeBloqueApuntador())

						var archivo structures.BloqueArchivo
						posBloque := Sb_System.S_block_start + (aux * size.SizeBloqueArchivo())
						archivo.B_content = DevolverContenidoArchivo(cadena)
						if _, err := file.Seek(int64(posBloque), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &archivo); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}
						CambioCont = true

						newPointer3.B_pointers[0] = posBloque
						for inicializer := 1; inicializer < 16; inicializer++ {
							newPointer2.B_pointers[inicializer] = -1
							newPointer3.B_pointers[inicializer] = -1
						}

						if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &newPointer2); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}

						if _, err := file.Seek(int64(pointer2.B_pointers[0]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &newPointer3); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}

						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return structures.TablaInodo{}
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return structures.TablaInodo{}
						}

						bit2 = 0
						bandera = false
						for bmi := start; bmi < end; bmi++ {
							if _, err := file.Seek(int64(bmi), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return structures.TablaInodo{}
							}
							if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return structures.TablaInodo{}
							}
							if bit == 0 && bandera {
								break
							}
							if bit == 0 && !bandera {
								if _, err := file.Seek(int64(bmi), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return structures.TablaInodo{}
								}
								if err := binary.Write(file, binary.LittleEndian, &one); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return structures.TablaInodo{}
								}
								bandera = false
							}
							bit2++
						}
						Sb_System.S_free_blocks_count -= 1
						Sb_System.S_first_blo = bit2

						return inodo
					} else {
						color.Red("[/]: No hay espacio disponible")
						return in
					}
				}
			}
		}
	}
	return in
}
