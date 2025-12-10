# 📚 ARQUITECTURA DEL BACKEND - GoDisk

## Visión General

El Backend es un servidor HTTP escrito en Go que gestiona un sistema de archivos virtual tipo EXT2. Recibe comandos desde el Frontend (React), los procesa y devuelve respuestas en JSON.

---

## 🏗️ Estructura del Proyecto

```
Backend/
├── server.go                    # Punto de entrada principal
├── go.mod, go.sum             # Dependencias del proyecto
├── air.toml                    # Configuración de hot-reload (air)
├── comandos.txt                # Comandos de instalación
├── Estructuras/                # Modelos de datos (structs) - no usado aún
├── VDIC-MIA/                   # Carpeta de datos (se crea automáticamente)
│   ├── Disks/                  # Discos virtuales (.mia)
│   ├── Rep/                    # Reportes generados
│   └── CarpetaImagenes.txt     # Metadata
├── comandos/
│   ├── general/
│   │   ├── general.go          # Funciones helper
│   │   └── ExecuteCommands.go  # Router de comandos
│   └── controllers/
│       └── ComandosController.go  # Maneja requests HTTP
└── tmp/                        # Carpeta temporal (air watch)
```

---

## 🔄 Flujo de Datos: Frontend → Backend

```
┌─────────────┐
│  Frontend   │
│  (React)    │
└──────┬──────┘
       │
       │ POST /commands
       │ {"Comandos": "mkdisk -size=3000 -unit=K\nfdisk -size=300 ..."}
       │
       ▼
┌──────────────────────────────────────────┐
│         server.go (main.go)              │
│  - Escucha en puerto 9700                │
│  - Configura CORS                        │
│  - Registra rutas                        │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│  ComandosController.go (HandleCommand)   │
│  - Valida método POST                    │
│  - Decodifica JSON                       │
│  - Extrae campo "Comandos"               │
│  - Divide en líneas (\n)                 │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│  general.go (ExecuteCommandList)         │
│  - Filtra líneas vacías                  │
│  - Elimina comentarios (#)               │
│  - Limpia espacios                       │
│  - Devuelve array de comandos válidos    │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│  general.go (GlobalCom - COMENTADO)      │
│  - Identifica tipo de comando            │
│  - Dirigiría a módulo correspondiente    │
│  - Módulos aún no implementados          │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│  Respuesta JSON devuelta al Frontend     │
│  {"Error": false,                        │
│   "Respuesta": [...comandos procesados]}│
└──────────────────────────────────────────┘
       │
       ▼
┌─────────────┐
│  Frontend   │
│  Actualiza  │
│     UI      │
└─────────────┘
```

---

## 📄 Archivos Principales Explicados

### 1️⃣ **server.go** - Punto de Entrada
```go
// Puerto donde escucha
puerto := 9700

// CORS: Permite requests desde cualquier origen
c := cors.AllowAll()

// Router HTTP
mux.HandleFunc("/commands", controllers.HandleCommand)

// Escucha continuamente
http.ListenAndServe(":9700", handler)
```

**¿Por qué 9700?**
- Es arbitrario, solo debe ser un puerto disponible
- No conflictúa con Vite (5173), MySQL (3306), MongoDB (27017), etc.

**¿Qué es CORS?**
- Frontend corre en http://localhost:5173 (o 5174)
- Backend corre en http://localhost:9700
- Son "orígenes diferentes", el navegador lo bloquea
- CORS permite que el Frontend acceda al Backend

---

### 2️⃣ **ComandosController.go** - Maneja Requests HTTP

```
POST http://localhost:9700/commands
Content-Type: application/json

{
  "Comandos": "mkdisk -size=3000 -unit=K\nfdisk -size=300 ..."
}
```

**Pasos:**
1. Valida que sea POST (rechaza GET, PUT, etc.)
2. Decodifica JSON con `json.Decoder`
3. Valida que "Comandos" no sea vacío
4. Divide por `\n` (saltos de línea)
5. Llama a `ExecuteCommandList()`
6. Devuelve JSON con resultado

