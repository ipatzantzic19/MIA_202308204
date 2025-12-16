package adminDisk

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Values_MKFS analiza y valida los parámetros del comando MKFS
func Values_MKFS(instructions []string) (string, string, error) {
	var _id string
	var _type = "FULL" // Valor por defecto

	for _, valor := range instructions {
		param := strings.ToLower(strings.TrimSpace(valor))

		if strings.HasPrefix(param, "id=") {
			parts := strings.Split(valor, "=")
			if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
				return "", "", fmt.Errorf("[MKFS]: el parámetro -id es obligatorio y no puede estar vacío")
			}
			_id = strings.TrimSpace(parts[1])
		} else if strings.HasPrefix(param, "type=") {
			parts := strings.Split(valor, "=")
			if len(parts) >= 2 {
				typeValue := strings.ToUpper(strings.TrimSpace(parts[1]))
				if typeValue == "FULL" {
					_type = typeValue
				} else {
					return "", "", fmt.Errorf("[MKFS]: type inválido, solo se acepta 'full'")
				}
			}
		} else {
			return "", "", fmt.Errorf("[MKFS]: parámetro no reconocido: %s", valor)
		}
	}

	// Validación final
	if _id == "" {
		return "", "", fmt.Errorf("[MKFS]: el parámetro -id es obligatorio")
	}

	if len(_id) > 4 {
		return "", "", fmt.Errorf("[MKFS]: el id no puede tener más de 4 caracteres")
	}

	return _id, _type, nil
}

// MKFS_EXECUTE ejecuta el formateo de la partición
func MKFS_EXECUTE(id_disco string, tipo_formateo string) error {
	// Buscar la partición montada con el ID proporcionado
	nodoM, encontrado := Buscar_ID_Montada(id_disco)
	if !encontrado {
		return fmt.Errorf("[MKFS]: partición con id '%s' no existe o no está montada", id_disco)
	}

	// Abrir el archivo del disco
	file, err := os.OpenFile(nodoM.Path, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("[MKFS]: error al abrir el archivo del disco: %v", err)
	}
	defer file.Close()

	// Determinar si es partición primaria o lógica y ejecutar el formateo correspondiente
	if nodoM.Es_Particion_P {
		return formatearParticionPrimaria(file, nodoM, id_disco)
	} else if nodoM.Es_Particion_L {
		return formatearParticionLogica(file, nodoM, id_disco)
	}

	return fmt.Errorf("[MKFS]: tipo de partición no reconocido")
}

// formatearParticionPrimaria formatea una partición primaria
func formatearParticionPrimaria(file *os.File, nodoM global.ParticionesMontadas, id_disco string) error {
	tamanio := nodoM.Particion_P.Part_s

	// Calcular número de estructuras (inodos y bloques)
	// n = (partition_size - superblock_size) / (4 + inode_size + 3*block_size)
	// donde 4 es para los bitmaps (1 byte de bitmap de inodos + 3 bytes de bitmap de bloques)
	n := float32(tamanio-size.SizeSuperBloque()) / float32(4+size.SizeTablaInodo()+(3*size.SizeBloqueArchivo()))
	numeroEstructuras := int32(math.Floor(float64(n)))
	nBloques := 3 * numeroEstructuras

	// Crear el superbloque
	sb := structures.SuperBloque{
		S_filesistem_type:   2, // ext2
		S_inodes_count:      numeroEstructuras,
		S_blocks_count:      nBloques,
		S_free_blocks_count: nBloques - 2,          // Se usan 2: raíz y users.txt
		S_free_inodes_count: numeroEstructuras - 2, // Se usan 2: raíz y users.txt
		S_mtime:             utils.ObFechaInt(),
		S_umtime:            0,
		S_mnt_count:         1,
		S_magic:             0xEF53,
		S_inode_s:           size.SizeTablaInodo(),
		S_block_s:           size.SizeBloqueArchivo(),
		S_first_ino:         2, // El primer inodo libre (0 y 1 ocupados)
		S_first_blo:         2, // El primer bloque libre (0 y 1 ocupados)
	}

	// Calcular posiciones en el disco
	sb.S_bm_inode_start = nodoM.Particion_P.Part_start + size.SizeSuperBloque()
	sb.S_bm_block_start = sb.S_bm_inode_start + numeroEstructuras
	sb.S_inode_start = sb.S_bm_block_start + nBloques
	sb.S_block_start = sb.S_inode_start + (numeroEstructuras * size.SizeTablaInodo())

	// Limpiar la partición (llenar con ceros)
	if err := limpiarParticion(file, nodoM.Particion_P.Part_start, nodoM.Particion_P.Part_s); err != nil {
		return err
	}

	// Crear el inodo raíz (directorio)
	inodoRaiz := crearInodoRaiz(sb.S_inode_start)
	if err := escribirEstructura(file, sb.S_inode_start, &inodoRaiz); err != nil {
		return err
	}

	// Crear el bloque de carpeta raíz
	carpetaRaiz := crearCarpetaRaiz(sb.S_inode_start, size.SizeTablaInodo())
	if err := escribirEstructura(file, sb.S_block_start, &carpetaRaiz); err != nil {
		return err
	}

	// Crear el inodo para users.txt
	inodoUsers := crearInodoUsers(sb.S_block_start + size.SizeBloqueCarpeta())
	if err := escribirEstructura(file, sb.S_inode_start+size.SizeTablaInodo(), &inodoUsers); err != nil {
		return err
	}

	// Crear el bloque de archivo users.txt
	archivoUsers := crearArchivoUsers()
	if err := escribirEstructura(file, sb.S_block_start+size.SizeBloqueCarpeta(), &archivoUsers); err != nil {
		return err
	}

	// Actualizar el MBR para marcar la partición como formateada
	if err := actualizarMBRParticionPrimaria(file, nodoM, id_disco); err != nil {
		return err
	}

	// Escribir el superbloque
	if err := escribirEstructura(file, nodoM.Particion_P.Part_start, &sb); err != nil {
		return err
	}

	// Escribir los bitmaps
	if err := escribirBitmaps(file, sb, numeroEstructuras, nBloques); err != nil {
		return err
	}

	nombreParticion := utils.ToString(nodoM.Particion_P.Part_name[:])
	color.Green("[MKFS]: Partición '%s' (id: %s) formateada exitosamente", nombreParticion, id_disco)
	return nil
}

