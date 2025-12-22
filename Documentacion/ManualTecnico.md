# Manual Técnico — Sistema de Archivos EXT2 Simulado

Fecha: 21 de diciembre de 2025

## Introducción

Este Manual Técnico documenta la implementación interna y el uso del sistema de archivos EXT2 simulado en esta aplicación web. El proyecto está dividido en frontend (interfaz) y backend (lógica y manipulación de discos en archivos `.mia`). El propósito es describir la arquitectura, las estructuras de datos principales, los comandos implementados y ejemplos de uso.

## Estructura del repositorio (relevante)

- Backend: [Backend](Backend)
	- Servidor y controladores: [Backend/server.go](Backend/server.go)
	- Comandos implementados: [Backend/comandos](Backend/comandos)
	- Utilidades y estructuras: [Backend/Estructuras](Backend/Estructuras) y [Backend/comandos/utils](Backend/comandos/utils)
	- Reportes: [Backend/comandos/report](Backend/comandos/report)
- Frontend: [Frontend](Frontend) — UI en React con componentes en [Frontend/src/components](Frontend/src/components) que permiten ejecutar comandos y mostrar resultados.

## 1. Arquitectura del sistema

Visión general:

- El frontend (aplicación web) captura comandos y parámetros del usuario y los envía al backend mediante peticiones HTTP (API local).
- El backend recibe las solicitudes, ejecuta la lógica correspondiente (crear disco, particionar, montar, manipular archivos, etc.) y devuelve resultados o mensajes.
- Los discos simulados son archivos binarios con extensión `.mia` ubicados en `tmp/VDIC-MIA/Disks`.

Flujo de interacción (simplificado):

1. Usuario ingresa un comando en la UI.
2. Frontend envía la solicitud al controlador del backend (ej. [Backend/comandos/controllers/ComandosController.go](Backend/comandos/controllers/ComandosController.go)).
3. El backend ejecuta la función de comando (ubicada en `Backend/comandos/*`) y actualiza el archivo `.mia` o las estructuras internas.
4. Si se solicita, se generan reportes (Graphviz) en `VDIC-MIA/Rep` y se envían como imagen al frontend.

Diagrama conceptual (texto):

Frontend (React) <--HTTP--> Backend (Go)
	- Frontend: componentes `PanelConsola`, `PanelEditor`, `InfoCards` (ver [Frontend/src/components](Frontend/src/components)).
	- Backend: módulos para administración de discos, archivos y reportes (ver [Backend/comandos](Backend/comandos)).

## 2. Estructuras de datos principales

Las estructuras de datos clave están en `Backend/Estructuras` y `Backend/comandos/utils`. A continuación se describen y su función.

### 2.1 MBR (Master Boot Record)
- Ubicación: comienzo del archivo `.mia`.
- Contenido clave: tamaño total del disco (`Mbr_tamano`), fecha de creación, firma del disco, tipo de ajuste y un arreglo de hasta 4 particiones (`Mbr_partitions`).
- Efecto: define la tabla primaria de particiones (primarias y una posible extendida). La estructura se encuentra como `MBR` en [Backend/Estructuras/structures/disk.go](Backend/Estructuras/structures/disk.go).

Campos relevantes (resumen):
- `Mbr_tamano` (int32): bytes totales del disco.
- `Mbr_partitions` ([4]Partition): array con 4 entradas de particiones.

### 2.2 Partition (entrada de partición)
- Cada entrada contiene: estado (activa/inactiva), tipo (`P` primaria o `E` extendida), fit (`B/F/W`), `Part_start`, `Part_s` (tamaño), `Part_name`.
- Se usa para ubicar particiones dentro del disco y para crear nuevas entradas vía `FDISK`.

### 2.3 EBR (Extended Boot Record)
- Las particiones lógicas dentro de una partición extendida se representan con EBRs encadenados.
- Estructura `EBR` está en [Backend/Estructuras/structures/disk.go](Backend/Estructuras/structures/disk.go).
- Campos: `Part_start`, `Part_s`, `Part_next` (offset al siguiente EBR o -1), `Name`.

### 2.4 Inodos y Bloques (sistema de archivos simulado)
- El proyecto implementa componentes de un sistema tipo EXT2: inodos e bloques. Las estructuras y utilidades relacionadas están en [Backend/Estructuras/size](Backend/Estructuras/size) y [Backend/comandos/utils](Backend/comandos/utils).
- Inodo (inode): guarda permisos, propietario, timestamps, punteros a bloques y tamaño del archivo.
- Bloque (block): unidades de almacenamiento para contenido de archivos y metadata (bitmap, journal, etc.).

Nota: para ver las definiciones concretas, revisar:
- [Backend/Estructuras/structures/files.go](Backend/Estructuras/structures/files.go)
- [Backend/Estructuras/size/blocks.go](Backend/Estructuras/size/blocks.go)

## 3. Comandos implementados

Los comandos están en subcarpetas dentro de `Backend/comandos` agrupados por responsabilidad.

- Administración de Discos: [Backend/comandos/admindisk](Backend/comandos/admindisk)
	- `mkdisk.go` — MKDISK: crea un archivo `.mia` de tamaño dado.
		- Parámetros típicos: `-size=N` (en bytes o sufijos), `-path=RUTA`, `-fit=[BFW]`.
		- Ejemplo: MKDISK -size=10M -path=VDIC-MIA/Disks/Disco1.mia -fit=F
	- `rmdisk.go` — RMDISK: elimina archivo de disco.
	- `mkfs.go` — MKFS: (si aplica) formatea el disco para crear estructuras de sistema de archivos.

