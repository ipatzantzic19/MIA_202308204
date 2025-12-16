package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

// Devuelve una particion
func BuscarParticion(disco structures.MBR, nombre []byte, path string) ([]interface{}, bool) {
	Devolucion := make([]interface{}, 2)
	tempDisk := PartitionVacia()
	inicio := int32(0)
	// error := false
	es_primaria_extendida := false
	for i := range disco.Mbr_partitions {
		if (ToString(disco.Mbr_partitions[i].Part_name[:]) == ToString(nombre)) && (ToString(nombre) != "") {
			tempDisk = disco.Mbr_partitions[i]
			Devolucion[0] = tempDisk
			es_primaria_extendida = true
			if i == 0 {
				inicio = size.SizeMBR_NotPartitions()
				Devolucion[1] = inicio
			} else if i == 1 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition()
				Devolucion[1] = inicio
			} else if i == 2 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else if i == 3 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else {
				inicio = 0
				// error = true
			}
			break
		} else {
			if ToString(nombre) == "" {
				color.Yellow("[Mount]: Particion sin nombre")
				Devolucion[0] = nil
				Devolucion[1] = nil
				return Devolucion, false
			}
		}
	}

	if es_primaria_extendida {
		return Devolucion, true
	}

	//caso de ser particion logica
	for _, log := range disco.Mbr_partitions {
		if log.Part_type == 'E' {
			tempDisk = log
			break
		}
	}

	// siguiente := tempDisk.Part_start
	ebr := structures.EBR{}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Mount]: Error al abrir archivo")
		Devolucion[0] = nil
		Devolucion[1] = nil
		return Devolucion, false
	}
	defer file.Close()
	if _, err := file.Seek(int64(tempDisk.Part_start), 0); err != nil {
		color.Red("[Mount]: Error en mover puntero")
		Devolucion[0] = nil
		Devolucion[1] = nil
		return Devolucion, false
	}
	if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
		color.Red("[Mount]: Error en la lectura del EBR")
		Devolucion[0] = nil
		Devolucion[1] = nil
		return Devolucion, false
	}

	if ebr.Part_next != -1 || ebr.Part_s != -1 {
		if ToString(ebr.Name[:]) == ToString(nombre) {
			Devolucion[0] = ebr
			Devolucion[1] = ebr.Part_start
			return Devolucion, true
		}
		for ebr.Part_next != -1 {
			if ToString(ebr.Name[:]) == ToString(nombre) {
				Devolucion[0] = ebr
				Devolucion[1] = ebr.Part_start
				return Devolucion, true
			}
			if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
				color.Red("[Mount]: Error en mover puntero")
				Devolucion[0] = nil
				Devolucion[1] = nil
				return Devolucion, false
			}
			if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
				color.Red("[Mount]: Error en la lectura del EBR")
				Devolucion[0] = nil
				Devolucion[1] = nil
				return Devolucion, false
			}
			if ToString(ebr.Name[:]) == ToString(nombre) {
				Devolucion[0] = ebr
				Devolucion[1] = ebr.Part_start
				return Devolucion, true
			}
		}
	}
	Devolucion[0] = nil
	Devolucion[1] = nil
	return Devolucion, false
}
