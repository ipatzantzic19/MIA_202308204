package adminfiles

import (
	"Proyecto/Estructuras/size"
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

func Values_CAT(instructions []string) ([]string, bool) {
	var stringVacio []string
	var stringEnviar []string
	re := regexp.MustCompile(`file\d*=.*`)
	for _, value := range instructions {
		if re.MatchString(value) {
			val := utils.TieneFile(value)
			stringEnviar = append(stringEnviar, val)
		} else {
			color.Yellow("[CAT]: Valor no reconocido")
			return stringVacio, false
		}
	}
	return stringEnviar, true
}

func CAT_EXECUTE(files []string) bool {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[CAT]: Usuario no logeado")
		return false
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[CAT]: Error al abrir archivo")
		return false
	}
	defer file.Close()

	for _, path_xyz := range files {
		var start int32
		if nodo.Es_Particion_L {
			if nodo.Particion_L.Part_mount != 1 {
				color.Red("[CAT]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
				return false
			} else {
				start = nodo.Particion_L.Part_start + size.SizeEBR()
			}
		} else if nodo.Es_Particion_P {
			if nodo.Particion_P.Part_status != 1 {
				color.Red("[CAT]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
				return false
			} else {
				start = nodo.Particion_P.Part_start
			}
		}

		if _, err := file.Seek(int64(start), 0); err != nil {
			color.Red("[CAT]: Error en mover puntero")
			return false
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
			color.Red("[CAT]: Error en la lectura del SuperBloque")
			return false
		}

		rutaS := utils.SplitRuta(path_xyz)
		if len(rutaS) == 0 {
			color.Red("[CAT]: Ruta invalida")
			return false
		}

		var inodo structures.TablaInodo
		var posInodoF int32

		if path_xyz != "/" {
			posInodoF = utils.GetInodoFSystem(rutaS, 0, int32(len(rutaS)-1), utils.Sb_System.S_inode_start, nodo.Path)
			if posInodoF == -1 {
				color.Red("[CAT]: Archivo no encontrado")
				return false
			}
		} else {
			posInodoF = utils.Sb_System.S_inode_start
		}

		if !utils.ValidarPermisoRSystem(posInodoF, nodo.Path) {
			color.Red("[CAT]: No se puede leer el archivo -> <<<" + path_xyz + ">>> por falta de permisos")
			return false
		}

		if _, err := file.Seek(int64(posInodoF), 0); err != nil {
			color.Red("[CAT]: Error en mover puntero")
			return false
		}
		if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[CAT]: Error en la lectura del archivo")
			return false
		}

		if inodo.I_type == 0 {
			color.Red("[CAT]: La direccion no hace referencia a un archivo -> " + path_xyz)
			return false
		}

		inodo.I_atime = utils.ObFechaInt()
		content := utils.GetContentAdminUsers(posInodoF)
		color.Magenta("CAT >-----------------------------------")
		color.Blue("------------------- > " + rutaS[len(rutaS)-1] + " < -------------------")
		fmt.Println(content)
		color.Blue("------------------- > " + rutaS[len(rutaS)-1] + " < -------------------")
		color.Magenta("CAT >-----------------------------------")

		if _, err := file.Seek(int64(posInodoF), 0); err != nil {
			color.Red("[MKDIR]: Error en mover puntero")
			return false
		}
		if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[MKDIR]: Error en la escritura del archivo")
			return false
		}
	}
	return true
}
