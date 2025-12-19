package utils

import (
	"Proyecto/Estructuras/structures"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

func GetInodoFSystem(rutaS []string, posAct int32, rutaSize int32, start int32, path string) int32 {
	inodo := structures.TablaInodo{}
	carpeta := structures.BloqueCarpeta{}
	apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return -1
	}

	if inodo.I_type == 1 {
		color.Red("[/]: No es inodo de carpeta")
		return -1
	}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 { //antes de ser indirecto
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
					color.Red("[/]: Error en la lectura del Bloque de Carpeta")
					return -1
				}
				for c := 0; c < 4; c++ {
					if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
						if posAct < rutaSize {
							return GetInodoF(rutaS, posAct+1, rutaSize, carpeta.B_content[c].B_inodo, path)
						}
						if posAct == rutaSize {
							return carpeta.B_content[c].B_inodo
						}
					}
				}
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del Bloque de Apuntadores 1")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[/]: Error en la lectura del Bloque de Carpeta")
							return -1
						}
						for c := 0; c < 4; c++ {
							if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
								if posAct < rutaSize {
									return GetInodoF(rutaS, posAct+1, rutaSize, carpeta.B_content[c].B_inodo, path)
								}
								if posAct == rutaSize {
									return carpeta.B_content[c].B_inodo
								}
							}
						}
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del Bloque de Apuntadores 1")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[/]: Error en la lectura del Bloque de Apuntadores 2")
							return -1
						}
						for k := 0; k < 16; k++ {
							if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return -1
							}
							if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
								color.Red("[/]: Error en la lectura del Bloque de Carpeta")
								return -1
							}
							for c := 0; c < 4; c++ {
								if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
									if posAct < rutaSize {
										return GetInodoF(rutaS, posAct+1, rutaSize, carpeta.B_content[c].B_inodo, path)
									}
									if posAct == rutaSize {
										return carpeta.B_content[c].B_inodo
									}
								}
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del Bloque de Apuntadores 1")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[/]: Error en la lectura del Bloque de Apuntadores 2")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[/]: Error en la lectura del Bloque de Apuntadores 3")
									return -1
								}
								for l := 0; l < 16; l++ {
									if apuntador3.B_pointers[l] != -1 {
										if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return -1
										}
										if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[/]: Error en la lectura del Bloque de Carpeta")
											return -1
										}
										for c := 0; c < 4; c++ {
											if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
												if posAct < rutaSize {
													return GetInodoF(rutaS, posAct+1, rutaSize, carpeta.B_content[c].B_inodo, path)
												}
												if posAct == rutaSize {
													return carpeta.B_content[c].B_inodo
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
		}
	}
	return -1
}

func ExistPathSystem(rutaS []string, posAct int32, start int32, path string) int32 {
	var inodo structures.TablaInodo
	var carpeta structures.BloqueCarpeta
	var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return -1
	}

	if inodo.I_type == 1 {
		color.Red("[/]: No es inodo de carpeta")
		return -1
	}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
					color.Red("[/]: Error en la lectura del Bloque de Apuntadores 1")
					return -1
				}
				for c := 0; c < 4; c++ {
					if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
						return carpeta.B_content[c].B_inodo
					}
				}
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[/]: Error en la lectura del archivo")
							return -1
						}
						for c := 0; c < 16; c++ {
							if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
								return carpeta.B_content[c].B_inodo
							}
						}
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[/]: Error en la lectura del archivo")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
									color.Red("[/]: Error en la lectura del archivo")
									return -1
								}
								for c := 0; c < 4; c++ {
									if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
										return carpeta.B_content[c].B_inodo
									}
								}
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[/]: Error en la lectura del archivo")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[/]: Error en la lectura del archivo")
									return -1
								}
								for l := 0; l < 16; l++ {
									if apuntador3.B_pointers[l] != -1 {
										if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return -1
										}
										if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[/]: Error en la lectura del archivo")
											return -1
										}
										for c := 0; c < 4; c++ {
											if ToString(carpeta.B_content[c].B_name[:]) == rutaS[posAct] {
												return carpeta.B_content[c].B_inodo
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
	}
	return start
}

func ReturnCarpetaSystem(nameP string, start int32, path string) int32 {
	inodo := structures.TablaInodo{}
	carpeta := structures.BloqueCarpeta{}
	apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return -1
	}

	if inodo.I_type == 1 {
		color.Red("[/]: No es inodo de carpeta")
		return -1
	}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for c := 0; c < 4; c++ {
					if ToString(carpeta.B_content[c].B_name[:]) == nameP {
						return inodo.I_block[i]
					}
				}
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[/]: Error en la lectura del archivo")
							return -1
						}
						for c := 0; c < 4; c++ {
							if ToString(carpeta.B_content[c].B_name[:]) == nameP {
								return apuntador1.B_pointers[j]
							}
						}
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[/]: Error en la lectura del archivo")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
									color.Red("[/]: Error en la lectura del archivo")
									return -1
								}
								for c := 0; c < 4; c++ {
									if ToString(carpeta.B_content[c].B_name[:]) == nameP {
										return apuntador2.B_pointers[k]
									}
								}
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[/]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[/]: Error en la lectura del archivo")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[/]: Error en la lectura del archivo")
									return -1
								}
								for l := 0; l < 16; l++ {
									if apuntador3.B_pointers[l] != -1 {
										if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return -1
										}
										if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[/]: Error en la lectura del archivo")
											return -1
										}
										for c := 0; c < 4; c++ {
											if ToString(carpeta.B_content[c].B_name[:]) == nameP {
												return apuntador3.B_pointers[l]
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
	}
	return -1
}

