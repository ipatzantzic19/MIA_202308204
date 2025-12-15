package partition

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func ParticionExtendida(path string, _size int32, _name []byte, _unit byte, _type byte, _fit byte) {
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] Disco en la ruta «" + path + "» no existente")
		return
	}

	particion := utils.PartitionVacia()
	mbr, embr := utils.Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return
	}

	pos := -1
	for i := range mbr.Mbr_partitions {
		if mbr.Mbr_partitions[i].Part_start == -1 {
			pos = i
			break
		}
	}

	if utils.ExisteNombreP(path, utils.ToString(_name)) {
		color.Red("[FDISK]: La particion «" + utils.ToString(_name) + "» ya existe")
		return
	}

	startByte := utils.EncontrarAjuste(mbr, utils.Tamano(_size, _unit), _fit)

	if startByte != -1 && pos != -1 {
		if !utils.ExisteNombreP(path, utils.ToString(_name)) {
			//existe extendida
			if !utils.ExisteParticionExt(path) {

				particion.Part_fit = _fit
				particion.Part_type = _type
				particion.Part_name = utils.DevolverNombreByte(utils.ToString(_name))
				particion.Part_status = '1'
				particion.Part_correlative = int32(pos + 1)
				particion.Part_s = utils.Tamano(_size, _unit)
				particion.Part_start = startByte

				posEBR := particion.Part_start
				mbr.Mbr_partitions[pos] = particion
				file, err := os.OpenFile(path, os.O_RDWR, 0666)
				if err != nil {
					color.Red("[FDISK]: Error al abrir archivo")
					return
				}
				defer file.Close()
				if _, err := file.Seek(0, 0); err != nil {
					color.Red("[FDISK]: Error en mover puntero")
					return
				}
				if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
					color.Red("[FDISK]: Error en la escritura del MBR")
					return
				}

				// Crear el primer EBR al inicio de la partición extendida
				ebr := structures.EBR{}
				//ebr.Part_status = '0' // Inactivo
				ebr.Part_fit = 'W'
				ebr.Part_start = posEBR
				ebr.Part_s = 0 // No ocupa espacio lógico aún
				ebr.Part_next = -1
				copy(ebr.Name[:], "")

				if _, err := file.Seek(int64(posEBR), 0); err != nil {
					color.Red("[FDISK]: Error en mover puntero al inicio del EBR")
					return
				}
				if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
					color.Red("[FDISK]: Error en la escritura del EBR inicial")
					return
				}
				file.Close()
				//comprobación
				comprobacion := structures.MBR{}
				file, err = os.OpenFile(path, os.O_RDWR, 0666)
				if err != nil {
					color.Red("[FDISK]: Error al abrir archivo")
					return
				}

				defer file.Close()
				if _, err := file.Seek(0, 0); err != nil {
					color.Red("[FDISK]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &comprobacion); err != nil {
					color.Red("[FDISK]: Error en la lectura del MBR")
					return
				}
				file.Close()
				color.Green("-----------------------------------------------------------")
				color.Blue("Se creo la particion #" + strconv.Itoa(int(comprobacion.Mbr_partitions[pos].Part_correlative)))
				color.Blue("Particion: " + utils.ToString(comprobacion.Mbr_partitions[pos].Part_name[:]))
				color.Blue("Tipo Extendida")
				color.Blue("Inicio: " + strconv.Itoa(int(comprobacion.Mbr_partitions[pos].Part_start)))
				color.Blue("Tamaño: " + strconv.Itoa(int(comprobacion.Mbr_partitions[pos].Part_s)))
				color.Green("-----------------------------------------------------------")
			} else {
				color.Yellow("[FDISK]: Particion extendida ya existente, no puede existir otra -> " + utils.ToString(_name))
				return
			}
		} else {
			color.Yellow("[FDISK]: Particion <" + utils.ToString(_name) + "> existente")
			return
		}
	} else {
		if pos == -1 {
			color.Red("[FDISK]: No se pueden crear más de 4 particiones primarias o extendidas.")
		} else {
			color.Red("[FDISK]: Espacio Insuficiente para la partición.")
		}
		return
	}
}
