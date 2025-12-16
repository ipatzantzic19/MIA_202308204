package global

import "Proyecto/Estructuras/structures"

// Estructura que representa una particion montada
type ParticionesMontadas struct {
	Path           string
	DriveLetter    byte
	ID_Particion   [4]byte
	Type           byte
	Es_Particion_P bool
	Es_Particion_L bool
	Particion_P    structures.Partition
	Particion_L    structures.EBR
}

// Particion montada por defecto
var ParticionMontadaDefault = ParticionesMontadas{
	Path:           "",
	DriveLetter:    '0',
	ID_Particion:   Global_ID(""),
	Type:           '0',
	Es_Particion_P: false,
	Es_Particion_L: false,
	Particion_P:    structures.Partition{},
	Particion_L:    structures.EBR{},
}

// Lista de particiones montadas
var Mounted_Partitions []ParticionesMontadas

// Estructura que representa un grupo de usuarios
type Grupo struct {
	GID    int32
	Tipo   byte
	Nombre string
}

// Estructura que representa un usuario
type Usuario struct {
	UID          int32
	GID          int32
	Tipo         byte
	Grupo        [10]byte
	User         [10]byte
	Password     [10]byte
	ID_Particion [4]byte
	Logged_in    bool
	Mounted      ParticionesMontadas
}

// Usuario root por defecto
var UsuarioLogeado = Usuario{
	UID:          1,
	GID:          1,
	Tipo:         'U',
	Grupo:        Global_Data("root"),
	User:         Global_Data("root"),
	Password:     Global_Data("123"),
	ID_Particion: Global_ID(""),
	Logged_in:    false,
}

// Usuario por defecto (no logeado)
var DefaultUser = Usuario{
	UID:          -1,
	GID:          -1,
	Tipo:         '0',
	Grupo:        Global_Data(""),
	User:         Global_Data(""),
	Password:     Global_Data(""),
	ID_Particion: Global_ID(""),
	Logged_in:    false,
	Mounted:      ParticionMontadaDefault,
}

// Grupo root por defecto
var GrupoUsuarioLoggeado = Grupo{
	GID:    1,
	Tipo:   'G',
	Nombre: "root",
}

// Grupo por defecto (no logeado)
var DefaultGrupoUsuario = Grupo{
	GID:    -1,
	Tipo:   '0',
	Nombre: "",
}
