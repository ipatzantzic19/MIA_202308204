/*
========== SERVIDOR BACKEND - GoDisk (Sistema de Archivos EXT2) ==========

DESCRIPCIÓN GENERAL:
  Este es el servidor HTTP principal que recibe comandos desde el Frontend,
  los procesa y devuelve respuestas JSON. Es el corazón de la aplicación backend.

RESPONSABILIDADES:
  1. Inicializar servidor HTTP en puerto 9700
  2. Configurar CORS para permitir requests desde Frontend (puerto 5173/5174)
  3. Registrar rutas HTTP (endpoints que el Frontend puede llamar)
  4. Crear carpetas necesarias (VDIC-MIA, Reportes, Discos)
  5. Escuchar y procesar requests POST con comandos del filesystem

FLUJO DE DATOS:
  1. Frontend envía comando POST a http://localhost:9700/commands
  2. Middleware CORS valida que la petición sea permitida (cross-origin)
  3. Router HTTP dirige la solicitud al controlador HandleCommand
  4. HandleCommand extrae los comandos del JSON
  5. ExecuteCommandList valida y parsea los comandos
  6. Se devuelve JSON con resultado al Frontend

ESTRUCTURA DEL PROYECTO:
  - server.go: Este archivo (punto de entrada, configuración HTTP)
  - comandos/controllers/ComandosController.go: Maneja requests HTTP
  - comandos/general/ExecuteCommands.go: Router de comandos internos
  - comandos/general/general.go: Funciones helper
  - Estructuras/: Modelos de datos (structs)
  - VDIC-MIA/: Carpeta donde se guardan discos y reportes (se crea automáticamente)

==========================================================================
*/

package main

import (
	// Paquetes internos del proyecto
	"Proyecto/comandos/controllers" // Controladores que procesan requests HTTP
	"Proyecto/comandos/general"     // Funciones generales y helpers
	"fmt"                           // Para imprimir en consola
	"net/http"                      // Servidor HTTP estándar de Go

	"github.com/rs/cors" // Middleware CORS: permite requests desde dominio diferente
)

func main() {
	// NewServeMux crea un nuevo router HTTP para manejar diferentes rutas
	// Un ServeMux dirige requests HTTP a diferentes handlers(↓) según la URL
	// handlers son funciones que procesan requests HTTP
	mux := http.NewServeMux()

	// Puerto en el que escucha el servidor
	puerto := 9700

	// Configuracion de CORS
	// cors.AllowAll() permite requests desde CUALQUIER origen (dominio)
	// IMPORTANTE: En producción,  restringir esto a dominios específicos por seguridad
	c := cors.AllowAll()

	// Cuando el Frontend hace POST a /commands, se ejecuta HandleCommand
	// El controlador HandleCommand está definido en ComandosController.go
	// HandleCommand procesa los comandos recibidos
	mux.HandleFunc("/commands", controllers.HandleCommand)

	// Rutas futuras comentadas (para expandir la API más adelante):
	// mux.HandleFunc("/login", handleLogin)           // Autenticación
	// mux.HandleFunc("/logout", handleLogout)         // Cierre de sesión
	// mux.HandleFunc("/obtainmbr", handleObtainMBR)   // Obtener MBR del disco
	// mux.HandleFunc("/reportesobtener", handleReportsObtener) // Reportes
	// mux.HandleFunc("/graphs", handleGraph)          // Gráficos
	// mux.HandleFunc("/obtain-carpetas-archivos", handleObtainCarpetasArchivos) // Listar archivos
	// mux.HandleFunc("/cat", handleCat)               // Ver contenido de archivo

	// Aplicar el middleware CORS al router
	// Esto envuelve todas las rutas con validación de CORS
	handler := c.Handler(mux)

	// Imprimir en consola que el servidor está corriendo
	fmt.Println("" + fmt.Sprintf("Backend server is on %v", puerto))

	// Crear las carpetas necesarias (VDIC-MIA, Rep, Disks)
	// Si ya existen, las saltea; si no, las crea
	general.CrearCarpeta()

	// Carga las particiones que ya estaban montadas desde los archivos de disco.
	general.LoadMountedPartitions()

	// Funciones futuras comentadas (para funcionalidades posteriores):
	// obtencionpf.ObtenerMBR_Mounted()         // Obtener MBR
	// obtencionpf.MostrarParticionesMontadas()  // Mostrar particiones

	// Iniciar el servidor HTTP
	// ListenAndServe bloquea el programa y escucha requests en el puerto
	// Devuelve error si algo falla (puerto en uso, falta de permisos, etc.)
	err := http.ListenAndServe(":"+fmt.Sprintf("%v", puerto), handler)
	if err != nil {
		fmt.Println("ERROR al iniciar servidor:", err)
	}
}
