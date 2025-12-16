package utils

import (
	"Proyecto/comandos/global"
)

func Buscar_ID_Montada(id_disco string) (global.ParticionesMontadas, bool) {
	for _, disco := range global.Mounted_Partitions {
		if disco.ID_Particion == IDParticionByte(id_disco) {
			//global.Mounted_Partitions = append(global.Mounted_Partitions, disco)
			return disco, true
		}
	}
	return global.ParticionesMontadas{}, false
}

func NameCarpeta12(nombre string) [12]byte {
	temp := make([]byte, 12)
	for i := range temp {
		temp[i] = '\x00'
	}
	copy(temp[:], []byte(nombre))
	return [12]byte(temp)
}
