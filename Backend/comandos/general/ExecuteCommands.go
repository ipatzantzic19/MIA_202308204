package general

import (
	adminDisk "Proyecto/comandos/admindisk"
	adminusers "Proyecto/comandos/adminuser"
	"Proyecto/comandos/commandGroups/files"
	"strings"

	"github.com/fatih/color"
)

// commandGroups es un mapa que agrupa los comandos por categoría.
// La clave del mapa es el nombre del grupo (por ejemplo, "disk") y el valor es una lista de los comandos que pertenecen a ese grupo.
var commandGroups = map[string][]string{
	"disk":    {"mkdisk", "fdisk", "rmdisk", "mount", "mounted", "mkfs"}, // Comandos relacionados con la administración de discos.
	"reports": {"rep"},                                                   // Comandos para generar reportes.
	"files":   {"mkfile", "mkdir"},                                       // Comandos para la administración de archivos.
	"cat":     {"cat"},                                                   // Comando para mostrar el contenido de los archivos.
	"users":   {"login", "logout"},                                       // Comandos para la administración de usuarios.
	"groups":  {"mkgrp", "mkusr"},                                        // Comandos para la administración de grupos.
}

// detectGroup identifica a qué grupo pertenece un comando dado.
// Recibe un comando como cadena y devuelve el grupo al que pertenece, el comando específico,
// un booleano que indica si hubo un error, y un mensaje de error si es necesario.
func detectGroup(cmd string) (string, string, bool, string) {
	cmdLower := strings.ToLower(strings.TrimSpace(cmd))
	cmdName := strings.Fields(cmdLower)[0] // Obtiene el nombre del comando

	// Itera sobre el mapa commandGroups para encontrar a qué grupo pertenece el comando.
	for group, cmds := range commandGroups {
		for _, prefix := range cmds {
			if cmdName == prefix {
				return group, prefix, false, ""
			}
		}
	}
	return "", "", true, "Comando no reconocido"
}

// GlobalCom procesa una lista de comandos, los ejecuta y devuelve los resultados.
// Retorna una lista de errores, una lista de mensajes de éxito y un contador de errores.
func GlobalCom(lista []string) ([]string, []string, int) {
	var errores []string  // Lista para almacenar los mensajes de error.
	var mensajes []string // Lista para almacenar los mensajes de éxito.
	var contErrores = 0   // Contador de errores.

	// Itera sobre cada comando en la lista de entrada.
	for _, comm := range lista {
		// Detecta el grupo y el comando específico.
		group, command, blnError, strError := detectGroup(comm)
		if blnError {
			color.Red("Comando no reconocido %v", command)
			errores = append(errores, strError)
			contErrores++
			continue // Continúa con el siguiente comando.
		}

		// Obtiene los parámetros del comando.
		comandos := ObtenerComandos(comm)
		// fmt.Println("00000000000000000000000000000")
		// fmt.Println(group, command)
		// Selecciona el manejador adecuado según el grupo del comando.
		switch group {
		case "disk":
			color.Cyan("Administración de discos: %v", command)
			// Enruta los comandos de disco a los manejadores correspondientes.
			msg, err := adminDisk.DiskCommandProps(strings.ToUpper(command), comandos)
			if err != nil {
				color.Red("Error "+command+": %v", err)
				errores = append(errores, err.Error())
				contErrores++
			} else {
				color.Green(msg)
				mensajes = append(mensajes, msg)
			}

		case "reports":
			color.Red("Administración de reportes: %v", command)

		case "files":
			color.Green("Administración de Archivos: %v", command)
			msg, err := files.FilesCommandProps(command, comandos)
			if err != nil {
				color.Red("Error: %v", err)
				errores = append(errores, err.Error())
				contErrores++
			} else {
				color.Green(msg)
				mensajes = append(mensajes, msg)
			}

		case "cat":
			color.Blue("Comando CAT")

		case "users":
			color.Yellow("Administración de Usuarios: %v", command)
			if !adminusers.UserCommandProps(command, comandos) {
				color.Red("Error en comando de usuario: %v", command)
				contErrores++
			}

		case "groups":
			color.White("Administración de Grrupos: %v", command)
		}
	}

	return errores, mensajes, contErrores // Devuelve las listas de errores, mensajes y el contador de errores.
}
