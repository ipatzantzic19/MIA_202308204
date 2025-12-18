package adminusers

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Values_MKGRP analiza y valida los parámetros del comando MKGRP
func Values_MKGRP(instructions []string) (string, bool) {
	var _name string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "name") {
			value := utils.TieneNombre("MKGRP", valor)
			if len(value) > 10 {
				color.Red("[MKGRP]: El nombre no puede ser mayor a 10 caracteres")
				return "", false
			}
			_name = value
		} else {
			color.Yellow("[MKGRP]: Atributo no reconocido: %s", valor)
		}
	}

	if _name == "" || len(_name) == 0 {
		color.Red("[MKGRP]: El parámetro -name es obligatorio")
		return "", false
	}

	return _name, true
}

// MKGRP_EXECUTE ejecuta la creación del grupo
func MKGRP_EXECUTE(_name string) {
	// 1. Verificar que hay una sesión activa
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[MKGRP]: No hay usuario loggeado")
		return
	}

	// 2. Verificar que el usuario es root
	if global.UsuarioLogeado.UID != 1 || global.UsuarioLogeado.GID != 1 {
		color.Red("[MKGRP]: No tienes permisos para ejecutar este comando. Solo el usuario root puede crear grupos.")
		return
	}

	mount := global.UsuarioLogeado.Mounted

	// 3. Abrir el archivo del disco
	file, err := os.OpenFile(mount.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[MKGRP]: Error al abrir archivo")
		return
	}
	defer file.Close()

	// 4. Leer el superbloque según el tipo de partición
	if mount.Es_Particion_L {
		inicio := mount.Particion_L.Part_start + size.SizeEBR()
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[MKGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[MKGRP]: Error en la lectura del SuperBloque")
			return
		}
	} else if mount.Es_Particion_P {
		inicio := mount.Particion_P.Part_start
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[MKGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[MKGRP]: Error en la lectura del SuperBloque")
			return
		}
	}

	// 5. Leer el contenido actual del archivo users.txt (inodo 1)
	content := utils.GetContentAdminUsers(utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo())
	cantBlockAnt := int32(len(utils.SplitContent(content)))

	// 6. Procesar el contenido y verificar que el grupo no existe
	split_content := strings.Split(content, "\n")
	if utils.GrupoExist(split_content, _name) {
		color.Red("[MKGRP]: El grupo «%s» ya existe", _name)
		return
	}

	// 7. Crear la nueva entrada de grupo
	nuevoGrp := fmt.Sprint(utils.GetGID(split_content)) + ",G," + _name + "\n"
	content += nuevoGrp

	// 8. Dividir el contenido en bloques de 64 caracteres
	usersTxt := utils.SplitContent(content)
	cantBlockAct := int32(len(usersTxt))

	// 9. Verificar límite de bloques
	if len(usersTxt) > 4380 {
		color.Red("[MKGRP]: No se pueden crear más grupos (límite de bloques alcanzado)")
		return
	}

	// 10. Verificar si hay suficientes bloques libres
	if utils.Sb_AdminUsr.S_free_blocks_count < (cantBlockAct - cantBlockAnt) {
		color.Red("[MKGRP]: No hay bloques suficientes para actualizar el archivo users.txt")
		return
	}

	// 11. Buscar bloques contiguos disponibles en el bitmap (solo si se necesitan nuevos bloques)
	var inicioBM int32 = -1
	var inicioB int32 = -1

	if (cantBlockAct - cantBlockAnt) > 0 {
		var bit byte
		start := utils.Sb_AdminUsr.S_bm_block_start
		end := start + utils.Sb_AdminUsr.S_blocks_count
		cantContiguos := int32(0)
		contadorA := int32(0)

		for i := start; i < end; i++ {
			if _, err := file.Seek(int64(i), 0); err != nil {
				color.Red("[MKGRP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
				color.Red("[MKGRP]: Error en la lectura del bitmap de bloques")
				return
			}

			if bit == '1' { // Ocupado
				cantContiguos = 0
				inicioBM = -1
				inicioB = -1
			} else { // Libre
				if cantContiguos == 0 {
					inicioBM = int32(i)
					inicioB = contadorA
				}
				cantContiguos++
			}

			if cantContiguos >= (cantBlockAct - cantBlockAnt) {
				break
			}
			contadorA++
		}

		// Validar que se encontraron suficientes bloques contiguos
		if (inicioBM == -1) || (cantContiguos != (cantBlockAct - cantBlockAnt)) {
			color.Red("[MKGRP]: No hay bloques contiguos suficientes para actualizar el archivo users.txt")
			return
		}

		// 12. Marcar los bloques como ocupados en el bitmap
		var uno byte = '1'
		for i := inicioBM; i < (inicioBM + (cantBlockAct - cantBlockAnt)); i++ {
			if _, err := file.Seek(int64(i), 0); err != nil {
				color.Red("[MKGRP]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &uno); err != nil {
				color.Red("[MKGRP]: Error en la escritura del bitmap de bloques")
				return
			}
		}

		// 13. Actualizar contadores del superbloque
		utils.Sb_AdminUsr.S_free_blocks_count -= (cantBlockAct - cantBlockAnt)

		// Buscar el primer bloque libre
		bit2 := int32(0)
		for k := start; k < end; k++ {
			if _, err := file.Seek(int64(k), 0); err != nil {
				color.Red("[MKGRP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
				color.Red("[MKGRP]: Error en la lectura del bitmap de bloques")
				return
			}
			if bit == '0' {
				break
			}
			bit2++
		}
		utils.Sb_AdminUsr.S_first_blo = bit2
	}

	// 14. Leer el inodo del archivo users.txt (inodo 1)
	inodo := structures.TablaInodo{}
	seekInodo := utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo()
	if _, err := file.Seek(int64(seekInodo), 0); err != nil {
		color.Red("[MKGRP]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[MKGRP]: Error en la lectura del Inodo")
		return
	}

	// 15. Calcular el nuevo tamaño del archivo
	tamanio := int32(0)
	for tm := range usersTxt {
		tamanio += int32(len(usersTxt[tm]))
	}

	// 16. Actualizar el inodo
	inodo.I_s = tamanio
	inodo.I_mtime = utils.ObFechaInt()

	// 17. Escribir el contenido en los bloques usando la función auxiliar
	var j, contador = 0, 0
	for j < len(usersTxt) {
		utils.CambioCont = false
		inodo = utils.AgregarArchivo(usersTxt[j], inodo, int32(j), (inicioB + int32(contador)))
		if utils.CambioCont {
			contador++
		}
		j++
	}

	// 18. Escribir el inodo actualizado
	if _, err := file.Seek(int64(seekInodo), 0); err != nil {
		color.Red("[MKGRP]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[MKGRP]: Error en la escritura del Inodo")
		return
	}

	// 19. Escribir el superbloque actualizado
	if mount.Es_Particion_P {
		if _, err := file.Seek(int64(mount.Particion_P.Part_start), 0); err != nil {
			color.Red("[MKGRP]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[MKGRP]: Error en la escritura del SuperBloque")
			return
		}
	} else if mount.Es_Particion_L {
		if _, err := file.Seek(int64(mount.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
			color.Red("[MKGRP]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[MKGRP]: Error en la escritura del SuperBloque")
			return
		}
	}

	color.Green("[MKGRP]: Grupo «%s» creado correctamente", _name)
}
