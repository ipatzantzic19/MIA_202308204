package utils

import (
	"Proyecto/Estructuras/structures"
	"encoding/binary"
	"os"
	"strings"

	"github.com/fatih/color"
)

func SplitRuta(ruta string) []string {
	var split []string
	aux := strings.Split(ruta, "/")
	for _, s := range aux {
		if s != "" {
			split = append(split, s)
		}
	}
	return split
}

func GetInodoF(rutaS []string, posAct int32, rutaSize int32, start int32, path string) int32 {
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
