package adminusers

import (
	"Proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

// UserCommandProps actúa como un despachador para los comandos de gestión de usuarios (login, logout).
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
		// Procesa el comando LOGOUT, que no requiere parámetros.
		return LOGOUT_EXECUTE()
	} else {
		// Si el comando no es ni LOGIN ni LOGOUT, es un error interno o un comando no reconocido.
		color.Red("[UserCommand]: Comando de usuario '%s' no reconocido.", command)
		return false
	}
}

/*
func GroupCommandProps(group string, instructions []string) {
	var _name string //mkgrp rmgrp
	var er bool
	var _user string //mkusr rmusr
	var _pass string //mkusr
	var _grp string  //mkusr

	if strings.ToUpper(group) == "MKGRP" {
		_user, er = Values_MKGRP(instructions)
		if !er {
			color.Red("[MKGRP]: Error to assing values")
		} else {
			MKGRP_EXECUTE(_user)
		}
	} else if strings.ToUpper(group) == "RMGRP" {
		_name, er = Values_RMGRP(instructions)
		if !er {
			color.Red("[RMGRP]: Error to assing values")
		} else {
			RMGRP_EXECUTE(_name)
		}
	} else if strings.ToUpper(group) == "MKUSR" {
		_name, _pass, _grp, er = Values_MKUSR(instructions)
		if !er {
			color.Red("[MKUSR]: Error to asign values")
		} else {
			// fmt.Println(_name, _pass, _grp, er)
			MKUSR_EXECUTE(_name, _pass, _grp)
		}

	} else if strings.ToUpper(group) == "RMUSR" {
		_name, er := Values_RMUSR(instructions)
		if !er {
			color.Red("[RMUSR]: Error to asign values")
		} else {
			RMUSR_EXECUTE(_name)
		}
		// fmt.Println("Eliminando usuario en la parcicion")
	} else {
		color.Red("[GroupCommandProps]: Internal Error")
	}
}*/
