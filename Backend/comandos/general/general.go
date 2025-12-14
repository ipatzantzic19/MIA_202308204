package general

import (
	"fmt"     // Para formatear strings (cadenas de texto).
	"os"      // Para operaciones del sistema, como crear carpetas y verificar si existen.
	"regexp"  // Para trabajar con expresiones regulares, que permiten buscar patrones en texto.
	"strings" // Para manipular strings, como cortar, unir y comparar.

	"github.com/fatih/color" // Librería externa para mostrar texto en colores en la consola.
)

// Variables globales que definen las rutas de las carpetas principales del proyecto.
var NamePath = "VDIC-MIA"           // Carpeta raíz del proyecto.
var ReportPath = "VDIC-MIA/Rep"     // Carpeta para almacenar reportes.
var DiskPath = "VDIC-MIA/Disks"   // Carpeta para almacenar los archivos de disco (.dsk).

// ObtenerComandos extrae los atributos y sus valores de una línea de comando.
// Por ejemplo, de "mkdisk -size=3000 -unit=K", extraería ["size=3000", "unit=K"].
// Utiliza expresiones regulares para encontrar todos los parámetros que siguen el patrón -parametro=valor o -parametro="valor con espacios".
func ObtenerComandos(x string) []string {
	var comandos []string
	// La expresión regular busca patrones que comiencen con '-' o '>', seguido del nombre del atributo,
	// y opcionalmente un valor, que puede estar entre comillas o no.
	re := regexp.MustCompile(`(-|>)(\w+)(?:="([^"]+)"|=(-?/?(\w+)?(?:/?[\w.-]+)*))?`)
	atributos := re.FindAllStringSubmatch(x, -1)

	for _, matches := range atributos {
		atributo := matches[2]         // El nombre del atributo (ej: "size").
		valorConComillas := matches[3] // El valor si estaba entre comillas (ej: "mi nombre").
		valorSinComillas := matches[4] // El valor si no tenía comillas (ej: "3000").

		// Construye el string "atributo=valor" y lo añade a la lista.
		if valorConComillas != "" {
			comandos = append(comandos, fmt.Sprintf("%s=%s", atributo, valorConComillas))
		} else if valorSinComillas != "" {
			comandos = append(comandos, fmt.Sprintf("%s=%s", atributo, valorSinComillas))
		} else {
			// Si un atributo no tiene valor (es un flag), se añade solo el nombre.
			comandos = append(comandos, atributo)
		}
	}
	return comandos
}

// getCommand identifica el comando principal de una línea de entrada (ej: "mkdisk").
// Compara el inicio de la línea (en minúsculas) con una lista de comandos posibles.
func getCommand(comm string, commands ...string) string {
	comm = strings.ToLower(comm)
	for _, c := range commands {
		if strings.HasPrefix(comm, c) {
			return c // Devuelve el comando si encuentra una coincidencia.
		}
	}
	return "" // Devuelve una cadena vacía si no encuentra ninguna coincidencia.
}

