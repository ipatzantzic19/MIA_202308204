package report

import (
	"Proyecto/Estructuras/structures"
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

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

type DiskSegment struct {
	Type      string
	Name      string
	Start     int32
	Size      int32
	IsEBR     bool
	Contained []*DiskSegment
}

func Report_DISK(id string, rutaReporte string) (string, error) {

	part, ok := utils.ObtenerDiscoID(id)
	if !ok {
		return "", fmt.Errorf("[REP-DISK]: No se encontró la partición con id %s", id)
	}

	mbr, ok := utils.Obtener_FULL_MBR_FDISK(part.Path)
	if !ok {
		return "", fmt.Errorf("[REP-DISK]: No se pudo leer el MBR del disco")
	}

	diskSize := float64(mbr.Mbr_tamano)
	mbrSize := int32(binary.Size(mbr))

	var parts []DiskSegment

	for _, p := range mbr.Mbr_partitions {
		if p.Part_status != '0' && p.Part_s > 0 {

			seg := DiskSegment{
				Name:  strings.TrimSpace(utils.ToString(p.Part_name[:])),
				Start: p.Part_start,
				Size:  p.Part_s,
			}

			if p.Part_type == 'E' {
				seg.Type = "Extendida"

				file, _ := os.Open(part.Path)
				defer file.Close()

				next := seg.Start
				for next != -1 {
					ebr := structures.EBR{}
					file.Seek(int64(next), 0)
					if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
						break
					}

					if ebr.Part_s > 0 {
						seg.Contained = append(seg.Contained,
							&DiskSegment{Type: "EBR", IsEBR: true, Size: int32(binary.Size(ebr))},
							&DiskSegment{
								Type:  "Logica",
								Name:  strings.TrimSpace(utils.ToString(ebr.Name[:])),
								Start: ebr.Part_start,
								Size:  ebr.Part_s,
							})
					}
					next = ebr.Part_next
				}
			} else {
				seg.Type = "Primaria"
			}

			parts = append(parts, seg)
		}
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Start < parts[j].Start
	})

	var layout []DiskSegment
	cursor := mbrSize

	for _, p := range parts {
		if p.Start > cursor {
			layout = append(layout, DiskSegment{
				Type: "Libre",
				Size: p.Start - cursor,
			})
		}
		layout = append(layout, p)
		cursor = p.Start + p.Size
	}

	if cursor < mbr.Mbr_tamano {
		layout = append(layout, DiskSegment{
			Type: "Libre",
			Size: mbr.Mbr_tamano - cursor,
		})
	}

	// Fusionar segmentos "Libre" adyacentes y eliminar tamaños cero
	var merged []DiskSegment
	for _, s := range layout {
		if s.Size <= 0 {
			continue
		}
		if len(merged) == 0 {
			merged = append(merged, s)
			continue
		}
		last := &merged[len(merged)-1]
		if last.Type == "Libre" && s.Type == "Libre" {
			last.Size += s.Size
		} else {
			merged = append(merged, s)
		}
	}
	layout = merged

	// Si existe un espacio libre inmediatamente después del MBR y es muy pequeño,
	// lo incorporamos al MBR para evitar mostrar un celda "Libre" de 0%.
	if len(layout) > 0 && layout[0].Type == "Libre" {
		leadingFreePct := (float64(layout[0].Size) / diskSize) * 100
		if leadingFreePct > 0 && leadingFreePct < 0.5 {
			// añadir al tamaño del MBR visual y eliminar el segmento libre
			mbrSize += layout[0].Size
			layout = layout[1:]
		}
	}

	// base para porcentajes: incluir todo el disco (MBR + particiones + libres)
	base := diskSize

	type pctSeg struct {
		seg DiskSegment
		pct float64
	}

	var list []pctSeg
	var total float64

	for _, s := range layout {
		p := (float64(s.Size) / base) * 100
		list = append(list, pctSeg{s, p})
		total += p
	}

	if len(list) > 0 {
		list[len(list)-1].pct += 100 - total
	}

	var dot bytes.Buffer
	dot.WriteString("digraph G {\n")
	dot.WriteString("\tnode [shape=plaintext];\n")
	dot.WriteString("\tdisk [label=<\n")
	dot.WriteString("<table border='4' cellborder='2' cellspacing='4' color='#4A148C'>\n<tr>\n")

	// MBR
	dot.WriteString("<td bgcolor='#E1BEE7'><b>MBR</b></td>\n")

	for _, item := range list {
		s := item.seg
		p := item.pct

		switch s.Type {

		case "Primaria":
			dot.WriteString(fmt.Sprintf(
				"<td bgcolor='#CE93D8'><b>Primaria</b><br/>%s<br/>%.2f%%</td>\n",
				escapeHTML(s.Name), p))

		case "Libre":
			dot.WriteString(fmt.Sprintf(
				"<td bgcolor='#F3E5F5'><b>Libre</b><br/>%.2f%%</td>\n", p))

		case "Extendida":
			dot.WriteString("<td bgcolor='#BA68C8'>")
			dot.WriteString("<table border='3' cellborder='2' cellspacing='3' color='#4A148C'>")
			dot.WriteString(fmt.Sprintf("<tr><td colspan='20'><b>Extendida</b><br/>%.2f%%</td></tr><tr>", p))

			for _, c := range s.Contained {
				if c.IsEBR {
					dot.WriteString("<td bgcolor='#D1C4E9'><b>EBR</b></td>")
				} else {
					cp := (float64(c.Size) / base) * 100
					dot.WriteString(fmt.Sprintf(
						"<td bgcolor='#B39DDB'><b>Lógica</b><br/>%s<br/>%.2f%%</td>",
						escapeHTML(c.Name), cp))
				}
			}
			dot.WriteString("</tr></table></td>\n")
		}
	}

	dot.WriteString("</tr></table>>];\n}\n")

	tempDot := "VDIC-MIA/Rep/disk_temp.dot"
	os.WriteFile(tempDot, dot.Bytes(), 0644)

	ext := strings.TrimPrefix(filepath.Ext(rutaReporte), ".")
	if ext == "" {
		ext = "png"
		rutaReporte += ".png"
	}

	cmd := exec.Command("dot", "-T"+ext, tempDot, "-o", rutaReporte)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	os.Remove(tempDot)

	return fmt.Sprintf("[REP-DISK]: Reporte generado correctamente en %s", rutaReporte), nil
}
