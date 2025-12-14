package utils

import (
	"Proyecto/Estructuras/structures"
	"math/rand"
	"time"
)

// ObDiskSignature genera una firma de disco aleatoria.
func ObDiskSignature() int32 {
	// Crea una nueva fuente de números pseudoaleatorios utilizando la hora actual en nanosegundos como semilla.
	// Esto asegura que la secuencia de números sea diferente en cada ejecución.
	source := rand.NewSource(time.Now().UnixNano())
	// Crea un nuevo generador de números aleatorios a partir de la fuente creada anteriormente.
	numberR := rand.New(source)
	// Genera un número entero aleatorio entre 1 y 1,000,000.
	signature := numberR.Intn(1000000) + 1
	// Convierte el número generado a un entero de 32 bits (int32) y lo devuelve.
	return int32(signature)
}

// PartitionVacia inicializa y devuelve una estructura de Partición con valores predeterminados
// que indican que está vacía o sin usar.
func PartitionVacia() structures.Partition {
	// Declara una variable partition del tipo structures.Partition.
	var partition structures.Partition
	// Establece el estado de la partición a -1, probablemente para indicar que está inactiva.
	partition.Part_status = int8(-1)
	// Establece el tipo de partición a 'P' (posiblemente "Primaria").
	partition.Part_type = 'P'
	// Establece el tipo de ajuste de la partición a 'F' (posiblemente "First Fit" o Primer Ajuste).
	partition.Part_fit = 'F'
	// Establece el inicio de la partición en -1, indicando que aún no está asignada en el disco.
	partition.Part_start = -1
	// Establece el tamaño de la partición en -1, indicando que no tiene un tamaño definido.
	partition.Part_s = -1
	// Limpia el campo del nombre de la partición, llenándolo con caracteres nulos.
	for i := 0; i < len(partition.Part_name); i++ {
		partition.Part_name[i] = '\x00'
	}
	// Establece un número correlativo a -1.
	partition.Part_correlative = -1
	// Limpia el campo del identificador de la partición, llenándolo con caracteres nulos.
	for i := 0; i < len(partition.Part_id); i++ {
		partition.Part_id[i] = '\x00'
	}
	// Devuelve la estructura de partición inicializada.
	return partition
}
