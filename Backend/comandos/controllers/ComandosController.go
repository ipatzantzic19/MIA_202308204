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

HandleCommand es el endpoint HTTP que recibe comandos desde el Frontend.

RESPONSABILIDADES:
 1. Validar que el request sea POST
 2. Validar CORS
 3. Extraer los comandos del JSON enviado
 4. Dividir los comandos en líneas individuales
 5. Procesar cada comando filtrando comentarios y espacios en blanco
 6. Devolver la lista de comandos procesados como JSON

FLUJO:

	Frontend                   Backend
	   |                         |
	   |---- POST /commands ----->|
	   |   (JSON: {"Comandos": "mkdisk ...\nfdisk ..."})
	   |                         |
	   |                    HandleCommand()
	   |                    - Valida método POST
	   |                    - Valida JSON
	   |                    - Divide en líneas
	   |                    - Filtra comentarios
	   |                    - Procesa comando
	   |                         |
	   |<--- JSON Response -------|
	   |   ({"Respuesta": [...], "Error": false})
	   |
	Actualiza UI

=======================================================
*/
func HandleCommand(w http.ResponseWriter, r *http.Request) {
	/** =======================================================
		VALIDACIONES INICIALES DE CORS Y MÉTODO HTTP
	=======================================================	*/

	// Manejar CORS: para requests OPTIONS (preflight de CORS)
	// El navegador envía OPTIONS antes de POST si es cross-origin
	if r.Method == http.MethodOptions {
		// Establecer encabezados CORS para las solicitudes OPTIONS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo permitir solicitudes POST
	// Rechaza GET, PUT, DELETE, etc.
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar encabezados CORS para las solicitudes POST
	// Permite que el Frontend (en otro puerto) acceda a esta ruta
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	/** =======================================================
		EXTRACCIÓN Y VALIDACIÓN DEL CUERPO JSON
	=======================================================	*/

	// Decodificar el cuerpo JSON de la solicitud
	// El Frontend envía: {"Comandos": "mkdisk -size=3000 ...\nfdisk ..."}

	// Estructura para mapear el JSON entrante
	var requestBody struct {
		Comandos *string `json:"Comandos"` // Puntero (*) para poder detectar si es nulo
	}

	decoder := json.NewDecoder(r.Body) // Crear un decodificador JSON para el cuerpo (los datos recibidos en back)
	decoder.DisallowUnknownFields()    // Rechaza campos desconocidos para el struct

	err := decoder.Decode(&requestBody) // Decodifica el JSON de requestBody, (& es para pasar la dirección del puntero)

	// Manejar errores de decodificación JSON
	if err != nil {
		// Si el JSON es inválido, responder con error
		w.Header().Set("Content-Type", "application/json") // Respuesta JSON
		w.WriteHeader(http.StatusBadRequest)               // Código 400 Bad Request
		json.NewEncoder(w).Encode(                         // Responder con mensaje de error
			general.ResultadoSalida("JSON inválido o campos no permitidos", true, nil),
		)
		return
	}

	// Validar que el campo "Comandos" no esté vacío
	if requestBody.Comandos == nil || strings.TrimSpace(*requestBody.Comandos) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			general.ResultadoSalida("El campo 'Comandos' es obligatorio y no puede ser nulo", true, nil),
		)
		return
	}

	// Preparar los comandos para procesar
	var ejecutar []string
	ejecutar = append(ejecutar, *requestBody.Comandos) // Obtener el string de comandos

	/** =======================================================
	 		Dividir los comandos por saltos de línea
		=======================================================	*/
	comando := strings.Split(ejecutar[0], "\n") // Divide en líneas

	// Procesar la lista de comandos
	// ExecuteCommandList filtra comentarios, espacios, valida formato
	// y devuelve solo los comandos válidos
	tempComandos := general.ExecuteCommandList(comando)

	// Extraer la respuesta del resultado procesado
	salida, ok := tempComandos.Respuesta.(general.SalidaComandoEjecutado)
	if !ok {
		// Si hay error en el procesamiento
		// http.Error(w, "Error al obtener comandos", http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			general.ResultadoSalida("Error al obtener comandos", true, nil),
		)
		return
	}

	// Debug: Imprimir en consola los comandos procesados (opcional)
	for _, temp := range salida.LstComandos {
		println(temp) // Imprime: "mkdisk -size=3000 ...", etc.
	}

	errores, mensajes, contadorErrorres := general.GlobalCom(salida.LstComandos)
	fmt.Println("Errores:", errores, "Mensajes:", mensajes, "Total Errores:", contadorErrorres)

	// Actualizar la salida con mensajes y errores de ejecución
	salida.Mensajes = mensajes
	salida.Errores = errores

	// Responder al Frontend con la lista de comandos, mensajes y errores de ejecución
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(general.ResultadoSalida("", false, salida))
	if err != nil {
		// Error al codificar la respuesta JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			general.ResultadoSalida("Error al leer el cuerpo de la solicitud", true, nil),
		)
		return
	}
}
