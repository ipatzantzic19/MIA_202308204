// Mensaje es un paquete que contiene estructuras y funciones
// para manejar resultados de comandos y respuestas API.
package general

// Estructuras para manejar resultados de comandos y respuestas API
type Resultado struct {
	StrMensajeError string
	BlnError        bool
	Respuesta       interface{}
}

// Estructura para devolver la lista de comandos ejecutados y sus mensajes
type SalidaComandoEjecutado struct {
	LstComandos []string `json:"comandos"`
	Mensajes    []string `json:"mensajes"`
	Errores     []string `json:"errores"`
}

// Estructura para manejar el cuerpo de la solicitud API
type ResultadoAPI struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Crea una estructura ResultadoAPI para respuestas JSON
func ResultadoSalida(message string, isError bool, data interface{}) ResultadoAPI {
	return ResultadoAPI{
		Message: message,
		Error:   isError,
		Data:    data,
	}
}
