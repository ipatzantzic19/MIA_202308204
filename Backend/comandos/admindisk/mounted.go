package adminDisk

import (
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"fmt"

	"github.com/fatih/color"
)

// Esta funcion muestra todas las particiones montadas actualmente.
func MOUNTED_EXECUTE() {
	if len(global.Mounted_Partitions) == 0 {
		color.Yellow("[MOUNTED]: No hay particiones montadas.")
		return
	}

	color.Green("Particiones Montadas:")
	for _, particion := range global.Mounted_Partitions {
		id := utils.ToString(particion.ID_Particion[:])
		path := particion.Path
		var name string
		if particion.Es_Particion_P {
			name = utils.ToString(particion.Particion_P.Part_name[:])
		} else {
			name = utils.ToString(particion.Particion_L.Name[:])
		}

		fmt.Printf("- id: %s, path: %s, name: %s\n", id, path, name)
	}
}
