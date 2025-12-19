package utils

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

func AgregarCarpetaSystem(posA int32, posC int32, name string) int32 {
	if !ValidarPermisoWSystem(posC, global.UsuarioLogeado.Mounted.Path) {
		color.Red("[*]: No se puede crear archivo <<" + name + ">> por falta de permisos")
	}

	nodo := global.UsuarioLogeado.Mounted
	var inodo structures.TablaInodo
	var carpeta structures.BloqueCarpeta
	var pointer1, pointer2, pointer3 structures.BloqueApuntador

	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(posC), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return -1
	}

	if inodo.I_type == 1 {
		color.Red("[/]: No es un inodo de carpeta")
		return -1
	}

	for i := 0; i < 15; i++ {
		if (inodo.I_block[i] != -1) && (i < 12) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return -1
			}
			if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return -1
			}
			for c := 0; c < 4; c++ {
				if carpeta.B_content[c].B_inodo == -1 {
					carpeta.B_content[c].B_name = NameCarpeta12(name)
					carpeta.B_content[c].B_inodo = posA
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return -1
					}
					if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
						color.Red("[/]: Error en la escritura del archivo")
						return -1
					}
					return 0
				}
			}
		} else if (inodo.I_block[i] == -1) && (i < 12) {
			if Sb_System.S_free_blocks_count > 0 {
				posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
				inodo.I_block[i] = posCarpetaO
				if _, err := file.Seek(int64(posC), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return -1
				}
				return 0
			} else {
				color.Red("[/]: No se puede crear el archivo -> " + name)
				return -1
			}
		} else if (inodo.I_block[i] == -1) && (i == 12) {
			if Sb_System.S_free_blocks_count > 1 {
				posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
				posApuntador := CrearBloqueApuntador(posCarpetaO)
				inodo.I_block[i] = posApuntador
				if _, err := file.Seek(int64(posC), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return -1
				}
				return 0
			} else {
				color.Red("[/]: No se puede crear el archivo -> " + name)
				return -1
			}
		} else if (inodo.I_block[i] != -1) && (i == 12) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return -1
			}
			if err := binary.Read(file, binary.LittleEndian, &pointer1); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return -1
			}
			for p1 := 0; p1 < 16; p1++ {
				if pointer1.B_pointers[p1] == -1 {
					if Sb_System.S_free_blocks_count > 0 {
						posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
						pointer1.B_pointers[p1] = posCarpetaO
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return -1
						}
						return 0
					} else {
						color.Red("[/]: No se puede crear el archivo -> " + name)
						return -1
					}
				} else if pointer1.B_pointers[p1] != -1 {
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return -1
					}
					if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return -1
					}
					for c := 0; c < 4; c++ {
						if carpeta.B_content[c].B_inodo == -1 {
							carpeta.B_content[c].B_name = NameCarpeta12(name)
							carpeta.B_content[c].B_inodo = posA
							if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return -1
							}
							if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
								color.Red("[/]: Error en la escritura del archivo")
								return -1
							}
							return 0
						}
					}
				}
			}
		} else if (inodo.I_block[i] == -1) && (i == 13) {
			if Sb_System.S_free_blocks_count > 2 {
				posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
				point1 := CrearBloqueApuntador(posCarpetaO)
				point2 := CrearBloqueApuntador(point1)
				inodo.I_block[i] = point2
				if _, err := file.Seek(int64(posC), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return -1
				}
				return 0
			} else {
				color.Red("[/]: No se puede crear el archivo -> " + name)
				return -1
			}
		} else if (inodo.I_block[i] != -1) && (i == 13) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return -1
			}
			if err := binary.Read(file, binary.LittleEndian, &pointer1); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return -1
			}
			for p1 := 0; p1 < 16; p1++ {
				if pointer1.B_pointers[p1] == -1 {
					if Sb_System.S_free_blocks_count > 1 {
						posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
						point1 := CrearBloqueApuntador(posCarpetaO)
						pointer1.B_pointers[p1] = point1
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return -1
						}
						return 0
					} else {
						color.Red("[/]: No se puede crear el archivo -> " + name)
						return -1
					}
				} else if pointer1.B_pointers[p1] != -1 {
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return -1
					}
					if err := binary.Read(file, binary.LittleEndian, &pointer2); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return -1
					}
					for p2 := 0; p2 < 16; p2++ {
						if pointer2.B_pointers[p2] == -1 {
							if Sb_System.S_free_blocks_count > 0 {
								posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
								pointer2.B_pointers[p2] = posCarpetaO
								if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Write(file, binary.LittleEndian, &pointer2); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return -1
								}
								return 0
							} else {
								color.Red("[/]: No se puede crear el archivo -> " + name)
								return -1
							}
						} else if pointer2.B_pointers[p2] != -1 {
							if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return -1
							}
							if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return -1
							}
							for c := 0; c < 4; c++ {
								if carpeta.B_content[c].B_inodo == -1 {
									carpeta.B_content[c].B_name = NameCarpeta12(name)
									carpeta.B_content[c].B_inodo = posA
									if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return -1
									}
									if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
										color.Red("[/]: Error en la escritura del archivo")
										return -1
									}
									return 0
								}
							}
						}
					}
				}
			}
		} else if (i == 14) && (inodo.I_block[i] == -1) {
			if Sb_System.S_free_blocks_count > 2 {
				posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
				point1 := CrearBloqueApuntador(posCarpetaO)
				point2 := CrearBloqueApuntador(point1)
				point3 := CrearBloqueApuntador(point2)
				inodo.I_block[i] = point3
				if _, err := file.Seek(int64(posC), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[/]: Error en la escritura del archivo")
					return -1
				}
				return 0
			} else {
				color.Red("[/]: No se puede crear el archivo -> " + name)
				return -1
			}
		} else if (i == 14) && (inodo.I_block[i] != -1) {
			if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
				color.Red("[/]: Error en mover puntero")
				return -1
			}
			if err := binary.Read(file, binary.LittleEndian, &pointer1); err != nil {
				color.Red("[/]: Error en la lectura del archivo")
				return -1
			}
			for p1 := 0; p1 < 16; p1++ {
				if pointer1.B_pointers[p1] == -1 {
					if (Sb_System.S_free_blocks_count > 2) && (Sb_System.S_free_inodes_count > 0) {
						posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
						point1 := CrearBloqueApuntador(posCarpetaO)
						point2 := CrearBloqueApuntador(point1)
						pointer1.B_pointers[p1] = point2
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return -1
						}
						return 0
					} else {
						color.Red("[/]: No se puede crear el archivo -> " + name)
						return -1
					}
				} else if pointer1.B_pointers[p1] != -1 {
					if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return -1
					}
					if err := binary.Read(file, binary.LittleEndian, &pointer2); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return -1
					}
					for p2 := 0; p2 < 16; p2++ {
						if pointer2.B_pointers[p2] == -1 {
							if (Sb_System.S_free_blocks_count > 1) && (Sb_System.S_free_inodes_count > 0) {
								posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
								point1 := CrearBloqueApuntador(posCarpetaO)
								pointer2.B_pointers[p2] = point1
								if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Write(file, binary.LittleEndian, &pointer2); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return -1
								}
								return 0
							} else {
								color.Red("[/]: No se puede crear el archivo -> " + name)
								return -1
							}
						} else if pointer2.B_pointers[p2] != -1 {
							if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return -1
							}
							if err := binary.Read(file, binary.LittleEndian, &pointer3); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return -1
							}
							for p3 := 0; p3 < 16; p3++ {
								if pointer3.B_pointers[p3] == -1 {
									if (Sb_System.S_free_blocks_count > 1) && (Sb_System.S_free_inodes_count > 0) {
										posCarpetaO := CrearBloqueCarpetaOtra(posA, name)
										pointer3.B_pointers[p3] = posCarpetaO
										if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return -1
										}
										if err := binary.Write(file, binary.LittleEndian, &pointer3); err != nil {
											color.Red("[/]: Error en la escritura del archivo")
											return -1
										}
										return 0
									} else {
										color.Red("[/]: No se puede crear el archivo -> " + name)
										return -1
									}
								} else if pointer3.B_pointers[p3] != -1 {
									if _, err := file.Seek(int64(pointer3.B_pointers[p3]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return -1
									}
									if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
										color.Red("[/]: Error en la lectura del archivo")
										return -1
									}
									for c := 0; c < 4; c++ {
										if carpeta.B_content[c].B_inodo == -1 {
											carpeta.B_content[c].B_name = NameCarpeta12(name)
											carpeta.B_content[c].B_inodo = posA
											if _, err := file.Seek(int64(pointer3.B_pointers[p3]), 0); err != nil {
												color.Red("[/]: Error en mover puntero")
												return -1
											}
											if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
												color.Red("[/]: Error en la escritura del archivo")
												return -1
											}
											return 0
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	color.Red("[/]: No se puede crear archivo -> " + name)
	return -1
}
