package adminfiles

import (
	"Proyecto/Estructuras/size"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"encoding/binary"
	"os"
	"strings"

	"github.com/fatih/color"
)

func Values_MKDIR(instructions []string) (string, bool, bool) {
	var path string
	var p bool

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("MKDIR", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "p") {
			var value = utils.TienePPermitionsFile("MKDIR", valor)
			p = value
		} else {
			color.Yellow("[MKDIR]: Atributo no reconocido")
			return "", false, false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[MKDIR]: No hay path")
		return "", false, false
	}
	return path, p, true
}

func MKDIR_EXECUTE(path string, p bool) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[MKDIR]: Usuario no logeado")
		return
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[MKDIR]: Error al abrir archivo")
		return
	}
	defer file.Close()

	var start int32
	if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[MKDIR]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_L.Part_start + size.SizeEBR()
		}
	} else if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[MKDIR]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_P.Part_start
		}
	}

	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[MKDIR]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[MKDIR]: Error en la lectura del SuperBloque")
		return
	}

	rutaS := utils.SplitRuta(path)
	if len(rutaS) == 0 {
		color.Red("[MKDIR]: Ruta invalida")
		return
	}

	exist := utils.GetInodoFSystem(rutaS, 0, int32(len(rutaS)-1), utils.Sb_System.S_inode_start, nodo.Path)
	if exist != -1 {
		color.Red("[MKDIR]: La ruta ya existe")
		return
	}

	posInodoI := utils.Sb_System.S_inode_start
	var existP bool = true
	if len(rutaS) > 1 {
		for i := 0; i < (len(rutaS) - 1); i++ {
			if existP {
				aux := posInodoI
				posInodoI = utils.ExistPathSystem(rutaS, int32(i), posInodoI, nodo.Path)
				if posInodoI == aux {
					existP = false
				}
			}
			if !existP {
				if p {
					posInodoI = utils.CrearCarpetaSystem(rutaS, int32(i), posInodoI, nodo.Path)
					if nodo.Es_Particion_P {
						if _, err := file.Seek(int64(nodo.Particion_P.Part_start), 0); err != nil {
							color.Red("[MKDIR]: Error en mover puntero")
							return
						}
					} else if nodo.Es_Particion_L {
						if _, err := file.Seek(int64(nodo.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
							color.Red("[MKDIR]: Error en mover puntero")
							return
						}
					}
					if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
						color.Red("[MKDIR]: Error en la escritura del archivo")
						return
					}
					if posInodoI == -1 {
						return
					}
				} else {
					color.Red("[MKDIR]: no se puede crear la carpeta: " + rutaS[i])
					return
				}
			}
		}
	}

	if posInodoI == -1 {
		color.Red("[MKDIR]: Algo salio mal")
		return
	}

	utils.CrearCarpetaSystem(rutaS, int32(len(rutaS)-1), posInodoI, nodo.Path)
	if nodo.Es_Particion_P {
		if _, err := file.Seek(int64(nodo.Particion_P.Part_start), 0); err != nil {
			color.Red("[MKDIR]: Error en mover puntero")
			return
		}
	} else if nodo.Es_Particion_L {
		if _, err := file.Seek(int64(nodo.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
			color.Red("[MKDIR]: Error en mover puntero")
			return
		}
	}

	if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[MKDIR]: Error en la escritura del archivo")
		return
	}

	color.Green("[MKDIR]: Creada la carpeta -> " + path)

}
