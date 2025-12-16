package global

//  convierte una cadena de texto en un arreglo de bytes de tamaño 10
func Global_Data(value string) [10]byte {
	temp := make([]byte, 10)
	for i := range temp {
		temp[i] = '\x00'
	}
	copy(temp[:], []byte(value))
	return [10]byte(temp)
}

// convierte una cadena de texto en un arreglo de bytes de tamaño 4
func Global_ID(value string) [4]byte {
	temp := make([]byte, 4)
	for i := range temp {
		temp[i] = '\x00'
	}
	copy(temp[:], []byte(value))
	return [4]byte(temp)
}
