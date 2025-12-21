package report

import (
	"Proyecto/Estructuras/structures"
	"Proyecto/comandos/utils"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Report_MBR(id_particion string, ruta_reporte string) (string, error) {
	particion_montada, encontrada := utils.ObtenerDiscoID(id_particion)
	if !encontrada {
		return "", fmt.Errorf("[REP-MBR]: No se encontró una partición montada con el id '%s'", id_particion)
	}

	diskPath := particion_montada.Path
	mbr, success := utils.Obtener_FULL_MBR_FDISK(diskPath)
	if !success {
		return "", fmt.Errorf("[REP-MBR]: No se pudo leer el MBR del disco en '%s'", diskPath)
	}

	var dotBuffer bytes.Buffer
	dotBuffer.WriteString("digraph G {\n")
	dotBuffer.WriteString("\tnode[shape=none];\n")
	dotBuffer.WriteString("\tstart[label=<\n")
	dotBuffer.WriteString("\t<table border='0' cellborder='1' cellspacing='0'>\n")
	dotBuffer.WriteString(`<tr><td colspan="2" bgcolor="#007BFF"><font color="white"><b>REPORTE DE MBR</b></font></td></tr>` + "\n")
	dotBuffer.WriteString(fmt.Sprintf(`<tr><td><b>mbr_tamano</b></td><td>%d</td></tr>`+"\n", mbr.Mbr_tamano))
	dotBuffer.WriteString(fmt.Sprintf(`<tr><td><b>mbr_fecha_creacion</b></td><td>%s</td></tr>`+"\n", utils.IntFechaToStr(mbr.Mbr_fecha_creacion)))
	dotBuffer.WriteString(fmt.Sprintf(`<tr><td><b>mbr_disk_signature</b></td><td>%d</td></tr>`+"\n", mbr.Mbr_disk_signature))
	dotBuffer.WriteString(fmt.Sprintf(`<tr><td><b>dsk_fit</b></td><td>%c</td></tr>`+"\n", mbr.Dsk_fit))

	for i, p := range mbr.Mbr_partitions {
		if p.Part_status != '0' {
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td colspan="2" bgcolor="#17A2B8"><font color="white"><b>Partición %d</b></font></td></tr>`+"\n", i+1))
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_status</td><td>%c</td></tr>`+"\n", p.Part_status))
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_type</td><td>%c</td></tr>`+"\n", p.Part_type))
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_fit</td><td>%c</td></tr>`+"\n", p.Part_fit))
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_start</td><td>%d</td></tr>`+"\n", p.Part_start))
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_s</td><td>%d</td></tr>`+"\n", p.Part_s))
			dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_name</td><td>%s</td></tr>`+"\n", utils.ToString(p.Part_name[:])))

			if p.Part_type == 'E' {
				ebrs, err := leerEBRs(diskPath, p.Part_start)
				if err != nil {
					color.Red("[REP-MBR]: Error al leer EBRs: %v", err)
				} else {
					for j, ebr := range ebrs {
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td colspan="2" bgcolor="#28A745"><font color="white"><b>Partición Lógica %d</b></font></td></tr>`+"\n", j+1))
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_mount</td><td>%c</td></tr>`+"\n", ebr.Part_mount))
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_fit</td><td>%c</td></tr>`+"\n", ebr.Part_fit))
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_start</td><td>%d</td></tr>`+"\n", ebr.Part_start))
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_s</td><td>%d</td></tr>`+"\n", ebr.Part_s))
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_next</td><td>%d</td></tr>`+"\n", ebr.Part_next))
						dotBuffer.WriteString(fmt.Sprintf(`<tr><td>part_name</td><td>%s</td></tr>`+"\n", utils.ToString(ebr.Name[:])))
					}
				}
			}
		}
	}

	dotBuffer.WriteString("\t</table>\n\t>];\n}\n")

	tempDotFile := "VDIC-MIA/Rep/mbr_temp.dot"
	err := os.WriteFile(tempDotFile, dotBuffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("[REP-MBR]: Error al escribir el archivo .dot temporal: %v", err)
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(ruta_reporte), ".")
	if fileExtension == "" {
		return "", fmt.Errorf("[REP-MBR]: La ruta del reporte no tiene una extensión válida (ej: .png, .jpg, .pdf)")
	}

	cmd := exec.Command("dot", "-T"+fileExtension, tempDotFile, "-o", ruta_reporte)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		errorMsg := fmt.Sprintf("[REP-MBR]: Error al ejecutar 'dot' de Graphviz. Asegúrate de que esté instalado y en el PATH. Error: %v. Stderr: %s", err, stderr.String())
		return "", errors.New(errorMsg)
	}

	os.Remove(tempDotFile)

	return fmt.Sprintf("Reporte MBR generado exitosamente en '%s'", ruta_reporte), nil
}

func leerEBRs(diskPath string, start int32) ([]structures.EBR, error) {
	file, err := os.Open(diskPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var ebrs []structures.EBR
	nextEbrStart := start

	for nextEbrStart != -1 {
		ebr := structures.EBR{}
		_, err := file.Seek(int64(nextEbrStart), 0)
		if err != nil {
			return nil, fmt.Errorf("error al buscar el EBR en la posición %d: %v", nextEbrStart, err)
		}

		if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
			// Fin de la lectura o error real
			break
		}

		if ebr.Part_s > 0 {
			ebrs = append(ebrs, ebr)
		}

		nextEbrStart = ebr.Part_next
	}

	return ebrs, nil
}
