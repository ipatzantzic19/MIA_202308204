package controllers

import (
	"Proyecto/comandos/general"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
=============== CONTROLADOR DE COMANDOS ===============

HandleCommand es el endpoint HTTP que actúa como punto de entrada para todos los comandos enviados desde el Frontend.

RESPONSABILIDADES PRINCIPALES:
 1. Validar que la solicitud HTTP sea un método POST.
 2. Gestionar la configuración de CORS para permitir solicitudes desde orígenes diferentes (como el frontend en desarrollo).
 3. Extraer el contenido de los comandos desde el cuerpo JSON de la solicitud.
 4. Dividir el texto de comandos en líneas individuales para su procesamiento.
 5. Filtrar comentarios y líneas en blanco para obtener una lista limpia de comandos a ejecutar.
 6. Invocar al procesador global de comandos (GlobalCom) para ejecutar cada comando.
 7. Recopilar los resultados (mensajes de éxito y errores) de la ejecución.
 8. Devolver una respuesta JSON al frontend con la lista de comandos procesados y los resultados de su ejecución.

FLUJO DE DATOS:

	Frontend (UI)                Backend (Este Controlador)
	      |                                |
	      |---- POST /commands ----------->|
	      |   (JSON: {"Comandos":          |
	      |    "mkdisk -size=3000\n"       |
	      |    "# un comentario\n"         |
	      |    "fdisk -driveletter=A..."}) |
	      |                                |
	      |                         HandleCommand()
	      |                         - Valida método POST y CORS.
	      |                         - Decodifica el JSON.
	      |                         - Valida que el campo "Comandos" no esté vacío.
	      |                         - Divide el string en líneas.
	      |                         - Filtra comentarios y líneas vacías.
	      |                         - Llama a `general.GlobalCom()` para ejecutar.
	      |                         - Recopila errores y mensajes.
	      |                                |
	      |<--- Respuesta JSON ------------|
	      |   ({                             |
	      |     "Respuesta": {              |
	      |       "LstComandos": ["mkdisk...", "fdisk..."],
	      |       "Mensajes": ["Disco Creado..."],
	      |       "Errores": ["Error en fdisk..."]
	      |     },                          |
	      |     "Error": false             |
	      |   })                            |
	      |                                |
	Actualiza la UI con los resultados

=======================================================
*/
func HandleCommand(w http.ResponseWriter, r *http.Request) {
	/** =======================================================
		1. VALIDACIONES INICIALES DE CORS Y MÉTODO HTTP
	=======================================================	*/

	// Manejo de CORS (Cross-Origin Resource Sharing):
	// El navegador envía una solicitud "preflight" (de tipo OPTIONS) antes de una solicitud POST
	// si el frontend se encuentra en un dominio o puerto diferente al del backend.
	if r.Method == http.MethodOptions {
		// Se establecen los encabezados que le indican al navegador que la solicitud POST desde cualquier origen (*) es segura.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK) // Responde con 200 OK para que el navegador proceda con la solicitud POST real.
		return
	}

	// Solo se permiten solicitudes con el método POST.
	// Se rechaza cualquier otro método como GET, PUT, DELETE, etc.
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar los encabezados CORS para la respuesta a la solicitud POST.
	// Esto es crucial para que el navegador del cliente acepte la respuesta del servidor.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	/** =======================================================
		2. EXTRACCIÓN Y VALIDACIÓN DEL CUERPO JSON
	=======================================================	*/

	// Se espera que el frontend envíe un JSON con la siguiente estructura:
	// {
	//   "Comandos": "mkdisk -size=3000 ...\nfdisk ..."
	// }

	// Se define una estructura para mapear el JSON entrante.
	var requestBody struct {
		Comandos *string `json:"Comandos"` // Se usa un puntero (*) para poder diferenciar entre un campo ausente (nil) y un campo con valor vacío ("").
	}

	decoder := json.NewDecoder(r.Body) // Se crea un decodificador JSON que lee directamente del cuerpo de la solicitud (r.Body).
	decoder.DisallowUnknownFields()    // Se configura para que rechace cualquier campo en el JSON que no esté definido en la estructura `requestBody`.

	err := decoder.Decode(&requestBody) // Se decodifica el JSON en la estructura `requestBody`. El `&` pasa la dirección de memoria para que `Decode` pueda modificarla.

	// Manejo de errores durante la decodificación del JSON.
	if err != nil {
		w.Header().Set("Content-Type", "application/json") // Se asegura que la respuesta sea JSON.
		w.WriteHeader(http.StatusBadRequest)               // Código de error 400 (Solicitud Incorrecta).
		json.NewEncoder(w).Encode(                         // Se codifica y envía una respuesta de error estandarizada.
			general.ResultadoSalida("JSON inválido o campos no permitidos", true, nil),
		)
		return
	}

	// Se valida que el campo "Comandos" no sea nulo ni esté vacío (después de quitar espacios en blanco).
	if requestBody.Comandos == nil || strings.TrimSpace(*requestBody.Comandos) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			general.ResultadoSalida("El campo 'Comandos' es obligatorio y no puede ser nulo", true, nil),
		)
		return
	}

	/** =======================================================
	 		3. PROCESAMIENTO Y EJECUCIÓN DE COMANDOS
		=======================================================	*/

	// Se extrae el string con todos los comandos.
	var ejecutar []string
	ejecutar = append(ejecutar, *requestBody.Comandos)

	// Se divide el string de comandos en un slice de strings, donde cada elemento es una línea.
	comando := strings.Split(ejecutar[0], "\n")

	// Se procesa la lista inicial de comandos para filtrar comentarios, espacios en blanco
	// y obtener una lista limpia de comandos listos para ser ejecutados.
	tempComandos := general.ExecuteCommandList(comando)

	// Se extrae la estructura de salida del resultado del pre-procesamiento.
	salida, ok := tempComandos.Respuesta.(general.SalidaComandoEjecutado)
	if !ok {
		// Este error ocurriría si el tipo de `Respuesta` no es el esperado, indicando un problema interno.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // Error 500 (Error Interno del Servidor).
		json.NewEncoder(w).Encode(
			general.ResultadoSalida("Error interno al procesar la lista de comandos", true, nil),
		)
		return
	}

	// Debug (Opcional): Imprimir en la consola del servidor los comandos que se van a ejecutar.
	fmt.Println("--- Comandos a Ejecutar ---")
	for _, temp := range salida.LstComandos {
		println(temp)
	}
	fmt.Println("---------------------------")

	// Se llama a la función principal que ejecuta cada comando y retorna los resultados.
	errores, mensajes, contadorErrores := general.GlobalCom(salida.LstComandos)
	fmt.Printf("Resultados de Ejecución -> Errores: %d, Mensajes: %d\n", contadorErrores, len(mensajes))

	// Se actualiza la estructura `salida` con los mensajes y errores generados durante la ejecución.
	salida.Mensajes = mensajes
	salida.Errores = errores

	/** =======================================================
		4. RESPUESTA AL FRONTEND
	=======================================================	*/

	// Se envía la respuesta final al frontend, incluyendo la lista de comandos procesados
	// y los resultados de su ejecución.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(general.ResultadoSalida("", false, salida))
	if err != nil {
		// Si ocurre un error al codificar la respuesta JSON final, se registra en el servidor.
		// En este punto, es difícil enviar una respuesta de error al cliente si el `Encode` ya falló.
		fmt.Printf("Error al codificar la respuesta JSON final: %v\n", err)
	}
}
