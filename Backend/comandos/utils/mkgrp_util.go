package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var CambioCont = false
var Sb_AdminUsr structures.SuperBloque

func GrupoExist(grupos []string, name string) bool {
	for _, con := range grupos {
		if strings.Contains(con, ",G,") {
			if strings.Contains(con, name) {
				return true
			}
		}
	}
	return false
}

func UsrExist(grupos []string, name string) bool {
	for _, con := range grupos {
		if strings.Contains(con, ",U,") {
			if strings.Contains(con, name) {
				return true
			}
		}
	}
	return false
}

func SplitContent(cadena string) []string {
	var split []string
	var aux string
	controlador := 0
	for i := 0; i < len(cadena); i++ {
		if controlador < 64 {
			aux += string(cadena[i])
			controlador++
		}
		if len(aux) == 64 {
			split = append(split, aux)
			aux = ""
			controlador = 0
		}
	}
	if controlador != 0 {
		split = append(split, aux)
	}
	return split
}

func GetGID(grupos []string) int32 {
	var datosG []string
	gid := 0
	for i := range grupos {
		if grupos[i] != "" {
			if strings.Contains(grupos[i], ",G,") {
				datosG = strings.Split(grupos[i], ",")
				id, _ := strconv.Atoi(datosG[0])
				if gid < id {
					gid++
				}
			}
		}
	}
	return int32(gid + 1)
}

func GetUID(grupos []string) int32 {
	var datosU []string
	uid := 0
	for i := range grupos {
		if grupos[i] != "" {
			if strings.Contains(grupos[i], ",U,") {
				datosU = strings.Split(grupos[i], ",")
				id, _ := strconv.Atoi(datosU[0])
				if uid < id {
					uid++
				}
			}
		}
	}
	return int32(uid + 1)
}

