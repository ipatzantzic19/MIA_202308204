package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"encoding/binary"
	"io"
	"os"

	"github.com/fatih/color"
)

// Obtener_FULL_MBR_FDISK lee y devuelve el MBR completo desde el archivo de disco especificado por la ruta.
func Obtener_FULL_MBR_FDISK(path string) (structures.MBR, bool) {
	mbr := structures.MBR{}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: Error al abrir archivo")
		return structures.MBR{}, false
	}
	defer file.Close()
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("[FDISK]: Error en mover puntero")
		return structures.MBR{}, false
	}
	if err := binary.Read(file, binary.LittleEndian, &mbr); err != nil {
		color.Red("[FDISK]: Error en la lectura del MBR")
		return structures.MBR{}, false
	}
	return mbr, true
}

// EspacioDisponible verifica si hay suficiente espacio disponible en el disco para una nueva partición.
func EspacioDisponible(s int32, path string, u byte, pos int32) bool {
	mbr, embr := Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return false
	}

	if pos > -1 {
		if Tamano(s, u) > 0 {
			espacioRestante := 0
			if pos == 0 {
				espacioRestante = int(mbr.Mbr_tamano) - int(size.SizeMBR())
			} else {
				espacioRestante = int(mbr.Mbr_tamano) - int(mbr.Mbr_partitions[pos-1].Part_start) - int(mbr.Mbr_partitions[pos-1].Part_s)
			}
			return espacioRestante >= int(Tamano(s, u))
		}
	}
	return false
}

// ExisteNombreP verifica si ya existe una partición con el nombre dado en el disco especificado por la ruta.
func ExisteNombreP(path string, name string) bool {
	mbr, embr := Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return true
	}

	for i := range mbr.Mbr_partitions {
		if ToString(mbr.Mbr_partitions[i].Part_name[:]) == name {
			return true
		}
		if mbr.Mbr_partitions[i].Part_type == 'E' {
			EBR := structures.EBR{}
			file, err := os.OpenFile(path, os.O_RDWR, 0666)
			if err != nil {
				color.Red("[*]: Error al abrir archivo")
				return true
			}
			defer file.Close()
			if _, err := file.Seek(int64(mbr.Mbr_partitions[i].Part_start), 0); err != nil {
				color.Red("[*]: Error en mover puntero")
				return true
			}
			if err := binary.Read(file, binary.LittleEndian, &EBR); err != nil {
				color.Red("[*]: Error en la lectura del MBR")
				return true
			}
			if EBR.Part_next != -1 || EBR.Part_s != -1 {
				if ToString(EBR.Name[:]) == name {
					return true
				}
				for EBR.Part_next != -1 {
					if ToString(EBR.Name[:]) == name {
						return true
					}
					if _, err := file.Seek(int64(EBR.Part_next), 0); err != nil {
						color.Red("[*]: Error en mover puntero")
						return true
					}
					if err := binary.Read(file, binary.LittleEndian, &EBR); err != nil {
						color.Red("[*]: Error en la lectura del MBR")
						return true
					}
					//si la particion que le sigue
					if ToString(EBR.Name[:]) == ToString([]byte(name)) {
						return true
					}
				}
			}
		}
	}
	return false
}

// ExisteParticionExt verifica si existe una partición extendida en el disco especificado por la ruta.
func ExisteParticionExt(path string) bool {
	mbr, embr := Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return false
	}

	for i := range mbr.Mbr_partitions {
		if mbr.Mbr_partitions[i].Part_type == 'E' {
			return true
		}
	}
	return false
}

// Ftell obtiene la posición actual del puntero de un archivo y la devuelve como un entero de 32 bits.
func Ftell(file *os.File) int32 {
	pos, _ := file.Seek(0, io.SeekCurrent)
	return int32(pos)
}

// SortPartitions ordena una lista de particiones basándose en su Part_start.
func SortPartitions(partitions [4]structures.Partition) []structures.Partition {
	var sortedPartitions []structures.Partition
	for _, p := range partitions {
		if p.Part_start != -1 { // -1 significa que no está en uso
			sortedPartitions = append(sortedPartitions, p)
		}
	}

	for i := 0; i < len(sortedPartitions)-1; i++ {
		for j := 0; j < len(sortedPartitions)-i-1; j++ {
			if sortedPartitions[j].Part_start > sortedPartitions[j+1].Part_start {
				sortedPartitions[j], sortedPartitions[j+1] = sortedPartitions[j+1], sortedPartitions[j]
			}
		}
	}
	return sortedPartitions
}

// EncontrarAjuste encuentra la mejor posición de inicio para una nueva partición
// basándose en la estrategia de ajuste especificada (First, Best, Worst Fit).
func EncontrarAjuste(mbr structures.MBR, requiredSize int32, fit byte) int32 {
	// 1. Obtener y ordenar las particiones existentes
	sortedPartitions := SortPartitions(mbr.Mbr_partitions)

	// 2. Identificar los huecos de espacio libre
	type Hole struct {
		Start int32
		Size  int32
	}
	var holes []Hole

	// Hueco antes de la primera partición
	firstPartitionStart := mbr.Mbr_tamano
	if len(sortedPartitions) > 0 {
		firstPartitionStart = sortedPartitions[0].Part_start
	}
	holeStart := size.SizeMBR()
	holeSize := firstPartitionStart - holeStart
	if holeSize >= requiredSize {
		holes = append(holes, Hole{Start: holeStart, Size: holeSize})
	}

	// Huecos entre particiones
	for i := 0; i < len(sortedPartitions)-1; i++ {
		holeStart = sortedPartitions[i].Part_start + sortedPartitions[i].Part_s
		holeEnd := sortedPartitions[i+1].Part_start
		holeSize = holeEnd - holeStart
		if holeSize >= requiredSize {
			holes = append(holes, Hole{Start: holeStart, Size: holeSize})
		}
	}

	// Hueco después de la última partición
	if len(sortedPartitions) > 0 {
		lastPartition := sortedPartitions[len(sortedPartitions)-1]
		holeStart = lastPartition.Part_start + lastPartition.Part_s
		holeEnd := mbr.Mbr_tamano
		holeSize = holeEnd - holeStart
		if holeSize >= requiredSize {
			holes = append(holes, Hole{Start: holeStart, Size: holeSize})
		}
	}

	if len(holes) == 0 {
		return -1 // No hay espacio disponible
	}

	// 3. Seleccionar el hueco según la estrategia de ajuste
	var bestHole Hole

	switch fit {
	case 'F': // First Fit
		bestHole = holes[0]
	case 'B': // Best Fit
		bestHole = holes[0]
		for _, h := range holes {
			if h.Size < bestHole.Size {
				bestHole = h
			}
		}
	case 'W': // Worst Fit
		bestHole = holes[0]
		for _, h := range holes {
			if h.Size > bestHole.Size {
				bestHole = h
			}
		}
	default: // Por defecto, Worst Fit
		bestHole = holes[0]
		for _, h := range holes {
			if h.Size > bestHole.Size {
				bestHole = h
			}
		}
	}

	return bestHole.Start
}