**Códigos HTTP:**
- `200 OK`: Éxito
- `400 Bad Request`: JSON inválido, campo vacío
- `405 Method Not Allowed`: No es POST

---

### 3️⃣ **ExecuteCommandList()** - Procesa Comandos

Ubicado en `general.go`

**¿Qué hace?**
- Recibe array de strings (líneas de comandos)
- Filtra líneas vacías
- Elimina comentarios (líneas que empiezan con `#`)
- Limpia espacios en blanco
- Usa regex para parsear correctamente
- Devuelve array de comandos válidos

**Ejemplo:**
```
Input:
[
  "# Esto es un comentario",
  "mkdisk -size=3000 -unit=K",
  "",
  "fdisk -size=300 -diskName=VDIC-A.mia -name=Particion1",
  "   # Otro comentario"
]

Output:
[
  "mkdisk -size=3000 -unit=K",
  "fdisk -size=300 -diskName=VDIC-A.mia -name=Particion1"
]
```

---

### 4️⃣ **GlobalCom()** - Router de Comandos (COMENTADO)

Ubicado en `ExecuteCommands.go`

**¿Qué debería hacer?**
- Identificar el tipo de comando por su prefijo
- Dirigirlo al módulo correspondiente

**Tipos de comandos:**
```
Comando          Módulo               Estado
────────────────────────────────────────────────
mkdisk, fdisk    admindisk.go         ❌ Comentado
rep              report.go            ❌ Comentado
mkfile, mkdir    filesystem.go        ❌ Comentado
cat              filesystem.go        ❌ Comentado
login, logout    adminusers.go        ❌ Comentado
mkgrp, mkusr     adminusers.go        ❌ Comentado
```

**¿Por qué está comentado?**
- Los módulos no existen aún
- Se implementarán en fases posteriores
- Por ahora el Backend solo procesa los comandos pero no los ejecuta

---

## 🔧 Funciones Helper

### `ObtenerComandos(comando string) []string`
Extrae los parámetros de un comando usando regex.

```go
// Input: "mkdisk -size=3000 -unit=K -fit=BF"
// Output: ["-size=3000", "-unit=K", "-fit=BF"]

// Input: 'mkfile -path="/archivos/test" -name="datos"'
// Output: ["-path=/archivos/test", "-name=datos"]
```

### `getCommand(comm string, commands ...string) string`
Identifica cuál comando es (entre varias opciones).

```go
// Input: "mkdisk -size=3000", ["mkdisk", "fdisk", "rmdisk"]
// Output: "mkdisk"

// Input: "rmdisk -diskName=VDIC-A", ["mkdisk", "fdisk", "rmdisk"]
// Output: "rmdisk"
```

### `CrearCarpeta()`
Crea las carpetas necesarias si no existen.

```
VDIC-MIA/
├── Disks/          # Aquí van los discos (.mia)
├── Rep/            # Aquí van los reportes
└── CarpetaImagenes.txt  # Metadata
```

---

## 📦 Dependencias

```go
go.mod:
- github.com/fatih/color  v1.18.0  // Colores en consola (rojo, verde, amarillo)
- github.com/rs/cors      v1.11.1  // Middleware CORS

go.mod indirectas (se instalan solas):
- golang.org/x/sys        // Funciones del sistema
- github.com/mattn/go-*   // Helper para colores
```

---

## 🚀 Cómo Funciona (Ejemplo Real)

### Paso 1: Frontend envía comandos
```javascript
// Frontend (React)
const response = await fetch('http://localhost:9700/commands', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ 
    Comandos: "mkdisk -size=3000 -unit=K\nfdisk -size=300 ..."
  })
})
```

### Paso 2: Backend recibe en HandleCommand()
```
POST /commands
Body: {"Comandos": "mkdisk ...\nfdisk ..."}
```

### Paso 3: Descodifica JSON y divide por \n
```go
comando := strings.Split(ejecutar[0], "\n")
// ["mkdisk -size=3000 -unit=K", "fdisk -size=300 ...", ...]
```

