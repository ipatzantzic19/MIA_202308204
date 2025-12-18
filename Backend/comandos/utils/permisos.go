package utils

// funciones para permisos de archivos
import (
	"strings"

	"github.com/fatih/color"
)

func TieneRPermitionsFile(comando string, valor string) bool {
	if !strings.HasPrefix(strings.ToLower(valor), "r") {
		color.Red("[" + comando + "]: No tiene r o tiene un valor no valido")
		return false
	}
	return true
}

func TienePathFilePermitions(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "path=") {
		color.Red("[" + comando + "]: No tiene path o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene path Valido")
		return ""
	}
	return value[1]
}

func TieneGRP(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "grp=") {
		color.Red("[" + comando + "]: No tiene grp o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene grp Valido")
		return ""
	}
	return value[1]
}

func TieneUGO(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "ugo=") {
		color.Red("[" + comando + "]: No tiene ugo o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene ugo Valido")
		return ""
	}
	return value[1]
}