// formatearParticionLogica formatea una partición lógica
func formatearParticionLogica(file *os.File, nodoM global.ParticionesMontadas, id_disco string) error {
	tamanio := nodoM.Particion_L.Part_s

	// Calcular número de estructuras
	n := float32(tamanio-size.SizeEBR()-size.SizeSuperBloque()) / float32(4+size.SizeTablaInodo()+(3*size.SizeBloqueArchivo()))
	numeroEstructuras := int32(math.Floor(float64(n)))
	nBloques := 3 * numeroEstructuras

	// Crear el superbloque
	sb := structures.SuperBloque{
		S_filesistem_type:   2, // ext2
		S_inodes_count:      numeroEstructuras,
		S_blocks_count:      nBloques,
		S_free_blocks_count: nBloques - 2,
		S_free_inodes_count: numeroEstructuras - 2,
		S_mtime:             utils.ObFechaInt(),
		S_umtime:            0,
		S_mnt_count:         1,
		S_magic:             0xEF53,
		S_inode_s:           size.SizeTablaInodo(),
		S_block_s:           size.SizeBloqueArchivo(),
		S_first_ino:         2,
		S_first_blo:         2,
	}

	// Calcular posiciones (después del EBR)
	inicioSuper := nodoM.Particion_L.Part_start + size.SizeEBR()
	sb.S_bm_inode_start = inicioSuper + size.SizeSuperBloque()
	sb.S_bm_block_start = sb.S_bm_inode_start + numeroEstructuras
	sb.S_inode_start = sb.S_bm_block_start + nBloques
	sb.S_block_start = sb.S_inode_start + (numeroEstructuras * size.SizeTablaInodo())

	// Limpiar la partición (después del EBR)
	if err := limpiarParticion(file, inicioSuper, nodoM.Particion_L.Part_s-size.SizeEBR()); err != nil {
		return err
	}

	// Crear estructuras (igual que partición primaria)
	inodoRaiz := crearInodoRaiz(sb.S_inode_start)
	if err := escribirEstructura(file, sb.S_inode_start, &inodoRaiz); err != nil {
		return err
	}

	carpetaRaiz := crearCarpetaRaiz(sb.S_inode_start, size.SizeTablaInodo())
	if err := escribirEstructura(file, sb.S_block_start, &carpetaRaiz); err != nil {
		return err
	}

	inodoUsers := crearInodoUsers(sb.S_block_start + size.SizeBloqueCarpeta())
	if err := escribirEstructura(file, sb.S_inode_start+size.SizeTablaInodo(), &inodoUsers); err != nil {
		return err
	}

	archivoUsers := crearArchivoUsers()
	if err := escribirEstructura(file, sb.S_block_start+size.SizeBloqueCarpeta(), &archivoUsers); err != nil {
		return err
	}

	// Actualizar el EBR
	if err := actualizarEBRParticionLogica(file, nodoM); err != nil {
		return err
	}

	// Escribir el superbloque
	if err := escribirEstructura(file, inicioSuper, &sb); err != nil {
		return err
	}

	// Escribir los bitmaps
	if err := escribirBitmaps(file, sb, numeroEstructuras, nBloques); err != nil {
		return err
	}

	nombreParticion := utils.ToString(nodoM.Particion_L.Name[:])
	color.Green("[MKFS]: Partición lógica '%s' (id: %s) formateada exitosamente como ext2", nombreParticion, id_disco)
	return nil
}

