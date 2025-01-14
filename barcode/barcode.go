// Funciones para generar códigos de barras 1D en formato SVG.
/*
Se han seguido las mismas reglas de generación que usan las impresoras
ESC/POS en el comando GS k, consultar https://download4.epson.biz/sec_pubs/pos/reference_en/escpos/gs_lk.html
*/
package barcode

import (
	"fmt"
)

type KIND int

const (
	C39     KIND = 1  // Code39
	C93     KIND = 2  // Code93
	I25     KIND = 3  // Interleaved 2/5
	C128    KIND = 4  // Code128
	C128A   KIND = 5  // Code128A
	C128B   KIND = 6  // Code128B
	C128C   KIND = 7  // Code128C
	EAN13   KIND = 8  // EAN-13
	EAN8    KIND = 9  // EAN-8
	UPCA    KIND = 10 // UPC-A
	UPCE    KIND = 11 // UPC-E
	CODABAR KIND = 12 // CODABAR
)

type HRI int

const (
	None  HRI = 0
	Above HRI = 1
	Below HRI = 2
	Both  HRI = 3
)

// Genera un documento SVG:
//   - code: código a representar
//   - kind: tipo de código de barras
//   - w: ancho del módulo
//   - h: altura del código, incluye texto HRI
//   - hri: posición del texto HRI
//   - inline: suprime cabeceras SVG para embeber en HTML
func GetBarcodeSVG(code string, kind KIND, w, h int, color string, hri HRI, inline bool) (string, error) {
	barcodeArray, err := GetBarcodeBARS(code, kind)
	if err != nil {
		return "", err
	}
	var svg string
	if !inline {
		svg += "<?xml version=\"1.0\" standalone=\"no\" ?>\n"
		svg += "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n"
	}
	modulos := 10 // Dejamos 5 módulos a cada lado de guarda
	for _, r := range barcodeArray {
		modulos += int(r) - 48
	}
	x := 5 * w
	y := 0
	bh := h
	switch hri {
	case Above:
		y = 14
		h += 14
	case Below:
		h += 14
	case Both:
		y = 14
		h += 28
	}
	svg += fmt.Sprintf("<svg width=\"%d\" height=\"%d\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">\n", modulos*w, h)
	svg += fmt.Sprintf("\t<g fill=\"%s\" stroke=\"none\">\n", color)
	negro := true
	for _, r := range barcodeArray {
		bw := int(r-48) * w
		if negro {
			svg += fmt.Sprintf("\t\t<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" />\n", x, y, bw, bh)
		}
		x += bw
		negro = !negro
	}
	if hri == Above || hri == Both {
		y := 2
		svg += fmt.Sprintf("\t<text x=\"%d\" text-anchor=\"middle\" dominant-baseline=\"hanging\" y=\"%d\" fill=\"%s\" font-size=\"12px\">%s</text>\n", modulos*w/2, y, color, code)
	}
	if hri == Below || hri == Both {
		y := h - 14
		svg += fmt.Sprintf("\t<text x=\"%d\" text-anchor=\"middle\" dominant-baseline=\"hanging\" y=\"%d\" fill=\"%s\" font-size=\"12px\">%s</text>\n", modulos*w/2, y, color, code)
	}
	svg += "\t</g>\n</svg>\n"
	return svg, nil
}

// Genera la secuencia de módulos de las barras de código:
//   - code: código a representar
//   - kind: tipo de código de barras
func GetBarcodeBARS(code string, kind KIND) (string, error) {
	switch kind {
	case C39:
		return barcodeCode39(code)
	case C93:
		return barcodeCode93(code)
	case I25:
		return barcodeI25(code)
	case C128:
		return barcodeC128(code, ' ')
	case C128A:
		return barcodeC128(code, 'A')
	case C128B:
		return barcodeC128(code, 'B')
	case C128C:
		return barcodeC128(code, 'C')
	case EAN13:
		return barcodeEAN13(code)
	case EAN8:
		return barcodeEAN8(code)
	case UPCA:
		return barcodeUPCA(code)
	case UPCE:
		return barcodeUPCE(code)
	case CODABAR:
		return barcodeCODABAR(code)
	default:
		return "", fmt.Errorf("unsuported type %q", kind)
	}
}
