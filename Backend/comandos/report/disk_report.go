package report

import (
	"Proyecto/comandos/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// DiskSegment represents a portion of the disk, which can be a partition or free space.
type DiskSegment struct {
	Type      string // "Primaria", "Extendida", "Logica", "Libre", "EBR"
	Name      string
	Start     int32
	Size      int32
	Percent   float64
	IsEBR     bool
	Contained []*DiskSegment
}

// Report_DISK generates a visual report of the disk structure, including partitions and free space.
func Report_DISK(id_particion string, ruta_reporte string) (string, error) {
	// 1. Get disk path from mounted partition ID
	particion_montada, encontrada := utils.ObtenerDiscoID(id_particion)
	if !encontrada {
		return "", fmt.Errorf("[REP-DISK]: No se encontró una partición montada con el id '%s'", id_particion)
	}
	diskPath := particion_montada.Path

	// 2. Read the MBR
	mbr, success := utils.Obtener_FULL_MBR_FDISK(diskPath)
	if !success {
		return "", fmt.Errorf("[REP-DISK]: No se pudo leer el MBR del disco en '%s'", diskPath)
	}
	diskTotalSize := float64(mbr.Mbr_tamano)
	sizeOfMBR := int32(binary.Size(mbr))

	// 3. Collect all partitions (primary and logical) into a single list
	allParts := []DiskSegment{}
	var extendedPartition *DiskSegment

	for _, p := range mbr.Mbr_partitions {
		if p.Part_status != '0' {
			segment := DiskSegment{
				Name:  strings.TrimSpace(utils.ToString(p.Part_name[:])),
				Start: p.Part_start,
				Size:  p.Part_s,
			}
			if p.Part_type == 'E' {
				segment.Type = "Extendida"
				extendedPartition = &segment
			} else {
				segment.Type = "Primaria"
			}
			allParts = append(allParts, segment)
		}
	}

	// Read EBRs if an extended partition exists
	if extendedPartition != nil {
		ebrs, err := leerEBRs(diskPath, extendedPartition.Start)
		if err != nil {
			return "", fmt.Errorf("[REP-DISK]: Error al leer EBRs: %v", err)
		}
		for _, ebr := range ebrs {
			extendedPartition.Contained = append(extendedPartition.Contained, &DiskSegment{
				IsEBR: true, Size: int32(binary.Size(ebr)),
			})
			extendedPartition.Contained = append(extendedPartition.Contained, &DiskSegment{
				Type:  "Logica",
				Name:  strings.TrimSpace(utils.ToString(ebr.Name[:])),
				Start: ebr.Part_start,
				Size:  ebr.Part_s,
			})
		}
	}

	// 4. Sort partitions by start byte to easily find free space
	sort.Slice(allParts, func(i, j int) bool {
		return allParts[i].Start < allParts[j].Start
	})

	// 5. Build the final list of disk segments, including free space
	diskLayout := []DiskSegment{}
	currentPos := sizeOfMBR

	for _, p := range allParts {
		// Free space before the current partition
		if p.Start > currentPos {
			freeSize := p.Start - currentPos
			diskLayout = append(diskLayout, DiskSegment{
				Type: "Libre",
				Size: freeSize,
			})
		}
		// Add the partition itself
		diskLayout = append(diskLayout, p)
		currentPos = p.Start + p.Size
	}

	// Final free space at the end of the disk
	if currentPos < mbr.Mbr_tamano {
		freeSize := mbr.Mbr_tamano - currentPos
		diskLayout = append(diskLayout, DiskSegment{
			Type: "Libre",
			Size: freeSize,
		})
	}

	// 6. Generate DOT string
	var dotBuffer bytes.Buffer
	dotBuffer.WriteString("digraph G {\n")
	dotBuffer.WriteString("\trankdir=TB;\n")
	dotBuffer.WriteString("\tnode[shape=plaintext];\n")
	dotBuffer.WriteString("\tdisk_table [label=<\n")
	dotBuffer.WriteString("\t\t<table border='1' cellborder='1' cellspacing='0'>\n")
	dotBuffer.WriteString("\t\t\t<tr>\n")

	// MBR cell
	mbrPercent := (float64(sizeOfMBR) / diskTotalSize) * 100
	dotBuffer.WriteString(fmt.Sprintf("\t\t\t\t<td><b>MBR</b><br/>%.2f%%</td>\n", mbrPercent))

	// Partition and Free Space cells
	for _, segment := range diskLayout {
		percent := (float64(segment.Size) / diskTotalSize) * 100
		if segment.Type == "Extendida" {
			dotBuffer.WriteString("\t\t\t\t<td>\n")
			dotBuffer.WriteString("\t\t\t\t\t<table border='0'>\n")
			dotBuffer.WriteString(fmt.Sprintf("\t\t\t\t\t\t<tr><td colspan='100'><b>Extendida</b><br/>%.2f%%</td></tr>\n", percent))
			dotBuffer.WriteString("\t\t\t\t\t\t<tr>\n")
			// Content of the extended partition
			for _, containedSegment := range segment.Contained {
				cPercent := (float64(containedSegment.Size) / diskTotalSize) * 100
				if containedSegment.IsEBR {
					dotBuffer.WriteString(fmt.Sprintf("\t\t\t\t\t\t\t<td><b>EBR</b><br/>%.2f%%</td>\n", cPercent))
				} else {
					dotBuffer.WriteString(fmt.Sprintf("\t\t\t\t\t\t\t<td>%s: %s<br/>%.2f%%</td>\n", containedSegment.Type, containedSegment.Name, cPercent))
				}
			}
			dotBuffer.WriteString("\t\t\t\t\t\t</tr>\n")
			dotBuffer.WriteString("\t\t\t\t\t</table>\n")
			dotBuffer.WriteString("\t\t\t\t</td>\n")
		} else {
			dotBuffer.WriteString(fmt.Sprintf("\t\t\t\t<td>%s %s<br/>%.2f%%</td>\n", segment.Type, segment.Name, percent))
		}
	}

	dotBuffer.WriteString("\t\t\t</tr>\n")
	dotBuffer.WriteString("\t\t</table>\n")
	dotBuffer.WriteString("\t>];\n")
	dotBuffer.WriteString("}\n")

	// 7. Write .dot file and execute Graphviz
	tempDotFile := "VDIC-MIA/Rep/disk_temp.dot"
	err := os.WriteFile(tempDotFile, dotBuffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("[REP-DISK]: Error al escribir el archivo .dot temporal: %v", err)
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(ruta_reporte), ".")
	if fileExtension == "" {
		return "", fmt.Errorf("[REP-DISK]: La ruta del reporte no tiene una extensión válida")
	}

	cmd := exec.Command("dot", "-T"+fileExtension, tempDotFile, "-o", ruta_reporte)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("[REP-DISK]: Error al ejecutar 'dot'. Stderr: %s", stderr.String())
	}

	os.Remove(tempDotFile)

	return fmt.Sprintf("Reporte DISK generado exitosamente en '%s'", ruta_reporte), nil
}
