package size

import (
	// "proyecto/estructuras/structures"
	"Proyecto/Estructuras/structures"
	"unsafe"
)

// / FUNCIONES PARA OBTENER EL TAMAÑO DE LAS ESTRUCTURAS EN BYTES
func SizeContent() int32 { //16bytes
	a01 := unsafe.Sizeof(structures.Content{}.B_name)
	a01 += unsafe.Sizeof(structures.Content{}.B_inodo)
	// result := a01
	return int32(a01)
}

func SizeBloqueCarpeta() int32 { //64 bytes
	a01 := unsafe.Sizeof(structures.BloqueCarpeta{})
	return int32(a01)
}

func SizeBloqueArchivo() int32 { //64 bytes
	a01 := unsafe.Sizeof(structures.BloqueArchivo{}.B_content)
	return int32(a01)
}

func SizeBloqueApuntador() int32 { //64 bytes
	a01 := unsafe.Sizeof(structures.BloqueApuntador{}.B_pointers)
	return int32(a01)
}
