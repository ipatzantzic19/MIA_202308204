package utils

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
)

func Ljust(s string, leng int) string {
	if len(s) >= leng {
		return s
	}
	return s + strings.Repeat(" ", leng-len(s))
}

func TieneNameRep(valor string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(valor), "name=") {
		color.Red("[REP]: No tiene name o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[REP]: No tiene grp Valido")
		return "", false
	}
	if (strings.ToLower(value[1]) == "mbr") || (strings.ToLower(value[1]) == "disk") || (strings.ToLower(value[1]) == "inode") || (strings.ToLower(value[1]) == "journaling") || (strings.ToLower(value[1]) == "block") || (strings.ToLower(value[1]) == "bm_inode") || (strings.ToLower(value[1]) == "bm_block") || (strings.ToLower(value[1]) == "tree") || (strings.ToLower(value[1]) == "sb") || (strings.ToLower(value[1]) == "file") || (strings.ToLower(value[1]) == "ls") {
		return strings.ToLower(value[1]), true
	}
	return "", false
}

func TienePathRep(valor string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(valor), "path=") {
		color.Red("[REP]: No tiene path o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[REP]: No tiene path Valido")
		return "", false
	}
	return value[1], true
}

func TieneIDRep(valor string) (string, bool) {
	comando := "REP"
	if !strings.HasPrefix(strings.ToLower(valor), "id=") {
		color.Red("[" + comando + "]: No tiene id o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene id Valido")
		return "", false
	}
	if len(value[1]) > 4 {
		color.Red("[" + comando + "]: No tiene id Valido")
		return "", false
	}
	return value[1], true
}

func TieneRutaRep(valor string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(valor), "ruta=") {
		color.Red("[REP]: No tiene ruta o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[REP]: No tiene ruta Valida")
		return "", false
	}
	return value[1], true
}

func Returnstring(s string) string {
	if len(s) <= 0 || !(s != "") {
		return ""
	}
	return s
}

func ObtenerDiscoID(discoID string) (global.ParticionesMontadas, bool) {
	for _, v := range global.Mounted_Partitions {
		if ToString(v.ID_Particion[:]) == ToString([]byte(discoID)) {
			return v, true
		}
	}
	return global.ParticionesMontadas{}, false
}

func Redondeo(num float32) float64 {
	pow := math.Pow(10, float64(2))
	rounted := math.Round(float64(num)*pow) / pow
	return rounted
}

func GetContentReport(inodoStart int32, path string) string {
	var inodo structures.TablaInodo
	var archivo structures.BloqueArchivo
	var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador
	var content string
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
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
