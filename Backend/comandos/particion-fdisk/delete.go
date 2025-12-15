package partition

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"os"
	"strings"

	"github.com/fatih/color"
)

// DeleteP elimina una partición del disco en la ruta especificada.
func DeleteP(path string, _name []byte, _unit byte, _type byte, _fit byte) {
	mbr := structures.MBR{}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Delete]: Error al abrir archivo, Disco inexistente")
		return
	}
	defer file.Close()
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("[Delete]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &mbr); err != nil {
		color.Red("[Delete]: Error en la lectura del MBR")
		return
	}
	// fmt.Print("¿Esta seguro de eliminar la particion «" + utils.ToString(_name) + "» (y/n) «Enter to Ok»: -> ")
	input := "y"
	// reader := bufio.NewReader(os.Stdin)
	// input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	pos := int32(-1)
	if input == "" || strings.ToUpper(input) == "Y" {
		// println("A borrar")
		// Situacion en la que es una particion primaria o extendida
		for i := range mbr.Mbr_partitions {
			if utils.ToString(mbr.Mbr_partitions[i].Part_name[:]) == utils.ToString(_name) {
				pos = int32(i)
				break
			} else if mbr.Mbr_partitions[i].Part_type == 'E' {
				ebr := structures.EBR{}
				if _, err := file.Seek(int64(mbr.Mbr_partitions[i].Part_start), 0); err != nil {
					color.Red("[Delete]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
					color.Red("[Delete]: Error en la lectura del EBR")
					return
				}
				if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
					if utils.ToString(ebr.Name[:]) == utils.ToString(_name) {
						ebr2 := structures.EBR{}
						ebr2.Part_next = ebr.Part_next
						ebr2.Part_start = ebr.Part_start
						ebr2.Part_s = -1
						ebr2.Part_mount = 0
						ebr2.Part_fit = 'W'
						ebr2.Name = utils.DevolverNombreByte("")

						posicionI := ebr.Part_start + size.SizeEBR()
						posicionF := ebr.Part_s + ebr.Part_start
						if ((posicionF - posicionI) / (1024 * 1024)) > 1 {
							// Eliminar MB
							EliminarMB(path, posicionI, posicionF)
						} else if ((posicionF - posicionI) / (1024)) > 1 {
							// Eliminar KB
							EliminarKB(path, posicionI, posicionF)
						} else {
							// Eliminar B
							EliminarB(path, posicionI, posicionF)
						}

						if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
							color.Red("[Delete]: Error en mover puntero")
							return
						}
						if err := binary.Write(file, binary.LittleEndian, &ebr2); err != nil {
							color.Red("[Delete]: Error en la escritura del EBR")
							return
						}
						file.Close()
						color.Blue("[Delete]: Se elimino la particion Logica «" + utils.ToString(_name) + "»")
						return
					} else if ebr.Part_next != -1 {
						antStart := ebr.Part_start
						if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
							color.Red("[Delete]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
							color.Red("[Delete]: Error en la lectura del EBR")
							return
						}
						for {
							if utils.ToString(ebr.Name[:]) == utils.ToString(_name) {
								auxEBR := structures.EBR{}
								if _, err := file.Seek(int64(antStart), 0); err != nil {
									color.Red("[Delete]: Error en mover puntero")
									return
								}
								if err := binary.Read(file, binary.LittleEndian, &auxEBR); err != nil {
									color.Red("[Delete]: Error en la lectura del EBR")
									return
								}
								auxEBR.Part_next = ebr.Part_next

								if _, err := file.Seek(int64(auxEBR.Part_start), 0); err != nil {
									color.Red("[Delete]: Error en mover puntero")
									return
								}
								if err := binary.Write(file, binary.LittleEndian, &auxEBR); err != nil {
									color.Red("[Delete]: Error en la escritura del EBR")
									return
								}

								posicionI := ebr.Part_start
								posicionF := ebr.Part_start + ebr.Part_s
								if (posicionF-posicionI)/(1024*1024) > 1 {
									EliminarMB(path, posicionI, posicionF)
								} else if (posicionF-posicionI)/(1024) > 1 {
									EliminarKB(path, posicionI, posicionF)
								} else {
									EliminarB(path, posicionI, posicionF)
								}
								file.Close()
								color.Blue("[Delete]: Se elimino la particion Logica «" + utils.ToString(_name) + "»")
								return
							}
							antStart = ebr.Part_start
							if ebr.Part_next == -1 {
								break
							}
							if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
								color.Red("[Delete]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
								color.Red("[Delete]: Error en la lectura del EBR")
								return
							}
						}
					}
				}
			}
		}
		// Particion primaria / extendida
		if pos != -1 {
			posicionI := mbr.Mbr_partitions[pos].Part_start
			posicionF := mbr.Mbr_partitions[pos].Part_start + mbr.Mbr_partitions[pos].Part_s
			if ((posicionF - posicionI) / (1024 * 1024)) > 1 {
				EliminarMB(path, posicionI, posicionF)
			} else if ((posicionF - posicionI) / (1024)) > 1 {
				EliminarKB(path, posicionI, posicionF)
			} else {
				EliminarB(path, posicionI, posicionF)
			}

			particion := structures.Partition{}
			particion.Part_s = -1
			particion.Part_start = -1
			particion.Part_type = 'P'
			particion.Part_fit = 'W'
			particion.Part_name = utils.DevolverNombreByte("")
			mbr.Mbr_partitions[pos] = particion
			if _, err := file.Seek(0, 0); err != nil {
				color.Red("[Delete]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
				color.Red("[Delete]: Error en la escritura del EBR")
				return
			}
			color.Blue("[Delete]: Se elimino la particion: <" + utils.ToString(_name) + ">")
			return
		} else {
			color.Red("[Delete]: No se encontro la particion <" + utils.ToString(_name) + ">")
			return
		}
	} else {
		color.Magenta("[Delete]: No se borro el disco")
		return
	}
}

func EliminarMB(path string, pos int32, final int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Delete]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[Delete]: Error en mover puntero")
		return
	}
	rest := (final - utils.Ftell(file)) / (1024 * 1024)
	if rest >= 1 {
		var buffer [1024]byte
		for i := range buffer {
			buffer[i] = '\x00'
		}
		for i := 0; i < (int(rest) * 1024); i++ {
			if err := binary.Write(file, binary.LittleEndian, &buffer); err != nil {
				color.Red("[Delete]: Error en eliminar disco")
				return
			}
		}
	}
	EliminarKB(path, utils.Ftell(file), final)
	file.Close()
}

func EliminarKB(path string, pos int32, final int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Delete]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[Delete]: Error en mover puntero")
		return
	}
	rest := (final - utils.Ftell(file)) / (1024)
	if rest >= 1 {
		var buffer [1024]byte
		for i := range buffer {
			buffer[i] = '\x00'
		}
		for i := 0; i < int(rest); i++ {
			if err := binary.Write(file, binary.LittleEndian, &buffer); err != nil {
				color.Red("[Delete]: Error en eliminar disco")
				return
			}
		}
	}
	EliminarB(path, utils.Ftell(file), final)
	file.Close()
}

func EliminarB(path string, pos int32, final int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Delete]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[Delete]: Error en mover puntero")
		return
	}
	rest := (final - utils.Ftell(file))
	if rest >= 1 {
		var buffer [1]byte
		buffer[0] = '\x00'
		for i := 0; i < int(rest); i++ {
			if err := binary.Write(file, binary.LittleEndian, &buffer); err != nil {
				color.Red("[Delete]: Error en eliminar disco")
				return
			}
		}
	}
	file.Close()
}
