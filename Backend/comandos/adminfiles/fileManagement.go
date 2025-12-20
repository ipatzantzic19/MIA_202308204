package adminfiles

import (
	"fmt"
	"strings"
)

// FilesCommandProps es el manejador principal para los comandos relacionados con archivos,
// como 'mkfile' (crear archivo) y 'mkdir' (crear directorio).
// Recibe el comando específico y una lista de sus propiedades (parámetros).
func FilesCommandProps(command string, props []string) (string, error) {
	// Se convierte el comando a minúsculas para un manejo de casos uniforme y no sensible a mayúsculas.
	switch strings.ToLower(command) {
	case "mkfile":
		path, r, size, cont, valid, overwrite := Values_MKFILE(props)
		if !valid {
			return "", fmt.Errorf("parámetros de 'mkfile' no válidos")
		}
		MKFILE_EXECUTE(path, r, size, cont, overwrite)
		return "Comando 'mkfile' ejecutado", nil

	case "mkdir":
		path, p, valid := Values_MKDIR(props)
		if !valid {
			return "", fmt.Errorf("parámetros de 'mkdir' no válidos")
		}
		MKDIR_EXECUTE(path, p)
		return "[MKDIR]: Comando ejecutado", nil

	case "cat":
		files, valid := Values_CAT(props)
		if !valid {
			return "", fmt.Errorf("parámetros de 'cat' no válidos")
		}
		CAT_EXECUTE(files)
		return "[CAT]: Comando ejecutado", nil

	default:
		// Si el comando no es 'mkfile' ni 'mkdir', se retorna un error para indicar que no es un comando de archivos válido.
		return "", fmt.Errorf("comando de archivos '%s' no reconocido", command)
	}
}
