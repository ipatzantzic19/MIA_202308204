package general

import (
	"fmt"     // Para formatear strings
	"os"      // Para operaciones del sistema (crear carpetas, verificar existencia)
	"regexp"  // Para expresiones regulares
	"strings" // Para manipular strings

	"github.com/fatih/color"
)

var NamePath = "VDIC-MIA"
var ReportPath = "VDIC-MIA/Rep"
var DiskPath = "VDIC-MIA/Disks"

// Funcion que obtiene los comandos y sus atributos para su posterior ejecucion
func ObtenerComandos(x string) []string { //
	var comandos []string
	// Expresión regular para extraer atributos y sus valores de un comando
	atributos := regexp.MustCompile(`(-|>)(\w+)(?:="([^"]+)"|=(-?/?(\w+)?(?:/?[\w.-]+)*))?`).FindAllStringSubmatch(x, -1)
	for _, matches := range atributos {
		atributo := matches[2]         // Nombre del atributo
		valorConComillas := matches[3] // Valor entre comillas (si existe)
		valorSinComillas := matches[4] // Valor sin comillas (si existe)

		// Construir el string del comando según el tipo de valor
		if valorConComillas != "" {
			comandos = append(comandos, fmt.Sprintf("%s=%s", atributo, valorConComillas))
		} else if valorSinComillas != "" {
			comandos = append(comandos, fmt.Sprintf("%s=%s", atributo, valorSinComillas))
		} else {
			comandos = append(comandos, atributo)
		}
	}
	return comandos
}

// Obtiene el comando principal de una línea de comando
func getCommand(comm string, commands ...string) string {
	comm = strings.ToLower(comm)
	for _, c := range commands {
		if strings.HasPrefix(comm, c) {
			return c
		}
	}
	return ""
}

// Crea las carpetas necesarias para el funcionamiento del proyecto
func CrearCarpeta() {
	// nombre := "VDIC-MIA"
	// reportes := "VDIC-MIA/Rep"
	// discos := "VDIC-MIA/Disks"
	nombreArchivo := "VDIC-MIA/CarpetaImagenes.txt"

	// os.Stat verifica si una ruta existe
	// os.IsNotExist(err) devuelve true si la ruta NO existe
	// os.MkdirAll crea carpetas recursivamente (todas las intermedias)

	// Crear carpeta VDIC-MIA si no existe
	if _, err := os.Stat(NamePath); os.IsNotExist(err) {
		err := os.MkdirAll(NamePath, 0777) // crea la carpeta con permisos completos
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}
		color.Green("\t\t\t\t\tCarpeta VDIC-MIA creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\tCarpeta VDIC-MIA ya existente")
	}

	// Si no existe, crear carpeta Rep
	if _, err := os.Stat(ReportPath); os.IsNotExist(err) {
		err := os.Mkdir(ReportPath, 0777)
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}
		color.Green("\t\t\t\t\tCarpeta Rep creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\tCarpeta Rep ya existente")
	}

	// Si no existe, crear carpeta Disks
	if _, err := os.Stat(DiskPath); os.IsNotExist(err) {
		err := os.Mkdir(DiskPath, 0777)
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}
		color.Green("\t\t\t\t\tCarpeta VDIC-MIA/Disks creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\tCarpeta VDIC-MIA/Disks ya existente")
	}

	// Crear archivo CarpetaImagenes.txt si no existe
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		archivo, err := os.Create(nombreArchivo)
		if err != nil {
			fmt.Println("Error al crear archivo")
			return
		}
		defer archivo.Close()

		content := []byte("Proyecto Único\t\t\t\tCreated by Iskandar")
		_, err = archivo.Write(content)
		if err != nil {
			color.Red("Error escribiendo archivo:", err)
			return
		}
		color.Green("\t\t\t\t\tArchivo creado correctamente")
	} else {
		color.Yellow("\t\t\t\t\tArchivo existente")
	}
	color.Green("Finalizada la creación de carpetas")
}

// Verifica si el path proporcionado existe en el sistema
func TienePath(x string) string {
	y := strings.Split(x, "=")
	fmt.Print("\t\t\t\t\t\t\tBuscando:")
	color.Yellow(y[1])
	// Verifica si el archivo o carpeta existe
	if _, err := os.Stat(y[1]); os.IsNotExist(err) {
		color.Red("Archivo No Encontrado")
		return "nil"
	} else {
		color.Green("Archivo Encontrado")
		return y[1]
	}
}

// procesa una lista de comandos y luego las filtra eliminando comentarios y líneas vacías
func ExecuteCommandList(comandos []string) Resultado {
	var lineas []string
	// _ -> índice
	for _, comando := range comandos {
		linea := strings.TrimSpace(comando)
		if len(linea) > 0 && !strings.HasPrefix(linea, "#") {
			lineas = append(lineas, linea) // Agrega línea válida
		}
	}

	var exportar []string
	// Expresión regular para eliminar comentarios al final de las líneas
	reg := regexp.MustCompile(`(.*?)\s*(?:#.*|$)`)
	// Recorre las líneas filtradas
	for _, y := range lineas {
		match := reg.FindStringSubmatch(y) // Aplicar regex a la línea
		//fmt.Println(y, "asdf")
		if len(match) > 1 {
			exportar = append(exportar, match[1]) // Agrega la parte sin comentario
			//fmt.Println(match[0], "///", match[1])
		}
	}

	return Resultado{"", false, SalidaComandoEjecutado{LstComandos: exportar}}
}