- Administración de Particiones: [Backend/comandos/particion-fdisk](Backend/comandos/particion-fdisk)
	- `fdisk.go` — FDISK: crear, eliminar y modificar particiones.
		- Sintaxis: FDISK -size=N -path=DISCO -type=[P|E|L] -name=NAME -unit=[b|k|m] -delete=[fast|full]
		- Efecto: actualiza el MBR (agrega entrada primaria o extendida) y/o crea EBRs para lógicas.

- Administración de Archivos/Directorios: [Backend/comandos/adminfiles](Backend/comandos/adminfiles)
	- `mkfiles.go` — MKFILE: crear archivos dentro del sistema de archivos simulado.
	- `mkdir.go` — MKDIR: crear carpetas.
	- `cat.go` — CAT: mostrar contenido.

- Usuarios y Permisos: [Backend/comandos/admingroup](Backend/comandos/admingroup) y [Backend/comandos/adminuser](Backend/comandos/adminuser)
	- `mkusr.go`, `mkgrp.go`, `login.go`, `logout.go` — gestión de usuarios y sesiones.

- Montaje: [Backend/comandos/admindisk/mount.go] y utilidades en [Backend/comandos/utils/mount.go]
	- `mount` — registra particiones montadas con un id para operaciones posteriores.
	- `unmount` — libera el montaje.

- Reports: [Backend/comandos/report](Backend/comandos/report)
	- `disk_report.go` — genera representación gráfica del disco (MBR, particiones, EBRs, espacios libres) usando Graphviz. Genera un temporal `.dot` y lo convierte a imagen (PNG, JPG, etc.).
	- `mbr.go` — genera reporte detallado del MBR y partición seleccionada.

## 4. Ejemplos de uso

- Crear disco (MKDISK):

	MKDISK -size=10M -path=VDIC-MIA/Disks/Disco1.mia -fit=F

- Crear partición primaria (FDISK):

	FDISK -size=2M -path=VDIC-MIA/Disks/Disco1.mia -type=P -name=Part1 -unit=M

- Crear partición extendida y lógicas (FDISK):

	FDISK -size=5M -path=VDIC-MIA/Disks/Disco1.mia -type=E -name=Ext1 -unit=M
	FDISK -size=1M -path=VDIC-MIA/Disks/Disco1.mia -type=L -name=Log1 -unit=M

- Montar partición:

	MOUNT -path=VDIC-MIA/Disks/Disco1.mia -name=Part1

- Generar reporte DISK (ejemplo):

	REPORT DISK -id=vd1 -path=VDIC-MIA/Rep/disk.png

	Nota: Graphviz (`dot`) debe estar instalado y en el PATH para que el backend pueda convertir el `.dot` a imagen.

## 5. Cálculo de porcentajes y representaciones gráficas

- El reporte DISK construye una tabla HTML embebida dentro de DOT para mostrar secciones del disco. El algoritmo considera:
	- MBR (visual) y el resto del disco.
	- Orden físico de particiones por `Part_start`.
	- Espacios libres antes/después/entre particiones.
	- Partición extendida contiene EBR/lógicas.
- Para asegurar que la sumatoria de porcentajes sea 100%, el reporte normaliza valores sobre la base del disco y ajusta pequeñas diferencias por redondeo en el último segmento.

## 6. Requisitos y ejecución

- Go (para compilar/ejecutar el backend).
- Node.js + npm (para el frontend si quieres ejecutar la UI localmente).
- Graphviz (`dot`) instalado y disponible en PATH para generar imágenes.

Comandos rápidos para probar localmente (en carpeta `Backend`):

```powershell
go run server.go
```

Generar un reporte de ejemplo (desde UI o llamada al endpoint correspondiente): especificar `ruta_reporte` con extensión válida (`.png`, `.jpg`).

## 7. Troubleshooting conocido

- Error: "Error al ejecutar 'dot'... syntax error": usualmente causado por etiquetas HTML/ASCII no escapadas en nombres de particiones. Se corrigió mediante `escapeHTML` en `Backend/comandos/report/disk_report.go`.
- Error: imagen vacía o .dot no convertido: verificar que `dot` esté instalado y que la ruta de salida tenga permisos de escritura.

## 8. Archivos clave para revisar

- Backend entry: [Backend/server.go](Backend/server.go)
- Comandos: [Backend/comandos](Backend/comandos)
- Reportes: [Backend/comandos/report/disk_report.go](Backend/comandos/report/disk_report.go) y [Backend/comandos/report/mbr.go](Backend/comandos/report/mbr.go)
- Estructuras: [Backend/Estructuras/structures/disk.go](Backend/Estructuras/structures/disk.go) y [Backend/Estructuras/structures/files.go](Backend/Estructuras/structures/files.go)
- Frontend: [Frontend/src/components/PanelConsola.jsx](Frontend/src/components/PanelConsola.jsx)

---

Si quieres, puedo:
- Incluir diagramas en formato ASCII más detallados o generar imágenes DOT que documenten la arquitectura.
- Añadir una sección con tests/ejemplos reproducibles paso a paso.
- Generar un índice y cross-referencias internas más extensas.

Indícame qué prefieres que añada o si deseas que ejecute pruebas y genere un reporte DISK de ejemplo en `VDIC-MIA/Rep/disk.png`.

