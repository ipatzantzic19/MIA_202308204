package admingroup

import (
	"strings"

	"github.com/fatih/color"
)

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

	} else if strings.ToUpper(group) == "MKUSR" {
		_user, _pass, _grp, valid := Values_MKUSR(instructions)
		if !valid {
			color.Red("[MKUSR]: Error al validar los parámetros")
			return false
		}
		MKUSR_EXECUTE(_user, _pass, _grp)
		return true

	} else {
		color.Red("[GroupCommandProps]: Comando de grupo '%s' no reconocido.", group)
		return false
	}
}
