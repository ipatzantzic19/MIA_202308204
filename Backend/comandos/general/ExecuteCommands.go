package general

import (
	"strings"

	"github.com/fatih/color"
)

var commandGroups = map[string][]string{
	"disk":    {"mkdisk", "fdisk", "rmdisk", "mount", "mounted", "mkfs"},
	"reports": {"rep"},
	"files":   {"mkfile", "mkdir"},
	"cat":     {"cat"},
	"users":   {"login", "logout"},
	"groups":  {"mkgrp", "mkusr"},
}

func detectGroup(cmd string) (string, string, bool, string) {
	cmdLower := strings.ToLower(cmd)

	for group, cmds := range commandGroups {
		for _, prefix := range cmds {
			if strings.HasPrefix(cmdLower, prefix) {
				return group, prefix, false, ""
			}
		}
	}

	return "", "", true, "Comando no reconocido"
}

// Retorna lista de errores, lista de mensajes exitosos, y contador de errores
func GlobalCom(lista []string) ([]string, []string, int) {
	var errores []string
	var mensajes []string
	var contErrores = 0

	for _, comm := range lista {
		group, command, blnError, strError := detectGroup(comm)
		if blnError {
			color.Red("Comando no reconocido %v", command)
			errores = append(errores, strError)
			contErrores++
		}

		//comandos := ObtenerComandos(comm)
		// fmt.Println("00000000000000000000000000000")
		// fmt.Println(group, command)
		switch group {
		case "disk":
			color.Cyan("Administración de discos: %v", command)
			// Enrutar comandos de disco a los handlers correspondientes
			switch command {
			case "mkdisk":
				msg, err := MkDisk(comm)
				if err != nil {
					color.Red("Error mkdisk: %v", err)
					errores = append(errores, err.Error())
					contErrores++
				} else {
					color.Green(msg)
					mensajes = append(mensajes, msg)
				}
			case "rmdisk":
				msg, err := RmDisk(comm)
				if err != nil {
					color.Red("Error rmdisk: %v", err)
					errores = append(errores, err.Error())
					contErrores++
				} else {
					color.Green(msg)
					mensajes = append(mensajes, msg)
				}
			default:
				color.Yellow("Comando de disco no implementado: %v", command)
			}

		case "reports":
			color.Red("Administración de reportes: %v", command)

		case "files":
			color.Green("Administración de Archivos: %v", command)

		case "cat":
			color.Blue("Comando CAT")

		case "users":
			color.Yellow("Administración de Usuarios: %v", command)

		case "groups":
			color.White("Administración de Grrupos: %v", command)
		}
	}

	return errores, mensajes, contErrores
}