// ==================== FUNCIONES AUXILIARES ====================

// Buscar_ID_Montada busca una partición montada por su ID
func Buscar_ID_Montada(id_disco string) (global.ParticionesMontadas, bool) {
	idBytes := utils.IDParticionByte(id_disco)
	for _, disco := range global.Mounted_Partitions {
		if disco.ID_Particion == idBytes {
			return disco, true
		}
	}
	return global.ParticionesMontadas{}, false
}

// NameCarpeta12 convierte un string a un arreglo de 12 bytes
func NameCarpeta12(nombre string) [12]byte {
	temp := make([]byte, 12)
	for i := range temp {
		temp[i] = '\x00'
	}
	copy(temp[:], []byte(nombre))
	return [12]byte(temp)
}

// limpiarParticion llena con ceros la región especificada del disco
func limpiarParticion(file *os.File, inicio int32, tamanio int32) error {
	// Dejar espacio para estructuras mínimas
	tamanioLimpieza := tamanio - 2
	if tamanioLimpieza <= 0 {
		return fmt.Errorf("[MKFS]: partición demasiado pequeña")
	}

	estructura := make([]byte, tamanioLimpieza)
	for i := range estructura {
		estructura[i] = '\x00'
	}

	if _, err := file.Seek(int64(inicio), 0); err != nil {
		return fmt.Errorf("[MKFS]: error al mover el puntero: %v", err)
	}

	if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
		return fmt.Errorf("[MKFS]: error al limpiar la partición: %v", err)
	}

	return nil
}

// crearInodoRaiz crea el inodo del directorio raíz
func crearInodoRaiz(bloqueInicio int32) structures.TablaInodo {
	inodo := structures.TablaInodo{
		I_uid:   1,
		I_gid:   1,
		I_s:     0,
		I_atime: utils.ObFechaInt(),
		I_ctime: utils.ObFechaInt(),
		I_mtime: utils.ObFechaInt(),
		I_type:  0, // 0 = carpeta
		I_perm:  664,
	}

	// Inicializar bloques en -1
	for i := range inodo.I_block {
		inodo.I_block[i] = -1
	}
	inodo.I_block[0] = bloqueInicio // Apunta al primer bloque

	return inodo
}

// crearCarpetaRaiz crea el bloque de carpeta raíz
func crearCarpetaRaiz(inodoInicio int32, sizeInodo int32) structures.BloqueCarpeta {
	carpeta := structures.BloqueCarpeta{}

	// Entrada "." (directorio actual)
	carpeta.B_content[0].B_name = NameCarpeta12(".")
	carpeta.B_content[0].B_inodo = inodoInicio

	// Entrada ".." (directorio padre, apunta a sí mismo en raíz)
	carpeta.B_content[1].B_name = NameCarpeta12("..")
	carpeta.B_content[1].B_inodo = inodoInicio

	// Entrada "users.txt"
	carpeta.B_content[2].B_name = NameCarpeta12("users.txt")
	carpeta.B_content[2].B_inodo = inodoInicio + sizeInodo

	// Entrada vacía
	carpeta.B_content[3].B_name = NameCarpeta12("")
	carpeta.B_content[3].B_inodo = -1

	return carpeta
}

// crearInodoUsers crea el inodo para el archivo users.txt
func crearInodoUsers(bloqueInicio int32) structures.TablaInodo {
	contenido := "1,G,root\n1,U,root,root,123\n"

	inodo := structures.TablaInodo{
		I_uid:   1,
		I_gid:   1,
		I_s:     int32(len(contenido)),
		I_atime: utils.ObFechaInt(),
		I_ctime: utils.ObFechaInt(),
		I_mtime: utils.ObFechaInt(),
		I_type:  1, // 1 = archivo
		I_perm:  664,
	}

	// Inicializar bloques en -1
	for i := range inodo.I_block {
		inodo.I_block[i] = -1
	}
	inodo.I_block[0] = bloqueInicio

	return inodo
}

