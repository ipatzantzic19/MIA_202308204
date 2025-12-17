package adminusers

import (
	"Proyecto/comandos/global"
	"Proyecto/comandos/utils"

	"github.com/fatih/color"
)

// LOGOUT_EXECUTE finaliza la sesión del usuario activo.
// Verifica si hay una sesión iniciada. Si no la hay, muestra un error.
// Si hay una sesión, la cierra restableciendo las variables globales de usuario y grupo a sus valores por defecto.
func LOGOUT_EXECUTE() bool {
	// 1. Validar que haya una sesión activa.
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[LOGOUT]: No hay una sesión activa. No se puede cerrar sesión.")
		return false
	}

	// 2. Guardar temporalmente los datos del usuario para el mensaje de despedida.
	usuario_temp := global.UsuarioLogeado

	// 3. Restablecer las variables globales a su estado por defecto.
	global.UsuarioLogeado = global.DefaultUser
	global.GrupoUsuarioLoggeado = global.DefaultGrupoUsuario

	// 4. Mostrar mensaje de éxito y retornar.
	color.Green("[LOGOUT]: Usuario «%s» ha cerrado sesión exitosamente de la partición (id): -> %s",
		utils.ToString(usuario_temp.User[:]),
		utils.ToString(usuario_temp.ID_Particion[:]))
	return true
}
