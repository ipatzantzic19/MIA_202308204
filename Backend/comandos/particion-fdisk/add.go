package partition

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"os"
	"strconv"

	"github.com/fatih/color"
)

// RestE reduce el tamaño de una partición en el disco en la ruta especificada.
func RestE(path string, unit byte, fit byte, add int32, name []byte) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Add-]: Error al abrir archivo")
		return
	}
	defer file.Close()

	//variables
	pos := int32(-1)
	reduccion := int32((utils.Tamano(add, unit)) * -1)
	//mbr
	mbr, embr := utils.Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return
	}

	for i := range mbr.Mbr_partitions {
		if utils.ToString(mbr.Mbr_partitions[i].Part_name[:]) == utils.ToString(name) {
			pos = int32(i)
			break
		} else if mbr.Mbr_partitions[i].Part_type == 'E' {
			ebr := structures.EBR{}
			if _, err := file.Seek(int64(mbr.Mbr_partitions[i].Part_start), 0); err != nil {
				color.Red("[Add-]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
				color.Red("[Add-]: Error en la lectura del EBR")
				return
			}
			if !((ebr.Part_s == -1) && (ebr.Part_next == -1)) {
				if utils.ToString(ebr.Name[:]) == utils.ToString(name) {
					tamanio := ebr.Part_s - size.SizeEBR()
					if tamanio > reduccion {
						ebr.Part_s = ebr.Part_s - reduccion
						if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
							color.Red("[Add-]: Error en mover puntero")
							return
						}
						if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
							color.Red("[Add-]: Error en la escritura del EBR")
							return
						}
						file.Close()
						color.Blue("[Add-]: Se quitaron " + strconv.Itoa(int(uint32(reduccion))) + "Bytes de la parcition -> " + utils.ToString(name))
						color.Blue("[Add-]: Nuevo Tamaño particion " + strconv.Itoa(int(ebr.Part_s)) + "Bytes")
						return
					} else {
						file.Close()
						color.Red("[Add-]: No se puede reducir la particion -> " + utils.ToString(name))
						return
					}
				} else if ebr.Part_next != -1 {
					if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
						color.Red("[Add-]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
						color.Red("[Add-]: Error en la lectura del EBR")
						return
					}
					for {
						if utils.ToString(ebr.Name[:]) == utils.ToString(name) {
							tamanio := ebr.Part_s - size.SizeEBR()
							if tamanio > int32(uint32(reduccion)) {
								ebr.Part_s = ebr.Part_s - reduccion
								if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
									color.Red("[Add-]: Error en mover puntero")
									return
								}
								if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
									color.Red("[Add-]: Error en la escritura del EBR")
									return
								}
								file.Close()
								color.Blue("[Add-]: Se quitaron " + strconv.Itoa(int(uint32(reduccion))) + "Bytes de la parcition -> " + utils.ToString(name))
								color.Blue("[Add-]: Nuevo Tamaño particion " + strconv.Itoa(int(ebr.Part_s)) + "Bytes")
								return
							} else {
								file.Close()
								color.Red("[Add-]: No se puede reducir la particion -> " + utils.ToString(name))
								return
							}
						}
						if ebr.Part_next == -1 {
							break
						}
						if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
							color.Red("[Add-]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
							color.Red("[Add-]: Error en la lectura del EBR")
							return
						}
					}
				}
			}
		}
	}
	if pos != -1 {
		tamanio := mbr.Mbr_partitions[pos].Part_s
		if tamanio > reduccion {
			mbr.Mbr_partitions[pos].Part_s = tamanio - reduccion
			if _, err := file.Seek(0, 0); err != nil {
				color.Red("[Add-]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
				color.Red("[Add-]: Error en la escritura del mbr")
				return
			}
			file.Close()
			color.Blue("[Add-]: Se quitaron " + strconv.Itoa(int(uint32(reduccion))) + "Bytes de la parcition -> " + utils.ToString(name))
			color.Blue("[Add-]: Nuevo Tamaño particion " + strconv.Itoa(int(mbr.Mbr_partitions[pos].Part_s)) + "Bytes")
			return
		} else {
			file.Close()
			color.Red("[Add-]: No se puede reducir la particion -> " + utils.ToString(name))
			return
		}
	} else {
		color.Red("[Add-]: No existe la particion")
		return
	}
}

