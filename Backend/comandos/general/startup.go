package general

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
)

// LoadMountedPartitions escanea los discos en busca de particiones previamente montadas
// y las carga en la lista global `Mounted_Partitions`.
func LoadMountedPartitions() {
	color.Cyan("Cargando particiones montadas previamente...")
	diskPath := "VDIC-MIA/Disks/"
	files, err := os.ReadDir(diskPath)
	if err != nil {
		color.Red("Error al leer el directorio de discos: %v", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".mia" {
			fullPath := diskPath + file.Name()
			mbr, success := utils.Obtener_FULL_MBR_FDISK(fullPath)
			if !success {
				continue
			}

			for _, part := range mbr.Mbr_partitions {
				if part.Part_status == '1' && part.Part_type == 'P' {
					// Se encontró una partición primaria montada.
					var driveLetter byte
					re := regexp.MustCompile(`VDIC-([A-Z])\.mia`)
					match := re.FindStringSubmatch(file.Name())
					if len(match) > 1 {
						driveLetter = match[1][0]
					} else {
						continue // No se pudo determinar la letra, se ignora.
					}

					mount_temp := global.ParticionesMontadas{
						Path:           fullPath,
						DriveLetter:    driveLetter,
						ID_Particion:   part.Part_id,
						Type:           part.Part_type,
						Es_Particion_P: true,
						Es_Particion_L: false,
						Particion_P:    part,
						Particion_L:    structures.EBR{},
					}
					global.Mounted_Partitions = append(global.Mounted_Partitions, mount_temp)
					idStr := utils.ToString(part.Part_id[:])
					nameStr := utils.ToString(part.Part_name[:])
					fmt.Printf("  -> Partición cargada: id=%s, nombre=%s, disco=%s\n", idStr, nameStr, file.Name())
				}
			}
		}
	}
	color.Green("Carga de particiones finalizada.")
}
