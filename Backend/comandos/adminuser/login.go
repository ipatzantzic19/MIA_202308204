package adminusers

import (
	"Proyecto/Estructuras/size"
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Values_LOGIN analiza las instrucciones proporcionadas para el comando LOGIN.
// Extrae y valida los parámetros -user, -pass, y -id.
// Retorna una estructura de usuario temporal y un booleano que indica si el análisis fue exitoso.
func Values_LOGIN(instructions []string) (global.Usuario, bool) {
	var usuario string
	var password string
	var id string

	user_temp := global.Usuario{}
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = utils.TieneID("LOGIN", valor)
			id = value
			continue
		} else if strings.HasPrefix(strings.ToLower(valor), "pass") {
			var value = utils.TienePassword("LOGIN", valor)
			password = value
			continue
		} else if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = utils.TieneUser("LOGIN", valor)
			usuario = value
			continue
		} else {
			color.Yellow("[LOGIN]: Atributo '%s' no reconocido.", valor)
			return user_temp, false
		}
	}

	// Validar que los parámetros obligatorios no estén vacíos
	if id == "" {
		color.Red("[LOGIN]: El parámetro '-id' es obligatorii")
		return user_temp, false
	} else if password == "" {
		color.Red("[LOGIN]: El parámetro '-pass' es obligatorio")
		return user_temp, false
	} else if usuario == "" {
		color.Red("[LOGIN]: El parámetro '-user' es obligatorio")
		return user_temp, false
	}

	// Poblar la estructura de usuario temporal con los datos parseados
	user_temp.ID_Particion = global.Global_ID(id)
	user_temp.User = global.Global_Data(usuario)
	user_temp.Password = global.Global_Data(password)
	return user_temp, true
}

// LOGIN_EXECUTE ejecuta la lógica principal para el inicio de sesión de un usuario.
// Realiza validaciones de sesión, busca la partición, lee el archivo de usuarios y autentica las credenciales.
func LOGIN_EXECUTE(usuario string, password string, id_disco string) bool {
	// 1. Validar que no haya una sesión activa
	if global.UsuarioLogeado.Logged_in {
		color.Red("[LOGIN]: Ya hay una sesión activa. Debe cerrar sesión (logout) antes de iniciar una nueva.")
		return false
	}

	// 2. Buscar la partición montada por su ID
	particion_montada, epm := utils.Buscar_ID_Montada(id_disco)
	if !epm {
		// El error específico ya se muestra dentro de la función Buscar_ID_Montada
		return false
	}

	// 3. Leer el SuperBloque de la partición para obtener la información del sistema de archivos
	superbloque, err_sb := utils.GetSuperBloque(particion_montada.Path, particion_montada)
	if err_sb != nil {
		color.Red("[LOGIN]: %v", err_sb)
		return false
	}

	// 4. Leer el contenido del archivo de usuarios (users.txt)
	// Por convención del sistema, el inodo 1 corresponde al archivo 'users.txt'.
	// Se accede al inodo directamente desde la tabla de inodos.
	contenido, eco := utils.GetContent(superbloque.S_inode_start+size.SizeTablaInodo(), particion_montada.Path)
	if !eco {
		color.Red("[LOGIN]: Error al leer el archivo de usuarios (users.txt).")
		return false
	}

	// 5. Procesar el archivo de usuarios para la autenticación
	contenido_split := strings.Split(contenido, "\n")
	var usuarioEncontrado bool = false
	var grupoEncontrado bool = false
	var datosUsuario []string

	// Iterar el contenido de users.txt para encontrar al usuario y verificar la contraseña
	for _, linea := range contenido_split {
		linea = strings.TrimSpace(linea)
		if linea == "" {
			continue
		}
		// El formato de usuario es: uid,tipo,grupo,usuario,contraseña
		if strings.Contains(linea, ",U,") {
			partes := strings.Split(linea, ",")
			if len(partes) == 5 && partes[3] == usuario {
				usuarioEncontrado = true
				if partes[4] == password {
					// Contraseña correcta, guardar datos del usuario para usarlos más adelante
					datosUsuario = partes
					break // Salir del bucle, usuario autenticado
				} else {
					// Usuario encontrado pero la contraseña no coincide
					color.Red("[LOGIN]: Autenticación fallida. Contraseña incorrecta para el usuario '%s'.", usuario)
					return false
				}
			}
		}
	}

	// Si el flag usuarioEncontrado es falso, el usuario no existe en el archivo
	if !usuarioEncontrado {
		color.Red("[LOGIN]: Autenticación fallida. El usuario '%s' no existe en la partición '%s'.", usuario, id_disco)
		return false
	}

	// 6. Si el usuario fue autenticado, buscar los datos de su grupo
	if usuarioEncontrado {
		nombreGrupo := datosUsuario[2]
		for _, linea := range contenido_split {
			linea = strings.TrimSpace(linea)
			if linea == "" {
				continue
			}
			// El formato de grupo es: gid,tipo,nombre_grupo
			if strings.Contains(linea, ",G,") {
				partes := strings.Split(linea, ",")
				if len(partes) == 3 && partes[2] == nombreGrupo {
					grupoEncontrado = true
					// Llenar la estructura global del grupo del usuario logueado
					gid, _ := strconv.Atoi(partes[0])
					global.GrupoUsuarioLoggeado.GID = int32(gid)
					global.GrupoUsuarioLoggeado.Tipo = 'G'
					global.GrupoUsuarioLoggeado.Nombre = partes[2]
					break // Salir del bucle, grupo encontrado
				}
			}
		}
	}

	// Si no se encuentra el grupo asociado, es un estado inconsistente del archivo users.txt
	if !grupoEncontrado {
		color.Red("[LOGIN]: Error interno. No se encontró el grupo '%s' asociado al usuario.", datosUsuario[2])
		global.UsuarioLogeado = global.DefaultUser // Resetear por seguridad
		return false
	}

	// 7. Todos los datos son correctos. Poblar la estructura global de la sesión del usuario.
	uid, _ := strconv.Atoi(datosUsuario[0])
	global.UsuarioLogeado.UID = int32(uid)
	global.UsuarioLogeado.GID = global.GrupoUsuarioLoggeado.GID // Asignar GID del grupo encontrado
	global.UsuarioLogeado.Tipo = 'U'
	global.UsuarioLogeado.Grupo = global.Global_Data(datosUsuario[2])
	global.UsuarioLogeado.User = global.Global_Data(datosUsuario[3])
	global.UsuarioLogeado.Password = global.Global_Data(datosUsuario[4])
	global.UsuarioLogeado.ID_Particion = global.Global_ID(id_disco)
	global.UsuarioLogeado.Mounted = particion_montada
	global.UsuarioLogeado.Logged_in = true

	color.Green("[LOGIN]: Usuario «%s» ha iniciado sesión exitosamente en la partición (id): -> %s", usuario, id_disco)
	return true
}