func AgregarArchivo(cadena string, inodo structures.TablaInodo, j int32, aux int32) structures.TablaInodo {
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
			seek := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
			if Sb_AdminUsr.S_free_blocks_count > 0 {
				var bit2 int32 = 0
				var bit byte
				var one byte = '1'
				start := Sb_AdminUsr.S_bm_block_start
				end := start + Sb_AdminUsr.S_blocks_count
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

				inodo.I_block[i] = Sb_AdminUsr.S_block_start + (bit2 * size.SizeBloqueApuntador())

				var archivo structures.BloqueArchivo
				posbloque := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
				Sb_AdminUsr.S_free_blocks_count -= 1
				Sb_AdminUsr.S_first_blo = bit2
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
					var posBloque = Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
			if Sb_AdminUsr.S_free_blocks_count > 1 {
				//apuntador 1
				var bit2 int32 = 0
				var bit byte
				var one byte = '1'
				start := Sb_AdminUsr.S_bm_block_start
				end := start + Sb_AdminUsr.S_blocks_count
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

				inodo.I_block[i] = Sb_AdminUsr.S_block_start + (bit2 * size.SizeBloqueApuntador())

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

				newPointer1.B_pointers[0] = Sb_AdminUsr.S_block_start + (bit2 * size.SizeBloqueApuntador())

				var archivo structures.BloqueArchivo
				posbloque := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
				Sb_AdminUsr.S_free_blocks_count -= 2
				Sb_AdminUsr.S_first_blo = bit2
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
							var posBloque = Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
					if Sb_AdminUsr.S_free_blocks_count > 0 {
						pointer1.B_pointers[p1] = Sb_AdminUsr.S_block_start + (Sb_AdminUsr.S_first_blo * size.SizeBloqueApuntador())

						var archivo structures.BloqueArchivo
						posBloque := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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

						start := Sb_AdminUsr.S_bm_block_start
						end := start + Sb_AdminUsr.S_blocks_count
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
						Sb_AdminUsr.S_free_blocks_count -= 1
						Sb_AdminUsr.S_first_blo = bit2

						return inodo
					} else {
						color.Red("[/]: No hay espacio disponible")
						return in
					}
				}
			}
		} else if (i == 14) && (inodo.I_block[i] == -1) && (aux != -1) {
			if Sb_AdminUsr.S_free_blocks_count > 2 {
				// apuntador 1
				var bit2 int32 = 0
				var bit byte
				var one byte = '1'
				start := Sb_AdminUsr.S_bm_block_start
				end := start + Sb_AdminUsr.S_blocks_count
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

				inodo.I_block[i] = Sb_AdminUsr.S_block_start + (bit2 * size.SizeBloqueApuntador())

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

				newPointer1.B_pointers[0] = Sb_AdminUsr.S_block_start + (bit2 * size.SizeBloqueApuntador())

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

				newPointer2.B_pointers[0] = Sb_AdminUsr.S_block_start + (bit2 * size.SizeBloqueApuntador())

				var archivo structures.BloqueArchivo
				posbloque := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
				Sb_AdminUsr.S_free_blocks_count -= 3
				Sb_AdminUsr.S_first_blo = bit2

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
									var posBloque = Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
							if Sb_AdminUsr.S_free_blocks_count > 0 {
								pointer2.B_pointers[p2] = Sb_AdminUsr.S_block_start + (Sb_AdminUsr.S_first_blo * size.SizeBloqueApuntador())

								var archivo structures.BloqueArchivo
								posBloque := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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

								start := Sb_AdminUsr.S_bm_block_start
								end := start + Sb_AdminUsr.S_blocks_count
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
								Sb_AdminUsr.S_free_blocks_count -= 1
								Sb_AdminUsr.S_first_blo = bit2

								return inodo
							} else {
								color.Red("[/]: No hay espacio disponible")
								return in
							}
						}
					}
				} else if (pointer1.B_pointers[p1] == -1) && (aux != -1) {
					if Sb_AdminUsr.S_free_blocks_count > 1 {
						pointer1.B_pointers[p1] = Sb_AdminUsr.S_block_start + (Sb_AdminUsr.S_first_blo * size.SizeBloqueApuntador())

						start := Sb_AdminUsr.S_bm_block_start
						end := start + Sb_AdminUsr.S_blocks_count
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
						Sb_AdminUsr.S_free_blocks_count -= 1
						Sb_AdminUsr.S_first_blo = bit2

						newPointer2.B_pointers[0] = Sb_AdminUsr.S_block_start + (Sb_AdminUsr.S_first_blo * size.SizeBloqueApuntador())

						var archivo structures.BloqueArchivo
						posBloque := Sb_AdminUsr.S_block_start + (aux * size.SizeBloqueArchivo())
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
						Sb_AdminUsr.S_free_blocks_count -= 1
						Sb_AdminUsr.S_first_blo = bit2

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

func GetContentAdminUsers(inodoStart int32) string {
	var inodo structures.TablaInodo
	var archivo structures.BloqueArchivo
	var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador
	var content string
	file, err := os.OpenFile(global.UsuarioLogeado.Mounted.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[-]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	if _, err := file.Seek(int64(inodoStart), 0); err != nil {
		color.Red("[-]: Error en mover puntero")
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[-]: Error en la lectura del archivo")
		return ""
	}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
					color.Red("[-]: Error en la lectura del archivo")
					return ""
				}
				content += ToString(archivo.B_content[:])
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[-]: Error en la lectura del archivo")
					return ""
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[-]: Error en mover puntero")
							return ""
						}
						if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
							color.Red("[-]: Error en la lectura del archivo")
							return ""
						}
						content += ToString(archivo.B_content[:])
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[-]: Error en la lectura del archivo")
					return ""
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[-]: Error en mover puntero")
							return ""
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[-]: Error en la lectura del archivo")
							return ""
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[-]: Error en mover puntero")
									return ""
								}
								if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
									color.Red("[-]: Error en la lectura del archivo")
									return ""
								}
								content += ToString(archivo.B_content[:])
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[-]: Error en la lectura del archivo")
					return ""
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[-]: Error en mover puntero")
							return ""
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[-]: Error en la lectura del archivo")
							return ""
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[-]: Error en mover puntero")
									return ""
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[-]: Error en la lectura del archivo")
									return ""
								}
								for l := 0; l < 16; l++ {
									if apuntador3.B_pointers[l] != -1 {
										if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
											color.Red("[-]: Error en mover puntero")
											return ""
										}
										if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
											color.Red("[-]: Error en la lectura del archivo")
											return ""
										}
										content += ToString(archivo.B_content[:])
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return content
}