// AddE aumenta el tamaño de una partición en el disco en la ruta especificada.
func AddE(path string, unit byte, fit byte, add int32, name []byte) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Add+]: Error al abrir archivo")
		return
	}
	defer file.Close()

	//variables
	pos := int32(-1)
	agregar := int32((utils.Tamano(add, unit)))
	//mbr
	mbr, embr := utils.Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return
	}

	for i := range mbr.Mbr_partitions {
		if utils.ToString(mbr.Mbr_partitions[i].Part_name[:]) == utils.ToString(name) {
			pos = int32(i)
			break
		} else if mbr.Mbr_partitions[i].Part_type == 'E' {
			ebr := structures.EBR{}
			if _, err := file.Seek(int64(mbr.Mbr_partitions[i].Part_start), 0); err != nil {
				color.Red("[Add+]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
				color.Red("[Add+]: Error en la lectura del EBR")
				return
			}
			if !((ebr.Part_s == -1) && (ebr.Part_next == -1)) {
				if utils.ToString(ebr.Name[:]) == utils.ToString(name) {
					if ebr.Part_next != -1 {
						if (ebr.Part_s + ebr.Part_start) >= ebr.Part_next {
							file.Close()
							color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
							return
						} else if (ebr.Part_start + ebr.Part_s + agregar) < ebr.Part_next {
							ebr.Part_s += agregar
							if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
								color.Red("[Add+]: Error en mover puntero")
								return
							}
							if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
								color.Red("[Add+]: Error en la escritura del EBR")
								return
							}
							file.Close()
							color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
							return
						} else {
							file.Close()
							color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
							return
						}
					} else {
						if (ebr.Part_start + ebr.Part_s + agregar) <= (mbr.Mbr_partitions[i].Part_start + mbr.Mbr_partitions[i].Part_s) {
							ebr.Part_s += agregar
							if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
								color.Red("[Add+]: Error en mover puntero")
								return
							}
							if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
								color.Red("[Add+]: Error en la escritura del EBR")
								return
							}
							file.Close()
							color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
							return
						} else {
							file.Close()
							color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
							return
						}
					}
				} else if ebr.Part_next != -1 {
					if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
						color.Red("[Add+]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
						color.Red("[Add+]: Error en la lectura del EBR")
						return
					}
					for { //while true
						if utils.ToString(ebr.Name[:]) == utils.ToString(name) {
							if ebr.Part_next != -1 {
								if (ebr.Part_start + ebr.Part_s) >= ebr.Part_next {
									file.Close()
									color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
									return
								} else if (ebr.Part_start + ebr.Part_s + agregar) < ebr.Part_next {
									ebr.Part_s += agregar
									if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
										color.Red("[Add+]: Error en mover puntero")
										return
									}
									if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
										color.Red("[Add+]: Error en la escritura del EBR")
										return
									}
									file.Close()
									color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
									return
								} else {
									file.Close()
									color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
									return
								}
							} else {
								if (ebr.Part_start + ebr.Part_s + agregar) <= (mbr.Mbr_partitions[i].Part_start + mbr.Mbr_partitions[i].Part_s) {
									ebr.Part_s += agregar
									if _, err := file.Seek(int64(ebr.Part_start), 0); err != nil {
										color.Red("[Add+]: Error en mover puntero")
										return
									}
									if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
										color.Red("[Add+]: Error en la escritura del EBR")
										return
									}
									file.Close()
									color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
									return
								} else {
									file.Close()
									color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
									return
								}
							}
						}
						if ebr.Part_next == -1 {
							break
						}
						if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
							color.Red("[Add+]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
							color.Red("[Add+]: Error en la lectura del EBR")
							return
						}
					}
				}
			}
		}
	}
	if pos != -1 {
		if pos == 3 {
			if (mbr.Mbr_partitions[pos].Part_start + mbr.Mbr_partitions[pos].Part_s + agregar) <= mbr.Mbr_tamano {
				mbr.Mbr_partitions[pos].Part_s += agregar
				if _, err := file.Seek(0, 0); err != nil {
					color.Red("[Add+]: Error en mover puntero")
					return
				}
				if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
					color.Red("[Add+]: Error en la escritura del EBR")
					return
				}
				file.Close()
				color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
				return
			} else {
				file.Close()
				color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
				return
			}
		} else {
			//caso en el que la particion no haya sido instancia
			if mbr.Mbr_partitions[pos+1].Part_start == -1 {
				if (mbr.Mbr_partitions[pos].Part_start + mbr.Mbr_partitions[pos].Part_s + agregar) <= mbr.Mbr_tamano {
					mbr.Mbr_partitions[pos].Part_s += agregar
					if _, err := file.Seek(0, 0); err != nil {
						color.Red("[Add+]: Error en mover puntero")
						return
					}
					if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
						color.Red("[Add+]: Error en la escritura del EBR")
						return
					}
					file.Close()
					color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
					return
				} else {
					file.Close()
					color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
					return
				}
			} else {
				if (mbr.Mbr_partitions[pos].Part_start + mbr.Mbr_partitions[pos].Part_s) >= mbr.Mbr_partitions[pos+1].Part_start {
					file.Close()
					color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
					return
				} else if (mbr.Mbr_partitions[pos].Part_start + mbr.Mbr_partitions[pos].Part_s + agregar) < mbr.Mbr_partitions[pos+1].Part_start {
					mbr.Mbr_partitions[pos].Part_s += agregar
					if _, err := file.Seek(0, 0); err != nil {
						color.Red("[Add+]: Error en mover puntero")
						return
					}
					if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
						color.Red("[Add+]: Error en la escritura del EBR")
						return
					}
					file.Close()
					color.Blue("[Add+]: Se Aumento '" + strconv.Itoa(int(agregar)) + "' Bytes el tamaño de la particion -> " + utils.ToString(name))
					return
				} else {
					file.Close()
					color.Red("[Add+]: No hay espacio suficiente para aumentar la particion")
					return
				}
			}
		}
	} else {
		color.Red("[Add+]: No existe la particion -> " + utils.ToString(name))
		return
	}
}