func CrearCarpetaSystem(rutaS []string, posAct int32, posI int32, path string) int32 {
	if !ValidarPermisoWSystem(posI, path) {
		color.Red("[/]: No se puede crear la carpeta <<'" + rutaS[posAct] + "'>> por falta de permisos")
		return -1
	}

	var inodo structures.TablaInodo
	var carpeta structures.BloqueCarpeta
	var pointer1, pointer2, pointer3 structures.BloqueApuntador
	// var posNewCar int32 = 0
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return -1
	}
	defer file.Close()
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return -1
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del inodo")
		return -1
	}
	if inodo.I_type == 1 {
		color.Red("[/]: No es inodo de carpeta")
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
					if (Sb_System.S_free_inodes_count > 0) && (Sb_System.S_free_blocks_count > 0) {
						posInodo := BuscarPosicionNewInodo()
						posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
						CrearInodoCarpeta(posInodo, posCarpetaI)
						carpeta.B_content[c].B_name = NameCarpeta12(rutaS[posAct])
						carpeta.B_content[c].B_inodo = posInodo
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[-]: Error en la escritura del archivo")
							return -1
						}
						return posInodo
					} else {
						color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
						return -1
					}
				}
			}
		} else if (inodo.I_block[i] == -1) && (i < 12) {
			if (Sb_System.S_free_blocks_count > 1) && (Sb_System.S_free_inodes_count > 0) {
				posInodo := BuscarPosicionNewInodo()
				posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
				CrearInodoCarpeta(posInodo, posCarpetaI)
				posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
				inodo.I_block[i] = posCarpetaO
				if _, err := file.Seek(int64(posI), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[-]: Error en la escritura del archivo")
					return -1
				}
				return posInodo
			} else {
				color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
				return -1
			}
		} else if (i == 12) && (inodo.I_block[i] == -1) {
			if (Sb_System.S_free_blocks_count > 2) && (Sb_System.S_free_inodes_count > 0) {
				posInodo := BuscarPosicionNewInodo()
				posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
				CrearInodoCarpeta(posInodo, posCarpetaI)
				posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
				posApuntador := CrearBloqueApuntador(posCarpetaO)
				inodo.I_block[i] = posApuntador
				if _, err := file.Seek(int64(posI), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[-]: Error en la escritura del archivo")
					return -1
				}
				return posInodo
			} else {
				color.Red("[/]: No hay suficientes bloques")
				return -1
			}
		} else if (i == 12) && (inodo.I_block[i] != -1) {
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
					if (Sb_System.S_free_blocks_count > 1) && (Sb_System.S_free_inodes_count > 0) {
						posInodo := BuscarPosicionNewInodo()
						posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
						CrearInodoCarpeta(posInodo, posCarpetaI)
						posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
						pointer1.B_pointers[p1] = posCarpetaO
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[-]: Error en la escritura del archivo")
							return -1
						}
						return posInodo
					} else {
						color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
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
							if (Sb_System.S_free_inodes_count > 0) && (Sb_System.S_free_blocks_count > 0) {
								posInodo := BuscarPosicionNewInodo()
								posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
								CrearInodoCarpeta(posInodo, posCarpetaI)
								carpeta.B_content[c].B_name = NameCarpeta12(rutaS[posAct])
								carpeta.B_content[c].B_inodo = posInodo
								if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
									color.Red("[-]: Error en la escritura del archivo")
									return -1
								}
								return posInodo
							} else {
								color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
								return -1
							}
						}
					}
				}
			}
		} else if (i == 13) && (inodo.I_block[i] == -1) {
			if (Sb_System.S_free_blocks_count > 3) && (Sb_System.S_free_inodes_count > 0) {
				posInodo := BuscarPosicionNewInodo()
				posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
				CrearInodoCarpeta(posInodo, posCarpetaI)
				posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
				point1 := CrearBloqueApuntador(posCarpetaO)
				point2 := CrearBloqueApuntador(point1)
				inodo.I_block[i] = point2
				if _, err := file.Seek(int64(posI), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[-]: Error en la escritura del archivo")
					return -1
				}
				return posInodo
			} else {
				color.Red("[/]: No hay suficientes bloques")
				return -1
			}
		} else if (i == 13) && (inodo.I_block[i] != -1) {
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
						posInodo := BuscarPosicionNewInodo()
						posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
						CrearInodoCarpeta(posInodo, posCarpetaI)
						posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
						point1 := CrearBloqueApuntador(posCarpetaO)
						pointer1.B_pointers[p1] = point1
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[-]: Error en la escritura del archivo")
							return -1
						}
						return posInodo
					} else {
						color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
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
								posInodo := BuscarPosicionNewInodo()
								posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
								CrearInodoCarpeta(posInodo, posCarpetaI)
								posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
								pointer2.B_pointers[p2] = posCarpetaO
								if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Write(file, binary.LittleEndian, &pointer2); err != nil {
									color.Red("[-]: Error en la escritura del archivo")
									return -1
								}
								return posInodo
							} else {
								color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
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
									if (Sb_System.S_free_inodes_count > 0) && (Sb_System.S_free_blocks_count > 0) {
										posInodo := BuscarPosicionNewInodo()
										posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
										CrearInodoCarpeta(posInodo, posCarpetaI)
										carpeta.B_content[c].B_name = NameCarpeta12(rutaS[posAct])
										carpeta.B_content[c].B_inodo = posInodo
										if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return -1
										}
										if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[-]: Error en la escritura del archivo")
											return -1
										}
										return posInodo
									} else {
										color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
										return -1
									}
								}
							}
						}
					}
				}
			}
		} else if (i == 14) && (inodo.I_block[i] == -1) {
			if (Sb_System.S_free_blocks_count > 4) && (Sb_System.S_free_inodes_count > 0) {
				posInodo := BuscarPosicionNewInodo()
				posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
				CrearInodoCarpeta(posInodo, posCarpetaI)
				posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
				point1 := CrearBloqueApuntador(posCarpetaO)
				point2 := CrearBloqueApuntador(point1)
				point3 := CrearBloqueApuntador(point2)
				inodo.I_block[i] = point3
				if _, err := file.Seek(int64(posI), 0); err != nil {
					color.Red("[/]: Error en mover puntero")
					return -1
				}
				if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
					color.Red("[-]: Error en la escritura del archivo")
					return -1
				}
				return posInodo
			} else {
				color.Red("[/]: No hay suficientes bloques")
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
					if (Sb_System.S_free_blocks_count > 3) && (Sb_System.S_free_inodes_count > 0) {
						posInodo := BuscarPosicionNewInodo()
						posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
						CrearInodoCarpeta(posInodo, posCarpetaI)
						posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
						point1 := CrearBloqueApuntador(posCarpetaO)
						point2 := CrearBloqueApuntador(point1)
						pointer1.B_pointers[p1] = point2
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[/]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &pointer1); err != nil {
							color.Red("[-]: Error en la escritura del archivo")
							return -1
						}
						return posInodo
					} else {
						color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
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
							if (Sb_System.S_free_blocks_count > 2) && (Sb_System.S_free_inodes_count > 0) {
								posInodo := BuscarPosicionNewInodo()
								posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
								CrearInodoCarpeta(posInodo, posCarpetaI)
								posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
								point1 := CrearBloqueApuntador(posCarpetaO)
								pointer2.B_pointers[p2] = point1
								if _, err := file.Seek(int64(pointer1.B_pointers[p1]), 0); err != nil {
									color.Red("[/]: Error en mover puntero")
									return -1
								}
								if err := binary.Write(file, binary.LittleEndian, &pointer2); err != nil {
									color.Red("[-]: Error en la escritura del archivo")
									return -1
								}
								return posInodo
							} else {
								color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
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
										posInodo := BuscarPosicionNewInodo()
										posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
										CrearInodoCarpeta(posInodo, posCarpetaI)
										posCarpetaO := CrearBloqueCarpetaOtra(posInodo, rutaS[posAct])
										pointer3.B_pointers[p3] = posCarpetaO
										if _, err := file.Seek(int64(pointer2.B_pointers[p2]), 0); err != nil {
											color.Red("[/]: Error en mover puntero")
											return -1
										}
										if err := binary.Write(file, binary.LittleEndian, &pointer3); err != nil {
											color.Red("[-]: Error en la escritura del archivo")
											return -1
										}
										return posInodo
									} else {
										color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
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
											if (Sb_System.S_free_inodes_count > 0) && (Sb_System.S_free_blocks_count > 0) {
												posInodo := BuscarPosicionNewInodo()
												posCarpetaI := CrearBloqueCarpetaInicial(posInodo, posI)
												CrearInodoCarpeta(posInodo, posCarpetaI)
												carpeta.B_content[c].B_name = NameCarpeta12(rutaS[posAct])
												carpeta.B_content[c].B_inodo = posInodo
												if _, err := file.Seek(int64(pointer3.B_pointers[p3]), 0); err != nil {
													color.Red("[/]: Error en mover puntero")
													return -1
												}
												if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
													color.Red("[-]: Error en la escritura del archivo")
													return -1
												}
												return posInodo
											} else {
												color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
												return -1
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
	}
	color.Red("[/]: No se puede crear la carpeta -> " + rutaS[posAct])
	return -1

}
