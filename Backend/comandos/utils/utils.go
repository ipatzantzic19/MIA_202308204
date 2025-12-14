package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// TieneSize valida y extrae el valor del parámetro -size.
// Espera un formato como "size=1024". Devuelve el valor numérico o 0 si es inválido.
func TieneSize(comando string, size string) int32 {
	valsize := TieneEntero(size)
	if valsize <= 0 {
		color.Red("[" + comando + "]: No tiene Size o tiene un valor no valido")
		return 0
	}
	return valsize
}

// TieneFit valida y extrae el valor del parámetro -fit.
// Los valores válidos son 'B' (Best Fit), 'F' (First Fit), 'W' (Worst Fit).
// Devuelve el carácter correspondiente o '0' si es inválido.
func TieneFit(comando string, fit string) byte {
	if !strings.HasPrefix(strings.ToLower(fit), "fit=") {
		color.Red("[" + comando + "]: No tiene Fit o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(fit, "=")
	if len(value) < 2 {
		return 'F' // Valor por defecto si no se especifica.
	}
	val := strings.ToUpper(value[1])
	if val == "BF" || val == "B" {
		return 'B'
	} else if val == "FF" || val == "F" {
		return 'F'
	} else if val == "WF" || val == "W" {
		return 'W'
	} else {
		color.Yellow("[" + comando + "]: No tiene Fit Valido")
		return '0'
	}
}

// TieneUnit valida y extrae el valor del parámetro -unit.
// Los valores pueden ser 'B' (bytes), 'K' (kilobytes), 'M' (megabytes).
// Devuelve el carácter correspondiente o '0' si es inválido.
func TieneUnit(command string, unit string) byte {
	if !strings.HasPrefix(strings.ToLower(unit), "unit=") {
		color.Red("[" + command + "]: No tiene Unit o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(unit, "=")
	if len(value) < 2 {
		color.Red("[" + command + "]: No tiene Unit")
		return '0'
	}
	val := strings.ToUpper(value[1])
	if val == "B" {
		// Restricciones específicas del comando.
		if command == "MKDISK" {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'M'
		} else if command == "FDISK" {
			return 'B'
		} else {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'K'
		}
	} else if val == "K" {
		return 'K'
	} else if val == "M" {
		return 'M'
	} else {
		color.Red("[" + command + "]: No tiene Unit Valido")
		return '0'
	}
}

// TieneEntero convierte el valor de un parámetro como "size=123" a un entero.
// Devuelve 0 si el formato es incorrecto o la conversión falla.
func TieneEntero(valor string) int32 {
	if !strings.HasPrefix(strings.ToLower(valor), "size=") {
		return 0
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		return 0
	}
	i, err := strconv.Atoi(value[1])
	if err != nil {
		fmt.Println("Error conversion", err)
		return 0
	}
	return int32(i)
}

// ObFechaInt devuelve la fecha y hora actual como un timestamp Unix (entero de 32 bits).
func ObFechaInt() int32 {
	fecha := time.Now()
	timestamp := fecha.Unix()
	return int32(timestamp)
}

// IntFechaToStr convierte un timestamp Unix (entero de 32 bits) a una cadena con formato "YYYY/MM/DD (HH:MM:SS)".
func IntFechaToStr(fecha int32) string {
	conversion := int64(fecha)
	formato := "2006/01/02 (15:04:05)"
	fech := time.Unix(conversion, 0)
	fechaFormat := fech.Format(formato)
	return fechaFormat
}

// Tamano calcula el tamaño en bytes a partir de un valor y una unidad ('B', 'K', 'M').
func Tamano(size int32, unit byte) int32 {
	if unit == 'B' {
		return size
	} else if unit == 'K' {
		return size * 1024
	} else if unit == 'M' {
		return size * 1024 * 1024
	} else {
		return 0
	}
}

// Type_FDISK valida y extrae el tipo de partición para FDISK.
// Puede ser 'P' (Primaria), 'E' (Extendida), 'L' (Lógica).
func Type_FDISK(_type string) byte {
	if !strings.HasPrefix(strings.ToLower(_type), "type=") {
		return '0'
	}
	value := strings.Split(_type, "=")
	if len(value) < 2 {
		color.Red("[FDISK]: No tiene Type Especificado")
		return 'P' // Valor por defecto.
	}
	val := strings.ToUpper(value[1])
	if val == "P" {
		return 'P'
	} else if val == "E" {
		return 'E'
	} else if val == "L" {
		return 'L'
	} else {
		color.Red("[FDISK]: No reconocido Type")
		return '0'
	}
}

// Type_MKFS valida el tipo de formateo para MKFS.
// Actualmente solo acepta "FULL".
func Type_MKFS(_type string) string {
	if strings.ToUpper(_type) == "FULL" {
		return "FULL"
	} else {
		color.Red("[MKFS]: No reconocido comando Type")
		return ""
	}
}

// TieneDriveLetter extrae la letra de unidad de un parámetro como "driveletter=A".
func TieneDriveLetter(comando string, deletter string) byte {
	if !strings.HasPrefix(strings.ToLower(deletter), "driveletter=") {
		color.Red("[" + comando + "]: No tiene driveletter o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(deletter, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene deletter Valido")
		return '0'
	} else {
		valor := []byte(value[1])
		if len(valor) != 1 {
			color.Red("[" + comando + "]: No tiene driveletter Valido")
			return '0'
		} else {
			return valor[0]
		}
	}
}

// TieneNombre extrae el valor del parámetro -name.
func TieneNombre(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "name=") {
		color.Red("[" + comando + "]: No tiene name o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene name Valido")
		return ""
	} else {
		return value[1]
	}
}

// --- Funciones de conversión a arreglos de bytes con padding ---

// NameArchivosByte convierte un string a un arreglo de 10 bytes, rellenando con nulos.
func NameArchivosByte(value string) [10]byte {
	padText := make([]byte, 10)
	copy(padText, []byte(value))
	var result [10]byte
	copy(result[:], padText)
	return result
}

// ObJournalData convierte un string a un arreglo de 100 bytes para datos de journaling.
func ObJournalData(value string) [100]byte {
	padText := make([]byte, 100)
	copy(padText, []byte(value))
	var result [100]byte
	copy(result[:], padText)
	return result
}

// IDParticionByte convierte un string a un arreglo de 4 bytes para el ID de partición.
func IDParticionByte(value string) [4]byte {
	padText := make([]byte, 4)
	copy(padText, []byte(value))
	var result [4]byte
	copy(result[:], padText)
	return result
}

// DevolverNombreByte convierte un string a un arreglo de 16 bytes para nombres.
func DevolverNombreByte(value string) [16]byte {
	padText := make([]byte, 16)
	copy(padText, []byte(value))
	var result [16]byte
	copy(result[:], padText)
	return result
}

// DevolverContenidoJournal convierte un string a un arreglo de 150 bytes para contenido de journal.
func DevolverContenidoJournal(value string) [150]byte {
	padText := make([]byte, 150)
	copy(padText, []byte(value))
	var result [150]byte
	copy(result[:], padText)
	return result
}

// DevolverContenidoArchivo convierte un string a un arreglo de 64 bytes para contenido de archivo.
func DevolverContenidoArchivo(value string) [64]byte {
	padText := make([]byte, 64)
	copy(padText, []byte(value))
	var result [64]byte
	copy(result[:], padText)
	return result
}

// TieneTypeFDISK (duplicado de Type_FDISK) valida y extrae el tipo de partición.
func TieneTypeFDISK(valor string) byte {
	if !strings.HasPrefix(strings.ToLower(valor), "type=") {
		return '0'
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[FDISK]: No tiene Type Especificado")
		return '0'
	}
	val := strings.ToUpper(value[1])
	if val == "P" {
		return 'P'
	} else if val == "E" {
		return 'E'
	} else if val == "L" {
		return 'L'
	} else {
		color.Red("[FDISK]: No reconocido Type")
		return '0'
	}
}

// TieneDelete valida el parámetro -delete para FDISK. Solo acepta "FULL".
func TieneDelete(valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "delete=") {
		color.Red("[FDISK]: No tiene Delete Especificado")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[FDISK]: No tiene Delete Especificado")
		return ""
	}
	if !(strings.ToUpper(value[1]) == "FULL") {
		color.Red("[FDISK]: No tiene Delete valido")
		return ""
	}
	return "FULL"
}

// TieneAdd extrae y convierte el valor numérico del parámetro -add.
func TieneAdd(valor string) int32 {
	if !strings.HasPrefix(strings.ToLower(valor), "add=") {
		return 0
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		return 0
	}
	num, err := strconv.Atoi(value[1])
	if err != nil {
		color.Red("[FDISK]: valor Add no aceptado")
		return 0
	}
	return int32(num)
}

// TieneID extrae el valor del parámetro -id.
func TieneID(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "id=") {
		color.Red("[" + comando + "]: No tiene id o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene id Valido")
		return ""
	}
	return value[1]
}

// TieneTypeMKFS valida el tipo para MKFS, que solo puede ser "FULL".
func TieneTypeMKFS(valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "type=") {
		color.Red("[MKFS]: No tiene Type Especificado")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[MKFS]: No tiene Type Especificado")
		return ""
	}
	if !(strings.ToUpper(value[1]) == "FULL") {
		color.Red("[MKSF]: No tiene Type valido")
		return ""
	}
	return "FULL"
}

// TieneFS valida el sistema de archivos para MKFS, que puede ser "2FS" o "3FS".
func TieneFS(valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "fs=") {
		color.Red("[MKFS]: No tiene FS Especificado")
		return ""
	}

	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[MKFS]: No tiene FS Especificado")
		return ""
	}

	fs_val := strings.ToUpper(value[1])
	if !(fs_val == "3FS" || fs_val == "2FS") {
		color.Red("[MKSF]: No tiene FS valido")
		return ""
	}

	if fs_val == "3FS" {
		return "3FS"
	} else { // Incluye "2FS"
		return "2FS"
	}
}

// TieneUser extrae el nombre de usuario del parámetro -user.
func TieneUser(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "user=") {
		color.Red("[" + comando + "]: No tiene user o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene user Valido")
		return ""
	}
	return value[1]
}

// TienePassword extrae la contraseña del parámetro -pass.
func TienePassword(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "pass=") {
		color.Red("[" + comando + "]: No tiene password o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene password Valido")
		return ""
	}
	return value[1]
}
