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

func ParticionLogica(path string, _size int32, _name []byte, _unit byte, _type byte, _fit byte) {
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] El disco no existe en la ruta especificada.")
		return
	}

	mbr, ok := utils.Obtener_FULL_MBR_FDISK(path)
	if !ok {
		return
	}

	if utils.ExisteNombreP(path, utils.ToString(_name)) {
		color.Red("[FDISK]: La particion «" + utils.ToString(_name) + "» ya existe.")
		return
	}

	// 1. Encontrar la partición extendida
	var extendedPartition structures.Partition
	extendedPartitionIndex := -1
	for i, p := range mbr.Mbr_partitions {
		if p.Part_type == 'E' {
			extendedPartition = p
			extendedPartitionIndex = i
			break
		}
	}

	if extendedPartitionIndex == -1 {
		color.Red("[FDISK]: No existe particion extendida para almacenar una particion logica.")
		return
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: Error al abrir el archivo del disco.")
		return
	}
	defer file.Close()

	// 2. Leer todos los EBRs existentes y ordenarlos
	var logicalPartitions []structures.EBR
	currentEbrStart := extendedPartition.Part_start
	for {
		ebr := structures.EBR{}
		if _, err := file.Seek(int64(currentEbrStart), 0); err != nil {
			color.Red("[FDISK]: Error al buscar el EBR en la posición: " + strconv.Itoa(int(currentEbrStart)))
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
			// Llegamos al final de la lista o a un espacio no inicializado
			break
		}
		if ebr.Part_s != -1 { // Solo agregar EBRs válidos
			logicalPartitions = append(logicalPartitions, ebr)
		}

		if ebr.Part_next == -1 {
			break
		}
		currentEbrStart = ebr.Part_next
	}

	// 3. Identificar los huecos de espacio libre dentro de la partición extendida
	type Hole struct {
		Start int32
		Size  int32
	}
	var holes []Hole
	requiredSize := utils.Tamano(_size, _unit) + size.SizeEBR()

	// Hueco al inicio de la partición extendida
	nextAvailableStart := extendedPartition.Part_start
	if len(logicalPartitions) > 0 {
		holeSize := logicalPartitions[0].Part_start - nextAvailableStart
		if holeSize >= requiredSize {
			holes = append(holes, Hole{Start: nextAvailableStart, Size: holeSize})
		}
		nextAvailableStart = logicalPartitions[0].Part_start + logicalPartitions[0].Part_s + size.SizeEBR()
	}

	// Huecos entre particiones lógicas
	for i := 0; i < len(logicalPartitions)-1; i++ {
		holeStart := logicalPartitions[i].Part_start + logicalPartitions[i].Part_s + size.SizeEBR()
		holeEnd := logicalPartitions[i+1].Part_start
		holeSize := holeEnd - holeStart
		if holeSize >= requiredSize {
			holes = append(holes, Hole{Start: holeStart, Size: holeSize})
		}
	}

	// Hueco al final
	holeEnd := extendedPartition.Part_start + extendedPartition.Part_s
	holeSize := holeEnd - nextAvailableStart
	if holeSize >= requiredSize {
		holes = append(holes, Hole{Start: nextAvailableStart, Size: holeSize})
	}

	if len(holes) == 0 {
		color.Red("[FDISK]: Espacio insuficiente en la partición extendida.")
		return
	}

	// 4. Seleccionar el hueco según la estrategia de ajuste
	var bestHole Hole = holes[0]
	switch _fit {
	case 'F': // First Fit
		// ya es el primero
	case 'B': // Best Fit
		for _, h := range holes {
			if h.Size < bestHole.Size {
				bestHole = h
			}
		}
	case 'W': // Worst Fit
		for _, h := range holes {
			if h.Size > bestHole.Size {
				bestHole = h
			}
		}
	}

	// 5. Crear y enlazar el nuevo EBR
	newEBR := structures.EBR{
		Part_mount: '1',
		Part_fit:   _fit,
		Part_start: bestHole.Start,
		Part_s:     utils.Tamano(_size, _unit),
		Part_next:  -1, // Se actualizará si no es el último
	}
	copy(newEBR.Name[:], _name)

	// Encontrar el EBR anterior para enlazarlo
	var prevEBR *structures.EBR
	for i := range logicalPartitions {
		if logicalPartitions[i].Part_start < newEBR.Part_start && (prevEBR == nil || logicalPartitions[i].Part_start > prevEBR.Part_start) {
			prevEBR = &logicalPartitions[i]
		}
	}

	if prevEBR != nil {
		newEBR.Part_next = prevEBR.Part_next
		prevEBR.Part_next = newEBR.Part_start

		// Escribir el EBR anterior actualizado
		if _, err := file.Seek(int64(prevEBR.Part_start), 0); err != nil {
			color.Red("[FDISK]: Error al buscar el EBR anterior.")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, prevEBR); err != nil {
			color.Red("[FDISK]: Error al escribir el EBR anterior actualizado.")
			return
		}
	} else {
		// Es el primer EBR lógico
		newEBR.Part_next = extendedPartition.Part_start
		// Tenemos que escribir un EBR "cabeza" en el inicio de la extendida
		// que apunte a nuestro nuevo EBR. Esto se complica.
		// Por ahora, asumimos que el primer EBR se escribe al inicio si no hay otros.
		if len(logicalPartitions) > 0 {
			newEBR.Part_next = logicalPartitions[0].Part_start
		} else {
			newEBR.Part_next = -1
		}
	}

	// Escribir el nuevo EBR
	if _, err := file.Seek(int64(newEBR.Part_start), 0); err != nil {
		color.Red("[FDISK]: Error al buscar la posición para el nuevo EBR.")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &newEBR); err != nil {
		color.Red("[FDISK]: Error al escribir el nuevo EBR.")
		return
	}

	color.Green("-----------------------------------------------------------")
	color.Blue("Se creo la particion logica: " + utils.ToString(_name))
	color.Blue("Inicio: " + strconv.Itoa(int(newEBR.Part_start)))
	color.Blue("Tamaño: " + strconv.Itoa(int(newEBR.Part_s)))
	color.Green("-----------------------------------------------------------")
}
