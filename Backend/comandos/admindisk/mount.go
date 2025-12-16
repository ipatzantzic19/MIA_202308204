package adminDisk

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Values_Mount(instructions []string) (string, [16]byte, bool) {
	var _diskName string
	var _name [16]byte
	var _error = false
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "diskname") {
			// Llama a la función auxiliar para validar y extraer el nombre del disco
			_diskName = utils.TieneDiskName("MOUNT", valor)
			if _diskName == "" {
				color.Red("[MOUNT]: Error en el parametro DiskName")
				_error = true
				break
			}
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
			_error = true
			break
		}
	}
	return _diskName, _name, _error
}

func MOUNT_EXECUTE(diskName string, _name []byte) {
	path := "VDIC-MIA/Disks/" + diskName
	if !utils.ExisteArchivo("MOUNT", path) {
		color.Yellow("[MOUNT]: No existe el disco")
		return
	}

	//Obtenemos los ultimos dos digitos del carnet
	carnet := "202308204"
	re := regexp.MustCompile(`\d{2}$`)
	match := re.FindStringSubmatch(carnet)
	hexadecimal := ""
	if len(match) > 0 {
		decimal, err := strconv.Atoi(match[0])
		if err != nil {
			panic(err)
		}
		hexadecimal = strconv.FormatInt(int64(decimal), 16)
		color.Cyan("[MOUNT]: Carnet en hexadecimal: " + hexadecimal)

	} else {
		color.Red("[MOUNT]: Error al obtener los ultimos dos digitos del carnet")
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
	contador := 1
	var driveLetter byte
	re = regexp.MustCompile(`VDIC-([A-Z])\.mia`)
	match = re.FindStringSubmatch(diskName)
	if len(match) > 1 {
		driveLetter = match[1][0]
	} else {
		color.Red("[MOUNT]: Error al obtener la letra del disco")
		return
	}

	for _, disco := range global.Mounted_Partitions {
		if disco.Path == path {
			if disco.Es_Particion_L {
				if utils.ToString(disco.Particion_L.Name[:]) == utils.ToString(_name) {
					color.Red("[Mount]: Particion (logica) ya montada -> " + utils.ToString(_name))
					return
				}
			} else {
				if utils.ToString(disco.Particion_P.Part_name[:]) == utils.ToString(_name) {
					color.Red("[Mount]: Particion (primaria) ya montada -> " + utils.ToString(_name))
					return
				}
			}
			contador++
		}
	}

	mount_temp := global.ParticionesMontadas{}
	nombre_part := hexadecimal + strconv.Itoa(contador) + string(driveLetter)
	nombre_bytes := utils.IDParticionByte(nombre_part)
	mount_temp.DriveLetter = driveLetter
	mount_temp.ID_Particion = nombre_bytes
	mount_temp.Path = path

	color.Magenta("[Mount]: Montando particion... " + nombre_part + ", (p): -> " + utils.ToString(_name))

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Mount]: Error al abrir archivo")
		return
	}
	defer file.Close()

	inicio := int32(0)

	particion := structures.Partition{}
	if temp, ok := conjunto[0].(structures.Partition); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&particion).Elem().Set(v)
		if particion.Part_type == 'E' {
			color.Red("[Mount]: No se puede montar una particion extendida -> " + utils.ToString(_name))
			return
		}

		if particion.Part_type != 'P' {
			color.Red("[Mount]: Solo se pueden montar particiones primarias -> " + utils.ToString(_name))
			return
		}

		inicio = particion.Part_start
		mount_temp.Es_Particion_P = true
		mount_temp.Es_Particion_L = false
		particion.Part_id = nombre_bytes
		mount_temp.Type = 'P'
		particion.Part_status = '1' // 1 = mounted
		mount_temp.Particion_P = particion
		count := 0
		for i, p := range mbr.Mbr_partitions {
			if utils.ToString(p.Part_name[:]) == utils.ToString(_name) {
				mbr.Mbr_partitions[i].Part_id = nombre_bytes
				mbr.Mbr_partitions[i].Part_status = '1'
				break
			}
			count++
		}

		if _, err := file.Seek(0, 0); err != nil {
			color.Red("[Mount]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
			color.Red("[Mount]: Error en la escritura del MBR")
			return
		}
	} else {
		color.Red("[Mount]: La particion encontrada no es primaria.")
		return
	}

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
}
