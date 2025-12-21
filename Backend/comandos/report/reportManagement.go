package report

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func ReportCommandProps(command string, instructions []string) (string, error) {
	var _namereport, _name, _id string
	var er bool

	if strings.ToUpper(command) != "REP" {
		return "", errors.New("[ReportCommandProps]: Comando de reporte no reconocido")
	}

	_namereport, _name, _id, er = Sub_Reports(instructions)
	if !er {
		return "", errors.New("[REP]: Error en los parámetros. Faltan valores obligatorios como -namereport, -name, o -id")
	}

	return REP_EXECUTE(_namereport, _name, _id)
}

func Sub_Reports(instructions []string) (string, string, string, bool) {
	var _namereport, _name, _id string

	for _, valor := range instructions {
		parts := strings.SplitN(valor, "=", 2)
		if len(parts) < 2 {
			continue
		}

		val := strings.ToLower(parts[0])
		dat := parts[1]

		switch val {
		case "namereport":
			_namereport = dat
		case "name":
			_name = dat
		case "id":
			_id = dat
		default:
			color.Yellow("ADVERTENCIA: Parámetro de REP desconocido -> %s", val)
		}
	}

	if _id == "" || _name == "" || _namereport == "" {
		return "", "", "", false
	}

	allowedReports := map[string]bool{
		"mbr":      true,
		"disk":     true,
		"inode":    true,
		"block":    true,
		"bm_inode": true,
		"bm_block": true,
		"tree":     true,
		"sb":       true,
	}
	if !allowedReports[strings.ToLower(_namereport)] {
		color.Red("Error: Valor de -namereport no válido -> %s", _namereport)
		return "", "", "", false
	}

	return _namereport, _name, _id, true
}

func REP_EXECUTE(_namereport string, _name string, _id string) (string, error) {
	// La ruta base para los reportes es siempre "Rep"
	ruta_base := "VDIC-MIA/Rep"

	// Crear la carpeta base si no existe
	if _, err := os.Stat(ruta_base); os.IsNotExist(err) {
		err := os.MkdirAll(ruta_base, 0755)
		if err != nil {
			return "", fmt.Errorf("[REP]: Error al crear la carpeta base de reportes: %v", err)
		}
		color.Green("Carpeta base de reportes creada en: " + ruta_base)
	}

	ruta_completa_archivo := ruta_base + "/" + _name

	switch strings.ToLower(_namereport) {
	case "mbr":
		return Report_MBR(_id, ruta_completa_archivo)
	case "disk":
		return "Reporte 'disk' no implementado todavía.", nil
	case "inode":
		return "Reporte 'inode' no implementado todavía.", nil
	case "block":
		return "Reporte 'block' no implementado todavía.", nil
	case "bm_inode":
		return "Reporte 'bm_inode' no implementado todavía.", nil
	case "bm_block":
		return "Reporte 'bm_block' no implementado todavía.", nil
	case "tree":
		return "Reporte 'tree' no implementado todavía.", nil
	case "sb":
		return "Reporte 'sb' no implementado todavía.", nil
	default:
		return "", fmt.Errorf("[REP]: Nombre de reporte desconocido -> %s", _namereport)
	}
}