// crearArchivoUsers crea el bloque de archivo users.txt con el contenido inicial
func crearArchivoUsers() structures.BloqueArchivo {
	contenido := "1,G,root\n1,U,root,root,123\n"
	return structures.BloqueArchivo{
		B_content: utils.DevolverContenidoArchivo(contenido),
	}
}

// escribirEstructura escribe una estructura en la posición especificada
func escribirEstructura(file *os.File, posicion int32, estructura interface{}) error {
	if _, err := file.Seek(int64(posicion), 0); err != nil {
		return fmt.Errorf("[MKFS]: error al mover el puntero: %v", err)
	}

	if err := binary.Write(file, binary.LittleEndian, estructura); err != nil {
		return fmt.Errorf("[MKFS]: error al escribir estructura: %v", err)
	}

	return nil
}

// escribirBitmaps escribe los bitmaps de inodos y bloques
func escribirBitmaps(file *os.File, sb structures.SuperBloque, numInodos int32, numBloques int32) error {
	var ch0 byte = '0'
	var ch1 byte = '1'

	// Inicializar bitmap de inodos con '0'
	for i := int32(0); i < numInodos; i++ {
		if err := escribirEstructura(file, sb.S_bm_inode_start+i, &ch0); err != nil {
			return err
		}
	}

	// Marcar los primeros 2 inodos como usados
	if err := escribirEstructura(file, sb.S_bm_inode_start, &ch1); err != nil {
		return err
	}
	if err := escribirEstructura(file, sb.S_bm_inode_start+1, &ch1); err != nil {
		return err
	}

	// Inicializar bitmap de bloques con '0'
	for i := int32(0); i < numBloques; i++ {
		if err := escribirEstructura(file, sb.S_bm_block_start+i, &ch0); err != nil {
			return err
		}
	}

	// Marcar los primeros 2 bloques como usados
	if err := escribirEstructura(file, sb.S_bm_block_start, &ch1); err != nil {
		return err
	}
	if err := escribirEstructura(file, sb.S_bm_block_start+1, &ch1); err != nil {
		return err
	}

	return nil
}

// actualizarMBRParticionPrimaria actualiza el MBR para marcar la partición como formateada
func actualizarMBRParticionPrimaria(file *os.File, nodoM global.ParticionesMontadas, id_disco string) error {
	mbr, ok := utils.Obtener_FULL_MBR_FDISK(nodoM.Path)
	if !ok {
		return fmt.Errorf("[MKFS]: error al leer el MBR")
	}

	// Buscar la partición y actualizar su estado
	nombreParticion := utils.ToString(nodoM.Particion_P.Part_name[:])
	for i := range mbr.Mbr_partitions {
		if utils.ToString(mbr.Mbr_partitions[i].Part_name[:]) == nombreParticion {
			mbr.Mbr_partitions[i].Part_status = 1 // 1 = formateada
			break
		}
	}

	// Escribir el MBR actualizado
	if err := escribirEstructura(file, 0, &mbr); err != nil {
		return err
	}

	// Actualizar en la lista global de particiones montadas
	for z := range global.Mounted_Partitions {
		if utils.ToString(global.Mounted_Partitions[z].Particion_P.Part_id[:]) == id_disco {
			global.Mounted_Partitions[z].Particion_P.Part_status = 1
			break
		}
	}

	return nil
}

// actualizarEBRParticionLogica actualiza el EBR para marcar la partición lógica como formateada
func actualizarEBRParticionLogica(file *os.File, nodoM global.ParticionesMontadas) error {
	ebr := nodoM.Particion_L
	ebr.Part_mount = 1 // 1 = formateada

	if err := escribirEstructura(file, nodoM.Particion_L.Part_start, &ebr); err != nil {
		return err
	}

	// Actualizar en la lista global
	nombreParticion := utils.ToString(nodoM.Particion_L.Name[:])
	for z := range global.Mounted_Partitions {
		if utils.ToString(global.Mounted_Partitions[z].Particion_L.Name[:]) == nombreParticion {
			global.Mounted_Partitions[z].Particion_L.Part_mount = 1
			break
		}
	}

	return nil
}
