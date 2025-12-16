package size

import (
	// "proyecto/estructuras/structures"
	"Proyecto/Estructuras/structures"
	"unsafe"
)

// FUNCIONES PARA OBTENER EL TAMAÑO DE LAS ESTRUCTURAS EN BYTES
func SizeEBR() int32 { //30 bytes
	a01 := unsafe.Sizeof(structures.EBR{})
	return int32(a01)
}
func SizePartition() int32 { //35 bytes
	a01 := unsafe.Sizeof(structures.Partition{})
	return int32(a01)
}

func SizeMBR() int32 { //153 bytes
	a01 := unsafe.Sizeof(structures.MBR{})
	return int32(a01)
}

func SizeMBR_NotPartitions() int32 {
	a01 := unsafe.Sizeof(structures.MBR{}.Mbr_tamano)
	a01 += unsafe.Sizeof(structures.MBR{}.Mbr_fecha_creacion)
	a01 += unsafe.Sizeof(structures.MBR{}.Mbr_disk_signature)
	a01 += unsafe.Sizeof(structures.MBR{}.Dsk_fit)
	return int32(a01)
}
