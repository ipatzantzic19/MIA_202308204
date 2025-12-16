package size

import (
	// "proyecto/estructuras/structures"
	"Proyecto/Estructuras/structures"
	"unsafe"
)

// FUNCIONES PARA OBTENER EL TAMAÑO DE LAS ESTRUCTURAS EN BYTES
func SizeSuperBloque() int32 { //68 bytes
	a01 := unsafe.Sizeof(structures.SuperBloque{})
	return int32(a01)
}

func SizeTablaInodo() int32 { //92 bytes
	a01 := unsafe.Sizeof(structures.TablaInodo{})
	return int32(a01)
}
