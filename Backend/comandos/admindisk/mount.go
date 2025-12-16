package adminDisk

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Values_Mount(instructions []string) (byte, [16]byte, bool) {
	var _driveletter byte
	var _name [16]byte
	var _error = false
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "driveletter") {
			var value = utils.TieneDriveLetter("MOUNT", valor)
			_driveletter = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = utils.TieneNombre("MOUNT", valor)
			if len(value) > 16 {
				color.Red("[MOUNT]: El nombre no puede ser mayor a 16 caracteres")
				_error = true
				break
			} else {
				_name = utils.DevolverNombreByte(value)
			}
		} else {
			color.Yellow("[MOUNT]: Atributo no reconocido")
		}
	}
	return _driveletter, _name, _error
}

func MOUNT_EXECUTE(_driveletter byte, _name []byte) {
	path := "VDIC-MIA/Disks/" + string(_driveletter) + ".mia"
	if !utils.ExisteArchivo("MOUNT", path) {
		color.Yellow("[MOUNT]: No existe el disco")
		return
	}

	mbr, embr := utils.Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return
	}

	conjunto, eco := utils.BuscarParticion(mbr, _name, path)
	if !eco {
		return
	}

	//verificar que exista la particion en la lista
	numero := int32(1)
	contador := int32(1)
	for _, disco := range global.Mounted_Partitions {
		if disco.DriveLetter == _driveletter {
			if disco.Es_Particion_L {
				if utils.ToString(disco.Particion_L.Name[:]) == utils.ToString(_name) {
					color.Red("[Mount]: Particion (logica) ya montada -> " + utils.ToString(_name))
					return
				}
			} else if disco.Es_Particion_P {
				if utils.ToString(disco.Particion_P.Part_name[:]) == utils.ToString(_name) {
					color.Red("[Mount]: Particion (primaria) ya montada -> " + utils.ToString(_name))
					return
				}
			}
			contador++
		}
	}
	numero = contador

	mount_temp := global.ParticionesMontadas{}
	nombre_part := string(_driveletter) + strconv.Itoa(int(numero)) + "51"
	nombre_bytes := utils.IDParticionByte(nombre_part)
	mount_temp.DriveLetter = _driveletter
	mount_temp.ID_Particion = nombre_bytes
	mount_temp.Path = path

	color.Magenta("[Mount]: Montando particion... " + nombre_part + ", (p): -> " + utils.ToString(_name))

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Mount]: Error al abrir archivo")
		return
	}

	inicio := int32(0)

	// Caso de ser particion primaria
	particion := structures.Partition{}
	if temp, ok := conjunto[0].(structures.Partition); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&particion).Elem().Set(v)
		if particion.Part_type == 'E' {
			color.Red("[Mount]: No se puede montar una particion extendida -> " + utils.ToString(_name))
			return
		}

		// Inicio SB
		inicio = particion.Part_start
		// Verificacion
		mount_temp.Es_Particion_P = true
		mount_temp.Es_Particion_L = false
		particion.Part_id = nombre_bytes
		mount_temp.Type = 'P'
		particion.Part_status = 0
		mount_temp.Particion_P = particion
		count := 0
		for _, c := range mbr.Mbr_partitions {
			if utils.ToString(c.Part_name[:]) == utils.ToString(_name) {
				break
			}
			count++
		}
		mbr.Mbr_partitions[count] = particion

		if _, err := file.Seek(0, 0); err != nil {
			color.Red("[Mount]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
			color.Red("[Mount]: Error en la escritura del MBR")
			return
		}
	}

	// Caso de ser particion logica
	// logica := structures.EBR{}
	// if temp, ok := conjunto[0].(structures.EBR); ok {
	// 	v := reflect.ValueOf(temp)
	// 	reflect.ValueOf(&logica).Elem().Set(v)

	// 	// Inicio SB
	// 	inicio = logica.Part_start + size.SizeEBR()
	// 	// Verificacion
	// 	mount_temp.Es_Particion_P = false
	// 	mount_temp.Es_Particion_L = true
	// 	//id no existe en ebr
	// 	mount_temp.Type = 'L'
	// 	logica.Part_mount = 0
	// 	mount_temp.Particion_L = logica

	// 	if _, err := file.Seek(int64(logica.Part_start), 0); err != nil {
	// 		color.Red("[Mount]: Error en mover puntero")
	// 		return
	// 	}
	// 	if err := binary.Write(file, binary.LittleEndian, &logica); err != nil {
	// 		color.Red("[Mount]: Error en la escritura del EBR")
	// 		return
	// 	}
	// }

	global.Mounted_Partitions = append(global.Mounted_Partitions, mount_temp)
	superblock := structures.SuperBloque{}
	if _, err := file.Seek(int64(inicio), 0); err != nil {
		color.Red("[Mount]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &superblock); err != nil {
		color.Red("[Mount]: Error en la lectura del superbloque")
		return
	}
	if superblock.S_mnt_count > 0 {
		superblock.S_mtime = utils.ObFechaInt()
		superblock.S_mnt_count += 1
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[Mount]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &superblock); err != nil {
			color.Red("[Mount]: Error en la escritura del SuperBloque")
			return
		}
	}

	color.Green("[Mount]: Particion «" + utils.ToString(_name) + "» montada - (id) -> [" + nombre_part + "]")
	file.Close()
	file = nil
}
