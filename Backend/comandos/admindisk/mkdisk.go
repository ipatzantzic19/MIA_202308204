package adminDisk

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// Values_MKDISK analiza las instrucciones (parámetros) para el comando MKDISK.
// Valida y extrae los valores de -size, -fit y -unit.
// Devuelve los valores extraídos o un error si los parámetros son inválidos.
func Values_MKDISK(instructions []string) (int32, byte, byte, error) {
	var _size int32
	var _fit byte = 'F'  // Valor por defecto para el ajuste: 'F' (First Fit)
	var _unit byte = 'M' // Valor por defecto para la unidad: 'M' (Megabytes)
	var sizeFound bool   // Flag para asegurar que el parámetro -size esté presente.

	// Itera sobre cada instrucción (parámetro) para identificar y extraer su valor.
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			// Extrae el valor de -size.
			var value = utils.TieneSize("MKDISK", valor)
			_size = value
			sizeFound = true
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			// Extrae el valor de -fit.
			var value = utils.TieneFit("MKDISK", valor)
			_fit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			// Extrae el valor de -unit.
			var value = utils.TieneUnit("mkdisk", valor)
			_unit = value
		} else {
			// Si se encuentra un atributo no reconocido, devuelve un error.
			return -1, '0', '0', fmt.Errorf("[MKDISK]: Atributo no reconocido: %s", valor)
		}
	}

	// Valida que el parámetro -size haya sido proporcionado y que sea un valor positivo.
	if !sizeFound || _size <= 0 {
		return -1, '0', '0', fmt.Errorf("[MKDISK]: El atributo -size es obligatorio y debe ser un entero positivo")
	}
	// Devuelve los valores validados.
	return _size, _fit, _unit, nil
}

// MKDISK_Create se encarga de la creación del archivo de disco.
// Busca un nombre de archivo disponible (de A.dsk a Z.dsk) y crea el disco.
// Devuelve un mensaje de éxito o un error si no se puede crear el disco.
func MKDISK_Create(_size int32, _fit byte, _unit byte) (string, error) {
	directorio := "VDIC-MIA/Disks/"
	// Asegura que el directorio base para los discos exista, si no, lo crea.
	if err := os.MkdirAll(directorio, 0755); err != nil {
		return "", fmt.Errorf("error al crear directorio '%s': %w", directorio, err)
	}

	// Itera de la 'A' a la 'Z' para encontrar un nombre de disco disponible.
	for i := 0; i < 26; i++ {
		letra := fmt.Sprintf("%c.dsk", 'A'+i)
		archivo := directorio + letra
		// Comprueba si el archivo ya existe.
		if _, err := os.Stat(archivo); os.IsNotExist(err) {
			// Si el archivo no existe, procede a crearlo.
			err := CreateFile(archivo, _size, _fit, _unit)
			if err != nil {
				// Si hay un error durante la creación, devuelve un mensaje de error detallado.
				return "", fmt.Errorf("error al crear el disco '%s': %w", letra, err)
			}
			// Si la creación es exitosa, retorna el mensaje de éxito.
			return fmt.Sprintf("[MKDISK]: Disco '%s' Creado -> %d%c", letra, _size, _unit), nil
		}
	}
	// Si el bucle termina, significa que no se encontraron letras disponibles.
	return "", fmt.Errorf("[MKDISK]: No hay letras de disco disponibles (A-Z)")
}

// CreateFile crea el archivo físico del disco, lo inicializa con ceros y escribe el MBR.
// Devuelve un error si alguna de estas operaciones falla.
func CreateFile(archivo string, _size int32, _fit byte, _unit byte) error {
	// Crea el archivo en la ruta especificada.
	file, err := os.Create(archivo)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo: %w", err)
	}
	defer file.Close() // Asegura que el archivo se cierre al final de la función.

	var estructura structures.MBR         // Crea una instancia de la estructura MBR.
	tamanio := utils.Tamano(_size, _unit) // Calcula el tamaño del disco en bytes.
	estructura.Mbr_tamano = tamanio
	estructura.Mbr_fecha_creacion = utils.ObFechaInt()      // Asigna la fecha de creación.
	estructura.Mbr_disk_signature = utils.ObDiskSignature() // Asigna una firma única al disco.
	estructura.Dsk_fit = _fit                               // Asigna el tipo de ajuste.
	// Inicializa las particiones del MBR como vacías.
	for i := 0; i < len(estructura.Mbr_partitions); i++ {
		estructura.Mbr_partitions[i] = utils.PartitionVacia()
	}

	// Llena el archivo con ceros para reservar el espacio.
	bytes_llenar := make([]byte, int(tamanio))
	if _, err := file.Write(bytes_llenar); err != nil {
		return fmt.Errorf("no se pudo llenar el archivo con ceros: %w", err)
	}

	// Regresa al inicio del archivo para escribir el MBR.
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("no se pudo posicionar el puntero al inicio del archivo: %w", err)
	}

	// Escribe la estructura MBR en formato binario (Little Endian) al inicio del archivo.
	if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
		return fmt.Errorf("no se pudo escribir la estructura MBR: %w", err)
	}
	return nil // Retorna nil si todas las operaciones fueron exitosas.
}
