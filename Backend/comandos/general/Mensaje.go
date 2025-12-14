// El paquete general contiene estructuras y funciones de utilidad
// que son compartidas a través de diferentes partes del sistema de manejo de comandos.
package general

// Resultado es una estructura genérica para encapsular el resultado de una operación.
// Se utiliza internamente para pasar información entre funciones antes de formatear la respuesta final de la API.
type Resultado struct {
	StrMensajeError string      // Mensaje de error, si lo hubo.
	BlnError        bool        // Un booleano que indica si la operación resultó en un error.
	Respuesta       interface{} // Contiene los datos de la respuesta exitosa. Es de tipo `interface{}` para ser flexible.
}

// SalidaComandoEjecutado define la estructura de los datos que se enviarán al frontend
// después de procesar y ejecutar una lista de comandos.
type SalidaComandoEjecutado struct {
	LstComandos []string `json:"comandos"` // La lista de comandos limpios que fueron procesados. `json:"comandos"` es una etiqueta de struct que personaliza el nombre del campo en la salida JSON.
	Mensajes    []string `json:"mensajes"` // Una lista de mensajes de éxito generados por los comandos.
	Errores     []string `json:"errores"`  // Una lista de mensajes de error generados por los comandos.
}

// ResultadoAPI define el formato estándar para todas las respuestas de la API.
// Esta estructura consistente facilita el manejo de respuestas en el lado del cliente (frontend).
type ResultadoAPI struct {
	Error bool        `json:"error"` // `true` si la solicitud falló en general, `false` si fue exitosa.
	Data  interface{} `json:"data"`  // Los datos específicos de la respuesta. Para los comandos, esto contendrá la estructura `SalidaComandoEjecutado`.
}

// ResultadoSalida es una función auxiliar para crear y poblar fácilmente una estructura `ResultadoAPI`.
// Ayuda a estandarizar la creación de respuestas de la API en todo el backend.
func ResultadoSalida(message string, isError bool, data interface{}) ResultadoAPI {
	return ResultadoAPI{
		Error: isError, // El indicador de error.
		Data:  data,    // Los datos a enviar.
	}
}