// CrearCarpeta se encarga de inicializar la estructura de directorios necesaria para el proyecto.
// Verifica si las carpetas VDIC-MIA, Rep y Disks existen, y las crea si es necesario.
// También crea un archivo de texto informativo.
func CrearCarpeta() {
	nombreArchivo := "VDIC-MIA/CarpetaImagenes.txt"

	// os.Stat devuelve información sobre un archivo o directorio.
	// os.IsNotExist(err) devuelve true si el error es porque el archivo o directorio no existe.
	// os.MkdirAll crea una carpeta y todas las carpetas padres necesarias.

	// Crear la carpeta raíz "VDIC-MIA" si no existe.
	if _, err := os.Stat(NamePath); os.IsNotExist(err) {
		err := os.MkdirAll(NamePath, 0777) // 0777 da permisos de lectura, escritura y ejecución a todos.
		if err != nil {
			color.Red("Error al crear carpeta VDIC-MIA: %v", err)
			return
		}
		color.Green("\t\t\t\t\tCarpeta VDIC-MIA creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\tCarpeta VDIC-MIA ya existente")
	}

	// Crear la subcarpeta "Rep" si no existe.
	if _, err := os.Stat(ReportPath); os.IsNotExist(err) {
		err := os.Mkdir(ReportPath, 0777)
		if err != nil {
			color.Red("Error al crear carpeta Rep: %v", err)
			return
		}
		color.Green("\t\t\t\t\tCarpeta Rep creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\tCarpeta Rep ya existente")
	}

	// Crear la subcarpeta "Disks" si no existe.
	if _, err := os.Stat(DiskPath); os.IsNotExist(err) {
		err := os.Mkdir(DiskPath, 0777)
		if err != nil {
			color.Red("Error al crear carpeta Disks: %v", err)
			return
		}
		color.Green("\t\t\t\t\tCarpeta Disks creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\tCarpeta Disks ya existente")
	}

	// Crear un archivo de texto "CarpetaImagenes.txt" si no existe.
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		archivo, err := os.Create(nombreArchivo)
		if err != nil {
			fmt.Println("Error al crear archivo CarpetaImagenes.txt")
			return
		}
		defer archivo.Close() // Se asegura que el archivo se cierre al final de la función.

		content := []byte("Proyecto Único\t\t\t\tCreated by Iskandar")
		_, err = archivo.Write(content)
		if err != nil {
			color.Red("Error escribiendo en CarpetaImagenes.txt:", err)
			return
		}
		color.Green("\t\t\t\t\tArchivo CarpetaImagenes.txt creado correctamente")
	} else {
		color.Yellow("\t\t\t\t\tArchivo CarpetaImagenes.txt ya existente")
	}
	color.Green("Finalizada la creación de carpetas.")
}

// TienePath verifica si una ruta de archivo o directorio proporcionada en un parámetro (ej: "path=/ruta/a/mi/archivo") existe.
func TienePath(x string) string {
	y := strings.Split(x, "=") // Separa el parámetro en "path" y la ruta.
	if len(y) < 2 {
		color.Red("Formato de path incorrecto.")
		return "nil"
	}
	ruta := y[1]
	fmt.Print("\t\t\t\t\t\t\tBuscando: ")
	color.Yellow(ruta)

	// Verifica si el archivo o directorio existe.
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		color.Red(" -> Archivo No Encontrado")
		return "nil"
	} else {
		color.Green(" -> Archivo Encontrado")
		return ruta
	}
}

// ExecuteCommandList procesa una lista de comandos de entrada (como un script) y los filtra.
// Elimina líneas vacías, líneas que son solo comentarios (que empiezan con '#'),
// y también quita los comentarios que están al final de una línea de comando.
func ExecuteCommandList(comandos []string) Resultado {
	var lineasValidas []string
	// Primer filtro: eliminar líneas vacías y las que son completamente comentarios.
	for _, comando := range comandos {
		linea := strings.TrimSpace(comando) // Quita espacios al inicio y al final.
		if len(linea) > 0 && !strings.HasPrefix(linea, "#") {
			lineasValidas = append(lineasValidas, linea) // Agrega la línea si no está vacía y no es un comentario.
		}
	}

	var comandosLimpios []string
	// Expresión regular para capturar todo lo que está antes de un '#' en una línea.
	re := regexp.MustCompile(`(.*?)\s*(?:#.*|$)`)

	// Segundo filtro: eliminar comentarios al final de la línea.
	for _, y := range lineasValidas {
		match := re.FindStringSubmatch(y)
		if len(match) > 1 {
			comandoLimpio := strings.TrimSpace(match[1]) // Captura el comando sin el comentario.
			if len(comandoLimpio) > 0 {
				comandosLimpios = append(comandosLimpios, comandoLimpio)
			}
		}
	}

	// Devuelve el resultado en una estructura estandarizada.
	return Resultado{"", false, SalidaComandoEjecutado{LstComandos: comandosLimpios}}
}
