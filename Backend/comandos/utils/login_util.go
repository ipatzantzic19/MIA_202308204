package utils

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"encoding/binary"
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"
)

// GetContent lee y devuelve el contenido completo de un archivo dentro del sistema de archivos simulado.
// Comienza en la posición de un inodo específico y recorre los bloques de datos del archivo,
// manejando apuntadores directos e indirectos (simples, dobles y triples).
//
// Parámetros:
//   - inodoStart: La posición (byte de inicio) del inodo del archivo en el disco.
//   - path: La ruta al archivo de disco (.mia) que contiene el sistema de archivos.
//
// Retorna:
//   - Una cadena con el contenido del archivo concatenado.
//   - Un booleano que es `true` si la lectura fue exitosa, `false` si ocurrió un error.
func GetContent(inodoStart int32, path string) (string, bool) {
	// Declaración de estructuras necesarias para leer el archivo
	inodo := structures.TablaInodo{}
	apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}
	var content strings.Builder // Usar strings.Builder para una concatenación de strings más eficiente

	// Abrir el archivo de disco en modo lectura/escritura
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[GetContent]: Error al abrir el archivo de disco '%s'.", path)
		return "", false
	}
	defer file.Close()

	// Mover el puntero al inicio del inodo y leer su estructura
	if _, err := file.Seek(int64(inodoStart), 0); err != nil {
		color.Red("[GetContent]: Error al mover el puntero al inodo en la posición %d.", inodoStart)
		return "", false
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[GetContent]: Error al leer la estructura del inodo.")
		return "", false
	}

	// Iterar a través de los 15 apuntadores del inodo
	for i := 0; i < 15; i++ {
		// Procesar solo si el apuntador es válido (no es -1)
		if inodo.I_block[i] != -1 {
			posicionBloque := int64(inodo.I_block[i])

			// --- APUNTADORES DIRECTOS (0-11) ---
			if i < 12 {
				// Leer un bloque de archivo directo
				if !leerBloqueArchivo(file, posicionBloque, &content) {
					return "", false
				}
			} else if i == 12 {
				// --- APUNTADOR INDIRECTO SIMPLE (12) ---
				// Leer el bloque de apuntadores
				if err := leerEstructura(file, posicionBloque, &apuntador1); err != nil {
					return "", false
				}
				// Iterar a través de los apuntadores en este bloque
				for _, ptr := range apuntador1.B_pointers {
					if ptr != -1 {
						if !leerBloqueArchivo(file, int64(ptr), &content) {
							return "", false
						}
					}
				}
			} else if i == 13 {
				// --- APUNTADOR INDIRECTO DOBLE (13) ---
				// Nivel 1: Leer el primer bloque de apuntadores
				if err := leerEstructura(file, posicionBloque, &apuntador1); err != nil {
					return "", false
				}
				for _, ptr1 := range apuntador1.B_pointers {
					if ptr1 != -1 {
						// Nivel 2: Leer el segundo bloque de apuntadores
						if err := leerEstructura(file, int64(ptr1), &apuntador2); err != nil {
							return "", false
						}
						for _, ptr2 := range apuntador2.B_pointers {
							if ptr2 != -1 {
								// Leer el bloque de archivo final
								if !leerBloqueArchivo(file, int64(ptr2), &content) {
									return "", false
								}
							}
						}
					}
				}
			} else if i == 14 {
				// --- APUNTADOR INDIRECTO TRIPLE (14) ---
				// Nivel 1
				if err := leerEstructura(file, posicionBloque, &apuntador1); err != nil {
					return "", false
				}
				for _, ptr1 := range apuntador1.B_pointers {
					if ptr1 != -1 {
						// Nivel 2
						if err := leerEstructura(file, int64(ptr1), &apuntador2); err != nil {
							return "", false
						}
						for _, ptr2 := range apuntador2.B_pointers {
							if ptr2 != -1 {
								// Nivel 3
								if err := leerEstructura(file, int64(ptr2), &apuntador3); err != nil {
									return "", false
								}
								for _, ptr3 := range apuntador3.B_pointers {
									if ptr3 != -1 {
										// Leer el bloque de archivo final
										if !leerBloqueArchivo(file, int64(ptr3), &content) {
											return "", false
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return content.String(), true
}

// leerEstructura es una función auxiliar para leer cualquier estructura de datos binaria desde el archivo.
func leerEstructura(file *os.File, posicion int64, data interface{}) error {
	if _, err := file.Seek(posicion, 0); err != nil {
		color.Red("[Util]: Error al mover puntero a la posición %d.", posicion)
		return err
	}
	if err := binary.Read(file, binary.LittleEndian, data); err != nil {
		color.Red("[Util]: Error en la lectura de la estructura en la posición %d.", posicion)
		return err
	}
	return nil
}

// leerBloqueArchivo lee un bloque de archivo y añade su contenido al strings.Builder.
func leerBloqueArchivo(file *os.File, posicion int64, content *strings.Builder) bool {
	bloque := structures.BloqueArchivo{}
	if err := leerEstructura(file, posicion, &bloque); err != nil {
		return false
	}
	content.WriteString(ToString(bloque.B_content[:]))
	return true
}

func GetSuperBloque(path string, particion global.ParticionesMontadas) (structures.SuperBloque, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return structures.SuperBloque{}, err
	}
	defer file.Close()

	var start int32
	if particion.Es_Particion_P {
		start = particion.Particion_P.Part_start
	} else {
		start = particion.Particion_L.Part_start + size.SizeEBR()
	}

	_, err = file.Seek(int64(start), 0)
	if err != nil {
		return structures.SuperBloque{}, err
	}

	superbloque := structures.SuperBloque{}
	err = binary.Read(file, binary.LittleEndian, &superbloque)
	if err != nil {
		return structures.SuperBloque{}, errors.New("error al leer el superbloque")
	}

	if superbloque.S_magic != 0xEF53 {
		return structures.SuperBloque{}, errors.New("la partición no está formateada con ext2")
	}

	return superbloque, nil
}