### Paso 4: ExecuteCommandList() procesa
```go
// Filtra vacías, comentarios, limpia espacios
// Input:  ["mkdisk -size=3000 -unit=K", "", "# comentario", "fdisk ..."]
// Output: ["mkdisk -size=3000 -unit=K", "fdisk ..."]
```

### Paso 5: Devuelve JSON
```json
{
  "Error": false,
  "Respuesta": [
    "mkdisk -size=3000 -unit=K",
    "fdisk -size=300 -diskName=VDIC-A.mia -name=Particion1"
  ]
}
```

### Paso 6: Frontend recibe y actualiza UI
```javascript
const data = await response.json()
console.log(data.Respuesta) // ["mkdisk ...", "fdisk ..."]
```

---

## ❌ Estado Actual (MVP - Producto Mínimo Viable)

### Implementado ✅
- [x] Servidor HTTP en puerto 9700
- [x] CORS habilitado
- [x] Endpoint POST /commands
- [x] Validación de JSON
- [x] Parsing de comandos (divide por \n)
- [x] Filtrado de comentarios y líneas vacías
- [x] Creación automática de carpetas (VDIC-MIA)

### No Implementado ❌
- [ ] Procesamiento real de comandos (admindisk.go, filesystem.go, etc.)
- [ ] Creación de discos virtuales
- [ ] Particionar discos
- [ ] Crear archivos/carpetas
- [ ] Sistema de usuarios/permisos
- [ ] Reportes
- [ ] Autenticación/Login

---

## 🔐 Seguridad (Notas)

### ⚠️ CORS demasiado abierto
```go
c := cors.AllowAll()  // ❌ Permite desde CUALQUIER sitio
```

**Debería ser:**
```go
c := cors.New(cors.Options{
  AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
  AllowedMethods: []string{"POST"},
  AllowedHeaders: []string{"Content-Type"},
})
```

### ⚠️ Sin autenticación
- Cualquiera puede enviar comandos al backend
- Debería verificar JWT token o sesión

### ⚠️ Sin validación de parámetros
- Acepta cualquier comando
- Debería validar que los parámetros sean válidos

---

## 📋 TODO (Fases Futuras)

### Fase 2: Módulos de Comandos
- [ ] Implementar `admindisk.go` (mkdisk, fdisk, mount, etc.)
- [ ] Implementar `filesystem.go` (mkfile, mkdir, cat, etc.)
- [ ] Implementar `adminusers.go` (login, logout, mkgrp, mkusr, etc.)
- [ ] Implementar `report.go` (generación de reportes)

### Fase 3: Persistencia
- [ ] Guardar discos virtuales en archivos .mia
- [ ] Guardar estructura de archivos
- [ ] Guardar usuarios/permisos

### Fase 4: Autenticación & Seguridad
- [ ] Implementar login real
- [ ] Generar JWT tokens
- [ ] Validar tokens en cada request
- [ ] Restringir CORS

### Fase 5: API Adicional
- [ ] GET /disks - Listar discos
- [ ] GET /partitions - Listar particiones
- [ ] GET /files - Listar archivos
- [ ] GET /reports - Descargar reportes

---

## 🤔 Preguntas Frecuentes

**P: ¿Por qué Go?**
A: Es rápido, compila a un binario único (sin dependencias externas), y maneja concurrencia bien.

**P: ¿Por qué puerto 9700?**
A: Arbitrario. Cualquier puerto libre funciona.

**P: ¿Por qué CORS?**
A: Frontend y Backend corren en puertos diferentes. CORS permite que se comuniquen.

**P: ¿Por qué los módulos están comentados?**
A: Porque no existen aún. Se implementarán en paralelo.

**P: ¿Cómo agrego un nuevo comando?**
A: 
1. Añade la condición en `GlobalCom()` en ExecuteCommands.go
2. Crea el módulo que lo implemente (ej: `disks.go`)
3. Llámalo desde `GlobalCom()`

---

## 🔗 Referencias

- [Documentación Go](https://golang.org/doc)
- [Paquete net/http](https://pkg.go.dev/net/http)
- [CORS en Go](https://github.com/rs/cors)
- [Expressions regulares en Go](https://golang.org/pkg/regexp)

