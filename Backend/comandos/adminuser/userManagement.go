package adminusers

import (
	"Proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

// UserCommandProps actúa como controlador para los comandos de gestión de usuarios (login, logout).
// Analiza el comando principal y, basándose en él, llama a la función de ejecución correspondiente.
//
// Parámetros:
//   - command: El comando específico a ejecutar (ej. "login", "logout").
//   - instructions: Una lista de parámetros o atributos proporcionados con el comando.
//
// Retorna:
//   - Un booleano que indica si la ejecución del comando fue exitosa.
func UserCommandProps(command string, instructions []string) bool {
	if strings.ToUpper(command) == "LOGIN" {
		// Procesa el comando LOGIN
		valor_usuario, err := Values_LOGIN(instructions)
		if !err {
			// Si hay un error al analizar los parámetros de LOGIN, no se puede continuar.
			color.Red("[Login]: Error al analizar los parámetros del comando.")
			return false
		} else {
			// Extrae los valores y ejecuta la lógica de inicio de sesión.
			_user := utils.ToString(valor_usuario.User[:])
			_pass := utils.ToString(valor_usuario.Password[:])
			_id := utils.ToString(valor_usuario.ID_Particion[:])
			return LOGIN_EXECUTE(_user, _pass, _id)
		}
	} else if strings.ToUpper(command) == "LOGOUT" {
		// Procesa el comando LOGOUT
		return LOGOUT_EXECUTE()
	} else {
		// Si el comando no es ni LOGIN ni LOGOUT, es un error interno o un comando no reconocido.
		color.Red("[UserCommand]: Comando de usuario '%s' no reconocido.", command)
		return false
	}
}

// GroupCommandProps actúa como controlador para los comandos de gestión de grupos
func GroupCommandProps(group string, instructions []string) bool {
	if strings.ToUpper(group) == "MKGRP" {
		_name, valid := Values_MKGRP(instructions)
		if !valid {
			color.Red("[MKGRP]: Error al validar los parámetros")
			return false
		}
		MKGRP_EXECUTE(_name)
		return true

	} else {
		color.Red("[GroupCommandProps]: Comando de grupo '%s' no reconocido.", group)
		return false
	}
}
