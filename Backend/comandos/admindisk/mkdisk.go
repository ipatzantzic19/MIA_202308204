package adminDisk

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

/*
MKDISK
------
Este comando crea un archivo binario que simula un disco duro.
- El archivo tendrá extensión .mia
- El archivo se llena inicialmente con bytes en 0
- El tamaño físico del archivo dependerá del parámetro -size y -unit
- Se escribe una estructura MBR al inicio del archivo

Parámetros:
- size (OBLIGATORIO): tamaño del disco (entero positivo > 0)
- fit  (OPCIONAL): BF | FF | WF (por defecto FF)
- unit (OPCIONAL): K | M (por defecto M)
*/

// Values_MKDISK analiza y valida los parámetros del comando MKDISK.
func Values_MKDISK(instructions []string) (int32, byte, byte, error) {
	var size int32
	var fit byte = 'F'  // FF = First Fit (por defecto)
	var unit byte = 'M' // Megabytes por defecto
	var sizeFound bool

	// Recorremos cada parámetro recibido
	for _, valor := range instructions {
		param := strings.ToLower(strings.TrimSpace(valor))

		// ---------- SIZE (obligatorio) ----------
		if strings.HasPrefix(param, "size") {
			size = utils.TieneSize("MKDISK", valor)
			sizeFound = true

			// ---------- FIT (opcional) ----------
		} else if strings.HasPrefix(param, "fit") {
			fit = utils.TieneFit("MKDISK", valor)
			// Validación explícita
			if fit != 'B' && fit != 'F' && fit != 'W' {
				return -1, '0', '0', fmt.Errorf("[MKDISK]: fit inválido (BF | FF | WF)")
			}

			// ---------- UNIT (opcional) ----------
		} else if strings.HasPrefix(param, "unit") {
			unit = utils.TieneUnit("MKDISK", valor)
			if unit != 'K' && unit != 'M' {
				return -1, '0', '0', fmt.Errorf("[MKDISK]: unit inválido (K | M)")
			}

			// ---------- PARÁMETRO DESCONOCIDO ----------
		} else {
			return -1, '0', '0', fmt.Errorf("[MKDISK]: atributo no reconocido: %s", valor)
		}
	}

	// Validación final del size
	if !sizeFound || size <= 0 {
		return -1, '0', '0', fmt.Errorf("[MKDISK]: el parámetro -size es obligatorio y debe ser mayor que 0")
	}

	return size, fit, unit, nil
}

// MKDISK_Create crea el disco buscando una letra disponible de A a Z.
func MKDISK_Create(size int32, fit byte, unit byte) (string, error) {
	directorio := "VDIC-MIA/Disks/"

	// Recorre letras A-Z
	for i := 0; i < 26; i++ {
		nombre := fmt.Sprintf("VDIC-%c.mia", 'A'+i)
		ruta := directorio + nombre

		// Si el archivo NO existe, se crea
		if _, err := os.Stat(ruta); os.IsNotExist(err) {
			if err := CreateFile(ruta, size, fit, unit); err != nil {
				return "", err
			}
			return fmt.Sprintf("[MKDISK]: Disco %s creado correctamente", nombre), nil
		}
	}

	return "", fmt.Errorf("[MKDISK]: no hay letras disponibles para crear discos")
}

// CreateFile crea físicamente el archivo del disco, lo llena con ceros
// y escribe el MBR al inicio.
func CreateFile(path string, size int32, fit byte, unit byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("[MKDISK]: no se pudo crear el archivo")
	}
	defer file.Close()

	// Conversión del tamaño a bytes
	tamanioBytes := utils.Tamano(size, unit)

	// ----------- CREACIÓN DEL MBR -----------
	var mbr structures.MBR
	mbr.Mbr_tamano = tamanioBytes
	mbr.Mbr_fecha_creacion = utils.ObFechaInt()
	mbr.Mbr_disk_signature = utils.ObDiskSignature()
	mbr.Dsk_fit = fit

	// Inicializar particiones vacías
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		mbr.Mbr_partitions[i] = utils.PartitionVacia()
	}

	// ----------- LLENAR ARCHIVO CON CEROS -----------
	bloque := make([]byte, 1024)
	bytesRestantes := tamanioBytes

	for bytesRestantes > 0 {
		if bytesRestantes < int32(len(bloque)) {
			bloque = make([]byte, bytesRestantes)
		}
		file.Write(bloque)
		bytesRestantes -= int32(len(bloque))
	}

	// ----------- ESCRIBIR MBR -----------
	file.Seek(0, 0)
	if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
		return fmt.Errorf("[MKDISK]: error al escribir el MBR")
	}

	return nil
}
