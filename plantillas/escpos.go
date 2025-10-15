// Funciones para procesar plantillas y convertirlas en esc/pos binario o en PDF
/*
En las plantillas escpos se sustituyen los siguientes comandos entre {} por las correspondientes secuencias de escape esc/pos:

Estilos {whsbuioxlrc}, cada letra es opcional y significan:
  - w - doble ancho
  - h - doble alto
  - s - pequeño
  - b - negrita
  - u - subrayado
  - i - itálica
  - o - blanco sobre negro
  - x - arriba-abajo
  - l - izquierda
  - r - derecha
  - c - centrado

Control de la impresora:
  - {reset}
  - {full-cut}
  - {partial-cut}
  - {form-feed}

Códigos de barras:
  - {bc-heigth 162}: altura del código de barras (1-255)
  - {bc-modulo 3}: módulo del código de barras (2-6)
  - {bc-hri none}: muestra el texto bajo el código de barras (none/above/below/both)
  - {code128 123456}: Conmutación automática de las variantes B y C
  - {code128a 123456}: Códigos de control, números, mayúsculas y simbolos
  - {code128b 123456}: Números, mayúsculas, minúsculas y más simbolos
  - {code128c 123456}: Solo números de longitud par
  - {itf 123456}: Solo números de longitud par
  - {upc-a 123456}: Códigos numéricos para productos (usa)
  - {upc-e 123456}: Códigos cortos para productos (usa)
  - {ean-13 123456}: Códigos numéricos para productos (eur)
  - {ean-8 123456}: Códigos cortos para productos (eur)
  - {code39 123456}: Números, mayúsculas y algunos símbolos
  - {code93 123456}: Números, mayúsculas y algunos símbolos
  - {codabar 123456} (se imprime como code128 en html/pdf)

Códigos QR
  - {qr-ecc L}: ECC del código QR (L|M|Q|H)
  - {qr-modulo 3}: módulo del código QR (1-16)
  - {qr https://devel.horus.es}

Imágenes:
  - {logo.png}: fichero en formato png

Se soportan las funciones de formato DATETIME, DATE, TIME y PRICE.

Ejemplo de plantillla en https://github.com/horus-es/go-util/blob/main/plantillas/plantilla.escpos
*/
package plantillas

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/beevik/etree"
	"github.com/horus-es/go-util/v3/barcode"
	"github.com/horus-es/go-util/v3/errores"
	"github.com/horus-es/go-util/v3/formato"
	go_qr "github.com/piglig/go-qr"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Fusiona una plantilla esc/pos con un struct o map de datos.
//   - name: nombre arbitrario para la plantilla que aparece en los mensajes de error
//   - escpos: plantilla en formato *.escpos
//   - datos: estructura de datos para fusionar con la plantilla
//   - assets: ruta de imágenes u otros recursos
//   - ff: formato de las fechas para las funciones DATETIME y DATE
//   - fp: formato de los precios para la funcion PRICE
func MergeEscPosTemplate(name, escpos string, datos any, assets string, ff formato.Fecha, fp formato.Moneda) (string, error) {
	var funciones = template.FuncMap{
		"DATETIME": func(x any) string {
			switch t := x.(type) {
			case time.Time:
				return formato.PrintFechaHora(t, ff)
			case string:
				t2, err := formato.ParseFechaHora(t, ff)
				if err == nil {
					return formato.PrintFechaHora(t2, ff)
				}
			}
			errores.PanicIfTrue(true, "fecha %q no soportada", x)
			return ""
		},
		"DATE": func(x any) string {
			switch t := x.(type) {
			case time.Time:
				return formato.PrintFecha(t, ff)
			case string:
				t2, err := formato.ParseFechaHora(t, ff)
				if err == nil {
					return formato.PrintFecha(t2, ff)
				}
			}
			errores.PanicIfTrue(true, "fecha %q no soportada", x)
			return ""
		},
		"TIME": func(x any) string {
			switch t := x.(type) {
			case time.Time:
				return formato.PrintHora(t, false)
			case string:
				t2, err := formato.ParseFechaHora(t, ff)
				if err == nil {
					return formato.PrintHora(t2, false)
				}
			}
			errores.PanicIfTrue(true, "fecha %q no soportada", x)
			return ""
		},
		"PRICE": func(f float64) string {
			return fmt.Sprintf("%10s", formato.PrintPrecio(f, fp))
		},
	}
	var opt string
	if reflect.TypeOf(datos).Kind() == reflect.Map {
		// En los mapas se permite que falten campos
		opt = "missingkey=zero"
	} else {
		// En las estructuras se exige la existencia del dato
		opt = "missingkey=error"
	}
	tmpl, err := template.New(filepath.Join(assets, name)).Funcs(funciones).Option(opt).Parse(escpos)
	if err != nil {
		return "", err
	}
	var marshaled bytes.Buffer
	err = tmpl.Execute(&marshaled, datos)
	if err != nil {
		return "", err
	}
	return marshaled.String(), nil
}

// Expresiones regulares esc/pos
var reEstilosEscPos = regexp.MustCompile(`{[whsbuioxlrc]*}`)
var reResetEscPos = regexp.MustCompile(`{reset}`)
var reFullCutEscPos = regexp.MustCompile(`{full-cut}`)
var rePartialCutEscPos = regexp.MustCompile(`{partial-cut}`)
var reFormFeedEscPos = regexp.MustCompile(`{form-feed}`)
var rePaperWidthEscPos = regexp.MustCompile(`{paper-width ([0-9]+)}`)
var reBcHeightEscPos = regexp.MustCompile(`{bc-height ([0-9]+)}`)
var reBcModuloEscPos = regexp.MustCompile(`{bc-modulo ([0-9]+)}`)
var reBcHriEscPos = regexp.MustCompile(`{bc-hri (none|above|below|both)}`)
var reQrModuloEscPos = regexp.MustCompile(`{qr-modulo ([0-9]+)}`)
var reQrEccEscPos = regexp.MustCompile(`{qr-ecc (L|M|Q|H)}`)
var reBarcodeEscPos = regexp.MustCompile(`{(code128|code128a|code128b|code128c|itf|upc-a|upc-e|ean-13|ean-8|code39|code93|codabar) ([^{}]+)}`)
var reQrEscPos = regexp.MustCompile(`{qr ([^{}]+)}`)
var reImgEscPos = regexp.MustCompile(`{img ([^{}]+)}`)

const (
	TAB = byte(0x09)
	LF  = byte(0x0a)
	CR  = byte(0x0d)
	FF  = byte(0x0c)
	ESC = byte(0x1b)
	FS  = byte(0x1c)
	GS  = byte(0x1d)
)

// Genera un []byte esc/pos (binario) a partir de una plantilla *.escpos.
func GenerateEscPos(escpos string) (bin []byte, width int, err error) {

	// Convertimos a windows-1252
	bin = win1252(escpos)

	// Quitamos espacios iniciales y finales
	bin = bytes.TrimSpace(bin)

	// Quitamos CRs
	bin = bytes.ReplaceAll(bin, []byte{'\r'}, nil)

	// Procesamos estilos
	bin = processEscPosStyles(bin)

	// Procesamos códigos de control
	bin, width = processEscPosControls(bin)

	// Procesamos códigos de barras
	bin = processEscPosBarcodes(bin)

	// Procesamos códigos QR
	bin = processEscPosQR(bin)

	// Procesamos imágenes
	bin = processEscPosImg(bin)

	return
}

func win1252(escpos string) (bin []byte) {
	bin = []byte(escpos)
	var buf bytes.Buffer
	writer := transform.NewWriter(&buf, charmap.Windows1252.NewEncoder())
	_, err := writer.Write(bin)
	if err != nil {
		return
	}
	defer writer.Close()
	bin = append([]byte{ESC, 't', 16}, buf.Bytes()...)
	return
}

// Procesa los estilos esc/pos
func processEscPosStyles(escpos []byte) []byte {
	// Estado de los estimos de la impresora
	type tEstadoEstilos struct {
		alignment    byte
		isBold       bool
		isUnderline  bool
		isItalics    bool
		isSmall      bool
		isDoubleX    bool
		isDoubleY    bool
		isReverse    bool
		isUpsideDown bool
	}
	estado := tEstadoEstilos{alignment: 'l'}
	var result bytes.Buffer
	posiciones := reEstilosEscPos.FindAllIndex(escpos, -1)
	p := 0
	for _, tuplas := range posiciones {
		if tuplas[0]-1 >= p {
			result.Write(escpos[p:tuplas[0]])
		}
		p = tuplas[1]
		letras := escpos[tuplas[0]+1 : tuplas[1]-1]
		nuevo := tEstadoEstilos{alignment: 'l'}
		for _, letra := range letras {
			switch letra {
			case 'w':
				nuevo.isDoubleX = true
			case 'h':
				nuevo.isDoubleY = true
			case 's':
				nuevo.isSmall = true
			case 'b':
				nuevo.isBold = true
			case 'u':
				nuevo.isUnderline = true
			case 'i':
				nuevo.isItalics = true
			case 'o':
				nuevo.isReverse = true
			case 'x':
				nuevo.isUpsideDown = true
			case 'l', 'c', 'r':
				nuevo.alignment = letra
			default:
				errores.PanicIfTrue(true, "estilo %c no soportado", letra)
			}
		}
		if estado.alignment != nuevo.alignment { // debe ser lo primero de la línea!
			// ESC a
			result.WriteByte(ESC)
			result.WriteByte('a')
			switch nuevo.alignment {
			case 'l':
				result.WriteByte('0') // left
			case 'c':
				result.WriteByte('1') // center
			case 'r':
				result.WriteByte('2') // right
			}
		}
		if estado.isDoubleX != nuevo.isDoubleX || estado.isDoubleY != nuevo.isDoubleY || estado.isSmall != nuevo.isSmall || estado.isBold != nuevo.isBold || estado.isUnderline != nuevo.isUnderline {
			// ESC !
			result.WriteByte(ESC)
			result.WriteByte('!')
			var octeto byte
			if nuevo.isSmall {
				octeto |= 0x01
			}
			if nuevo.isBold { // (también ESC E)
				octeto |= 0x08
			}
			if nuevo.isDoubleY {
				octeto |= 0x10
			}
			if nuevo.isDoubleX {
				octeto |= 0x20
			}
			if nuevo.isUnderline { // (también ESC -)
				octeto |= 0x80
			}
			result.WriteByte(octeto)
		}
		if estado.isItalics != nuevo.isItalics {
			// ESC 4
			result.WriteByte(ESC)
			result.WriteByte('4')
			if nuevo.isItalics {
				result.WriteByte('1')
			} else {
				result.WriteByte('0')
			}
		}
		if estado.isUpsideDown != nuevo.isUpsideDown {
			// ESC {
			result.WriteByte(ESC)
			result.WriteByte('{')
			if nuevo.isUpsideDown {
				result.WriteByte('1')
			} else {
				result.WriteByte('0')
			}
		}
		if estado.isReverse != nuevo.isReverse {
			// GS B
			result.WriteByte(GS)
			result.WriteByte('B')
			if nuevo.isReverse {
				result.WriteByte('1')
			} else {
				result.WriteByte('0')
			}
		}
		estado = nuevo
	}
	result.Write(escpos[p:])
	return result.Bytes()
}

// Procesa las secuencias de control esc/pos y extrae el ancho del papel
func processEscPosControls(escpos []byte) ([]byte, int) {
	var width int
	result := reResetEscPos.ReplaceAll(escpos, []byte{ESC, '@'})         // ESC @
	result = reFullCutEscPos.ReplaceAll(result, []byte{GS, 'V', '0'})    // GS V 0
	result = rePartialCutEscPos.ReplaceAll(result, []byte{GS, 'V', '1'}) // GS V 1
	result = reFormFeedEscPos.ReplaceAll(result, []byte{FF})             // FF
	result = rePaperWidthEscPos.ReplaceAllFunc(result, func(match []byte) []byte {
		submatches := rePaperWidthEscPos.FindSubmatch(match)
		w, _ := strconv.Atoi(string(submatches[1]))
		width = max(width, w)
		return nil
	})
	if width == 0 {
		width = 80 // Por defecto 80mm
	}
	return result, width
}

// Procesa los códigos de barras esc/pos
func processEscPosBarcodes(escpos []byte) []byte {
	result := reBcHeightEscPos.ReplaceAllFunc(escpos, func(match []byte) []byte {
		submatches := reBcHeightEscPos.FindSubmatch(match)
		h := bytesToByte(submatches[1])
		if h > 0 {
			return []byte{GS, 'h', h} // GS h altura
		}
		return match
	})
	result = reBcModuloEscPos.ReplaceAllFunc(result, func(match []byte) []byte {
		submatches := reBcModuloEscPos.FindSubmatch(match)
		m := bytesToByte(submatches[1])
		if m >= 2 && m <= 6 {
			return []byte{GS, 'w', m} // GS w modulo
		}
		return match
	})
	result = reBcHriEscPos.ReplaceAllFunc(result, func(match []byte) []byte {
		submatches := reBcHriEscPos.FindSubmatch(match)
		hri := string(submatches[1])
		switch hri {
		case "none":
			return []byte{GS, 'H', '0'} // GS H 0
		case "above":
			return []byte{GS, 'H', '1'} // GS H 1
		case "below":
			return []byte{GS, 'H', '2'} // GS H 2
		case "both":
			return []byte{GS, 'H', '3'} // GS H 3
		}
		return match
	})
	result = reBarcodeEscPos.ReplaceAllFunc(result, func(match []byte) []byte {
		submatches := reBarcodeEscPos.FindSubmatch(match)
		tipo := string(submatches[1])
		codigo := submatches[2]
		var l byte
		if len(codigo) < 30 { // Suficiente
			l = byte(len(codigo))
		}
		if l == 0 {
			return match
		}
		switch tipo {
		case "code128":
			//return append([]byte{GS, 'k', 79, l}, codigo...) // GS k 79 l codigo // Muchas impresoras no lo soportan
			codigo = c128auto(codigo)
			l = byte(len(codigo))
			return append([]byte{GS, 'k', 73, l}, codigo...) // GS k 73 l codigo // ... por eso usamos este
		case "code128a":
			if bytesInRange(codigo, [][]byte{{0, 95}}) {
				return append([]byte{GS, 'k', 73, l + 2, '{', 'A'}, codigo...) // GS k 73 l { A codigo
			}
		case "code128b":
			if bytesInRange(codigo, [][]byte{{32, 122}}, 124, 126) {
				return append([]byte{GS, 'k', 73, l + 2, '{', 'B'}, codigo...) // GS k 73 l { B codigo
			}
		case "code128c":
			if l%2 == 0 && bytesInRange(codigo, [][]byte{{'0', '9'}}) {
				l = l / 2
				pares := make([]byte, l)
				for i := byte(0); i < l; i++ {
					pares[i] = (codigo[i*2]-48)*10 + codigo[i*2+1] - 48
				}
				return append([]byte{GS, 'k', 73, l + 2, '{', 'C'}, pares...) // GS k 73 l { C codigo
			}
		case "itf":
			if l%2 == 0 && bytesInRange(codigo, [][]byte{{'0', '9'}}) {
				return append([]byte{GS, 'k', 70, l}, codigo...) // GS k 70 l codigo
			}
		case "upc-a":
			if (l == 11 || l == 12) && bytesInRange(codigo, [][]byte{{'0', '9'}}) {
				return append([]byte{GS, 'k', 65, l}, codigo...) // GS k 65 l codigo
			}
		case "upc-e":
			if (l == 7 || l == 11) && bytesInRange(codigo, [][]byte{{'0', '9'}}) && codigo[0] == '0' {
				return append([]byte{GS, 'k', 66, l}, codigo...) // GS k 66 l codigo
			}
		case "ean-13":
			if l == 12 && bytesInRange(codigo, [][]byte{{'0', '9'}}) {
				return append([]byte{GS, 'k', 67, l}, codigo...) // GS k 67 l codigo
			}
		case "ean-8":
			if l == 7 && bytesInRange(codigo, [][]byte{{'0', '9'}}) {
				return append([]byte{GS, 'k', 68, l}, codigo...) // GS k 68 l codigo
			}
		case "code39":
			if bytesInRange(codigo, [][]byte{{'0', '9'}, {'A', 'Z'}}, ' ', '$', '%', '*', '+', '-', '.', '/') {
				return append([]byte{GS, 'k', 69, l}, codigo...) // GS k 69 l codigo
			}
		case "code93":
			if bytesInRange(codigo, [][]byte{{0, 127}}) {
				return append([]byte{GS, 'k', 72, l}, codigo...) // GS k 72 l codigo
			}
		case "codabar":
			if bytesInRange(codigo, [][]byte{{'0', '9'}}, '$', '+', '-', '.', '/', ':') {
				return append([]byte{GS, 'k', 71, l + 2, 'a'}, append(codigo, 'a')...) // GS k 71 l codigo
			}
			if bytesInRange(codigo, [][]byte{{'0', '9'}, {'A', 'D'}, {'a', 'd'}}, '$', '+', '-', '.', '/', ':') {
				return append([]byte{GS, 'k', 71, l}, codigo...) // GS k 71 l codigo
			}
		}
		return match
	})
	return result
}

// Procesa los códigos QR esc/pos
func processEscPosQR(escpos []byte) []byte {
	result := reQrModuloEscPos.ReplaceAllFunc(escpos, func(match []byte) []byte {
		submatches := reQrModuloEscPos.FindSubmatch(match)
		m := bytesToByte(submatches[1])
		if m >= 1 && m <= 16 {
			return []byte{GS, '(', 'k', 3, 0, '1', 67, m} // GS ( k ... 1 67 modulo
		}
		return match
	})
	result = reQrEccEscPos.ReplaceAllFunc(result, func(match []byte) []byte {
		submatches := reQrEccEscPos.FindSubmatch(match)
		ecc := string(submatches[1])
		switch ecc {
		case "L":
			return []byte{GS, '(', 'k', 3, 0, '1', 69, '0'} // GS ( k ... 1 69 0
		case "M":
			return []byte{GS, '(', 'k', 3, 0, '1', 69, '1'} // GS ( k ... 1 69 1
		case "Q":
			return []byte{GS, '(', 'k', 3, 0, '1', 69, '2'} // GS ( k ... 1 69 2
		case "H":
			return []byte{GS, '(', 'k', 3, 0, '1', 69, '3'} // GS ( k ... 1 69 3
		default:
			return match
		}
	})
	result = reQrEscPos.ReplaceAllFunc(result, func(match []byte) []byte {
		submatches := reQrEscPos.FindSubmatch(match)
		codigo := submatches[1]
		p := len(codigo) + 3
		if p < 4 || p > 7092 {
			return match
		}
		var pL byte = byte(p % 256)
		var pH byte = byte(p / 256)
		datos := []byte{GS, '(', 'k', pL, pH, '1', 80, '0'}     // GS ( k ... 1 80 0
		datos = append(datos, codigo...)                        // codigo
		datos = append(datos, GS, '(', 'k', 3, 0, '1', 81, '0') // GS ( k ... 1 81 0
		return datos
	})
	return result
}

func processEscPosImg(escpos []byte) []byte {
	result := reImgEscPos.ReplaceAllFunc(escpos, func(match []byte) []byte {
		submatches := reImgEscPos.FindSubmatch(match)
		raster, width, height, err := rasterize(string(submatches[1]))
		if err == nil {
			size := len(raster) + 10
			var datos []byte
			if size < 256*256 {
				pL := byte(size)
				pH := byte(size >> 8)
				datos = []byte{GS, '(', 'L', pL, pH} // GS ( L ...
			} else {
				p1 := byte(size)
				p2 := byte(size >> 8)
				p3 := byte(size >> 16)
				p4 := byte(size >> 24)
				datos = []byte{GS, '8', 'L', p1, p2, p3, p4} // GS 8 L ...
			}
			xL := byte(width)
			xH := byte(width >> 8)
			yL := byte(height)
			yH := byte(height >> 8)
			datos = append(datos, '0', 112, '0', 1, 1, '1', xL, xH, yL, yH) // 0 112 ...
			datos = append(datos, raster...)                                // raster
			datos = append(datos, GS, '(', 'L', 2, 0, '0', '2')             // GS ( L ... 0 2
			return datos
		}
		return match
	})
	return result
}

// Función auxiliar, convierte un array representando un numero a un byte. Si error, devuelve 0.
func bytesToByte(array []byte) byte {
	n, err := strconv.Atoi(string(array))
	if err != nil || n > 255 || n < 1 {
		return 0
	}
	return byte(n)
}

// Función auxiliar, determina si todos los caracteres del código están en los rangos o en el conjunto
func bytesInRange(codigo []byte, rangos [][]byte, conjunto ...byte) bool {
	for _, b := range codigo {
		f := false
		for _, r := range rangos {
			if b >= r[0] && b <= r[1] {
				f = true
				break
			}
		}
		if !f && conjunto != nil {
			f = bytes.IndexByte(conjunto, b) >= 0
		}
		if !f {
			return false
		}
	}
	return true
}

// "Rasteriza" una imagen
func rasterize(file string) (data []byte, width, height int, err error) {
	// Cargar la imagen
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	// Decodificar la imagen
	img, _, err := image.Decode(f)
	if err != nil {
		return
	}
	// Obtener las dimensiones de la imagen
	bounds := img.Bounds()
	width = bounds.Dx()
	height = bounds.Dy()
	rowSize := (width + 7) / 8
	m16 := uint32(1<<16 - 1)
	// Crear un array de bytes para almacenar los datos rasterizados
	data = make([]byte, rowSize*height)
	// Hallar umbral
	max := uint32(0)
	min := uint32(1<<32 - 1)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			gray := a*(299*r+587*g+114*b)/m16 + 1000*(m16-a)
			if gray < min {
				min = gray
			}
			if gray > max {
				max = gray
			}
		}
	}
	umbral := (max + min) / 2
	// Binarizar la imagen
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Obtener el color del píxel
			r, g, b, a := img.At(x, y).RGBA()
			// Convertir a escala de grises (sobre fondo blanco)
			gray := a*(299*r+587*g+114*b)/m16 + 1000*(m16-a)
			// Umbral
			if gray < umbral {
				// Empaquetar el bit en el byte correspondiente
				byteIndex := y*rowSize + x/8
				bitIndex := 7 - x%8
				data[byteIndex] |= 1 << bitIndex
			}
		}
	}
	return
}

// Genera un fichero PDF a partir de una secuencia esc/pos. Parámetros:
//   - escpos: secuencia binaria esc/pos, tal y como se enviaría a la impresora
//   - out: fichero PDF de salida
//   - width: ancho del papel en mm
//   - opciones: opciones adicionales utilidad wkhtmltopdf (ver https://wkhtmltopdf.org/usage/wkhtmltopdf.txt)
func GenerateEscPosPdf(escpos []byte, out string, width int, opciones ...string) error {
	// Fichero temporal
	tmp, err := os.CreateTemp("", "horus-*.html")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	// Documento html
	tmp.WriteString("<!DOCTYPE html>\n")
	tmp.WriteString("<html>\n")
	tmp.WriteString("<head>\n")
	tmp.WriteString("<title>Recibo</title>\n")
	tmp.WriteString("<meta http-equiv=\"Content-Type\" content=\"text/html; charset=Windows-1252\" />\n")
	tmp.WriteString("<style>\n")
	// Añadimos CSS para tickets
	addEscPosCSS(tmp, width)
	tmp.WriteString("</style>\n")
	tmp.WriteString("</head>\n")
	tmp.WriteString("<body>\n")
	// Añadimos HTML ticket
	addEscPosHTML(tmp, escpos)
	tmp.WriteString("</body>\n")
	tmp.WriteString("</html>\n")
	err = tmp.Close()
	if err != nil {
		return err
	}
	// Ejecución wkhtmltopdf
	args := append([]string{"-q", "--enable-local-file-access", "--no-outline"}, opciones...)
	args = append(args, tmp.Name(), out)
	cmd := exec.Command("wkhtmltopdf", args...)
	var log bytes.Buffer
	cmd.Stdout = &log
	cmd.Stderr = &log
	err = cmd.Run()
	if err != nil {
		s := "wkhtmltopdf " + strings.Join(args, " ") + "\n" + err.Error() + "\n" + log.String()
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\n", " -> ")
		return errors.New(s)
	}
	return nil
}

// Añade el CSS necesario para esc/pos
func addEscPosCSS(html io.Writer, width int) {
	io.WriteString(html, "escpos { font-family: 'DejaVu Sans Mono', monospace; font-size: 12px; white-space: pre-wrap; display: inline-block; border: 1px solid black; padding: 1em; margin: 1em; word-break: break-all; vertical-align: top; width: "+strconv.Itoa(width)+"mm; }\n")
	io.WriteString(html, "escpos .bold { font-weight: bold; }\n")
	io.WriteString(html, "escpos .underline { text-decoration: underline; }\n")
	io.WriteString(html, "escpos .italics { font-style: italic; }\n")
	io.WriteString(html, "escpos .center { display:inline-block; width: 100%; text-align: center; }\n")
	io.WriteString(html, "escpos .right { display:inline-block; width: 100%; text-align: right; }\n")
	io.WriteString(html, "escpos .img-left { display:block; margin-right: auto; }\n")
	io.WriteString(html, "escpos .img-center { display:block; margin-left: auto; margin-right: auto; }\n")
	io.WriteString(html, "escpos .img-right { display:block;  margin-left: auto; }\n")
	io.WriteString(html, "escpos .doubleY { display: inline-block; scale: 1 2; transform-origin: bottom; margin-top: 1em; }\n")
	io.WriteString(html, "escpos .doubleX { font-size: 2em; display: inline-block; transform-origin: bottom; scale: 1 0.5; ; margin-top: -1em;}\n")
	io.WriteString(html, "escpos .double { font-size: 2em; }\n")
	io.WriteString(html, "escpos .small { font-size: 0.75em; }\n")
	io.WriteString(html, "escpos .small.doubleX { font-size: 1.5em; }\n")
	io.WriteString(html, "escpos .small.double { font-size: 1.5em; }\n")
	io.WriteString(html, "escpos .reverse { background-color: black; color: white; }\n")
	io.WriteString(html, "escpos .upsidedown { display: inline-block; scale: -1 -1; }\n")
}

// Añade el HTML de las etiquetas esc/pos
func addEscPosHTML(html io.Writer, escpos []byte) {
	// Estado inicial
	alignment := "left"
	isBold := false
	isUnderline := false
	isItalics := false
	isSmall := false
	isDoubleX := false
	isDoubleY := false
	isReverse := false
	isUpsideDown := false
	inLabel := false
	bcHeight := 162
	bcWidth := 3
	bcHRI := barcode.None
	qrModulo := 3
	qrECC := 48
	var qrData []byte
	var img *image.Gray

	textBuffer := strings.Builder{}
	currentClass := alignment
	col := 0

	writeToHtml := func(s string) {
		if len(s) == 0 {
			return
		}
		if !inLabel {
			inLabel = true
			io.WriteString(html, "<escpos>")
		}
		io.WriteString(html, s)
	}

	// Vacía el buffer y lo agrega a la salida
	flushBuffer := func() {
		if textBuffer.Len() == 0 {
			return
		}
		writeToHtml(fmt.Sprintf("<span class=\"%s\">%s</span>", currentClass, textBuffer.String()))
		textBuffer.Reset()
	}

	// Obtiene la clase actual para textos
	getTextClass := func() string {
		class := []string{alignment}
		if isBold {
			class = append(class, "bold")
		}
		if isUnderline {
			class = append(class, "underline")
		}
		if isItalics {
			class = append(class, "italics")
		}
		if isDoubleX && !isDoubleY {
			class = append(class, "doubleX")
		}
		if !isDoubleX && isDoubleY {
			class = append(class, "doubleY")
		}
		if isDoubleY && isDoubleX {
			class = append(class, "double")
		}
		if isSmall {
			class = append(class, "small")
		}
		if isReverse {
			class = append(class, "reverse")
		}
		if isUpsideDown {
			class = append(class, "upsidedown")
		}
		return strings.Join(class, " ")
	}

	for i := 0; i < len(escpos); i++ {
		switch escpos[i] {
		case LF: // Nueva línea
			flushBuffer()
			writeToHtml("\n")
			col = 0

		case CR: // Retorno de carro
			// ignoramos

		case FF: // Fin de etiqueta
			flushBuffer()
			if inLabel {
				inLabel = false
				io.WriteString(html, "</escpos>\n")
			}

		case TAB: // Tabulaciones (convertir en espacios cada 8 posiciones)
			for col%8 != 0 {
				textBuffer.WriteByte(' ')
				col++
			}

		case ESC:
			if i+1 < len(escpos) {
				var next byte
				if i+2 < len(escpos) {
					next = escpos[i+2]
				}
				switch escpos[i+1] {
				case '@': // ESC @ (reset)
					flushBuffer()
					isBold = false
					isUnderline = false
					isItalics = false
					isSmall = false
					isReverse = false
					isUpsideDown = false
					isDoubleX = false
					isDoubleY = false
					alignment = "left"
					bcHeight = 162
					bcWidth = 3
					bcHRI = barcode.None
					qrModulo = 3
					qrECC = 48
					qrData = nil
					img = nil
					currentClass = alignment
					col = 0
					i += 1
				case '!': // ESC ! (tamaño de fuente, negrita, subrrayado)
					flushBuffer()
					isSmall = (next & 0x01) > 0
					isBold = (next & 0x08) > 0
					isDoubleY = (next & 0x10) > 0
					isDoubleX = (next & 0x20) > 0
					isUnderline = (next & 0x80) > 0
					i += 2
				case '-': // ESC - (subrayado)
					flushBuffer()
					isUnderline = next == 1 || next == 2 || next == '1' || next == '2'
					i += 2
				case '{': // ESC { (arriba/abajo)
					flushBuffer()
					isUpsideDown = next%2 == 1
					i += 2
				case '4': // ESC 4 (itálico)
					flushBuffer()
					isItalics = next%2 == 1
					i += 2
				case 'E': // ESC E (negrita)
					flushBuffer()
					isBold = next%2 == 1
					i += 2
				case 'a': // ESC a (alineación)
					flushBuffer()
					switch next {
					case 0, '0':
						alignment = "left"
					case 1, '1':
						alignment = "center"
					case 2, '2':
						alignment = "right"
					}
					i += 2
				case 'd': // ESC d (n saltos de línea)
					flushBuffer()
					for next > 0 {
						writeToHtml("\n")
						next--
					}
					col = 0
					i += 2
				case 'i', 'm': // ESC i/m (corte total/parcial)
					flushBuffer()
					if inLabel {
						inLabel = false
						io.WriteString(html, "</escpos>\n")
					}
					i += 2
				case 'p': // ESC p (pulso)
					i += 4
				case 't': // ESC t (página de código) 16=WIN1252
					i += 2
				}
			}

		case GS:
			if i+1 < len(escpos) {
				var next byte
				if i+2 < len(escpos) {
					next = escpos[i+2]
				}
				switch escpos[i+1] {
				case 'B': // GS B (blanco sobre negro)
					flushBuffer()
					isReverse = next%2 == 1
					i += 2
				case 'V': // GS V (corte total/parcial)
					flushBuffer()
					if inLabel {
						inLabel = false
						io.WriteString(html, "</escpos>\n")
					}
					if next == 0 || next == 1 || next == '0' || next == '1' {
						i += 2
					} else {
						i += 3
					}
				case '(':
					// GS (L (bitmap <64k)
					if next == 'L' && i+4 < len(escpos) {
						z := int(escpos[i+3]) + int(escpos[i+4])*256
						i += 4
						if z > 10 && escpos[i+1] == '0' && escpos[i+2] == 112 {
							// store raster image
							img = decodeRastrerImage(&escpos, i, z)
						}
						if z > 10 && escpos[i+1] == '0' && escpos[i+2] == 113 {
							// TODO: store column image
							//img = decodeColumnImage(&escpos, i, z)
						}
						if z == 2 && escpos[i+1] == '0' && (escpos[i+2] == 2 || escpos[i+2] == '2') {
							// print image
							flushBuffer()
							writeToHtml(encodeImage(img, alignment))
						}
						i += z
					}
					// GS (k (2d barcode)
					if next == 'k' && i+4 < len(escpos) {
						z := int(escpos[i+3]) + int(escpos[i+4])*256
						i += 4
						// QR
						if z == 3 && escpos[i+1] == '1' && escpos[i+2] == 67 {
							qrModulo = int(escpos[i+3])
						}
						if z == 3 && escpos[i+1] == '1' && escpos[i+2] == 69 {
							qrECC = int(escpos[i+3])
						}
						if z > 3 && escpos[i+1] == '1' && escpos[i+2] == 80 && escpos[i+3] == '0' {
							qrData = escpos[i+4 : i+z+1]
						}
						if z == 3 && escpos[i+1] == '1' && escpos[i+2] == 81 && escpos[i+3] == '0' {
							flushBuffer()
							writeToHtml(imprimeQR(qrData, qrModulo, qrECC))
						}
						i += z
					}
				case '8': // GS 8L (bitmap >64k)
					if next == 'L' && i+6 < len(escpos) {
						z := int(escpos[i+3]) + int(escpos[i+4])*256 + int(escpos[i+5])*65536 + int(escpos[i+6])*16777216
						i += 6
						if z > 10 && escpos[i+1] == '0' && escpos[i+2] == 112 {
							// store raster image
							img = decodeRastrerImage(&escpos, i, z)
						}
						if z > 10 && escpos[i+1] == '0' && escpos[i+2] == 113 {
							// TODO: store column image
							//img = decodeColumnImage(&escpos, i, z)
						}
						i += z
					}
				case 'h': // GS h (barcode height)
					bcHeight = int(next)
					i += 2
				case 'w': // GS w (barcode width)
					bcWidth = int(next)
					i += 2
				case 'H': // GS H (barcode show)
					switch next {
					case 0, '0':
						bcHRI = barcode.None
					case 1, '1':
						bcHRI = barcode.Above
					case 2, '2':
						bcHRI = barcode.Below
					case 3, '3':
						bcHRI = barcode.Both
					}
					i += 2
				case 'k': // GS k (print barcode)
					z := 0
					if next <= 6 {
						i += 2
						for i < len(escpos) {
							if escpos[i+z] == 0 {
								break
							}
							z++
						}
					}
					if next >= 65 && next <= 79 {
						i += 3
						z = int(escpos[i])
					}
					if z > 0 {
						codigo := string(escpos[i+1 : i+z+1])
						flushBuffer()
						writeToHtml(imprimeBC(codigo, next, bcWidth, bcHeight, bcHRI))
						i += z
					}
				}
			}

		case FS:
			if i+1 < len(escpos) {
				switch escpos[i+1] {
				case '.': // Cancel Kanji character mode
					i++
				}
			}

		default: // Texto normal
			newClass := getTextClass()
			if newClass != currentClass {
				flushBuffer()
				currentClass = newClass
			}
			col++
			textBuffer.WriteByte(escpos[i])
		}
	}

	flushBuffer()
	if inLabel {
		io.WriteString(html, "</escpos>\n")
	}
}

// Codificar la imagen en png+base64
func encodeImage(img *image.Gray, class string) string {
	var buffer bytes.Buffer
	err := png.Encode(&buffer, img)
	if err == nil {
		b64 := base64.StdEncoding.EncodeToString(buffer.Bytes())
		return fmt.Sprintf(`<img class="img-%s" style="width: %dpx; height: %dpx;" src="data:image/png;base64,%s"/>`, class, img.Rect.Dx()/2, img.Rect.Dy()/2, b64)
	}
	return ""
}

// Decodifica una imagen raster
func decodeRastrerImage(escpos *[]byte, i, z int) *image.Gray {
	w := int((*escpos)[i+7]) + int((*escpos)[i+8])*256
	h := int((*escpos)[i+9]) + int((*escpos)[i+10])*256
	img := image.NewGray(image.Rect(0, 0, w-1, h-1))
	if w%8 != 0 {
		// Ajustar ancho a múltiplo de 8
		w = w + 8 - w%8
	}
	for p := 11; p < z; p++ {
		x := (p - 11) * 8 % w
		y := (p - 11) * 8 / w
		for b := 0; b < 8; b++ {
			if ((*escpos)[i+p] & (128 >> b)) == 0 {
				img.SetGray(x+b, y, color.Gray{255})
			}
		}
	}
	return img
}

func imprimeBC(codigo string, bcKind byte, bcWidth, bcHeight int, hri barcode.HRI) string {
	var tipo barcode.KIND
	switch bcKind {
	case 0, 65:
		tipo = barcode.UPCA // UPC-A
	case 1, 66:
		tipo = barcode.UPCE // UPC-E
	case 2, 67:
		tipo = barcode.EAN13 // EAN-13
	case 3, 68:
		tipo = barcode.EAN8 // EAN-8
	case 4, 69:
		tipo = barcode.C39 // CODE39
	case 5, 70:
		tipo = barcode.I25 // INT 2/5
	case 6, 71:
		tipo = barcode.CODABAR // CODABAR
	case 72:
		tipo = barcode.C93 // CODE93
	case 73:
		tipo = barcode.C128 // CODE128 codificado esc/pos
	default:
		tipo = barcode.C128X // Para el resto usamos CODE128 automático
	}
	svg, _ := barcode.GetBarcodeSVG(codigo, tipo, float64(bcWidth)*3/5, float64(bcHeight)/2, "#000", hri, true)
	return svg
}

func imprimeQR(qrData []byte, qrModulo, qrECC int) string {
	qr, err := go_qr.EncodeBinary(qrData, go_qr.Ecc(qrECC-48))
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	err = qr.WriteAsSVG(go_qr.NewQrCodeImgConfig(1, 3, go_qr.WithOptimalSVG()), &buf, "#FFFFFF", "#000000")
	if err != nil {
		return ""
	}
	n := qr.GetSize() * qrModulo * 7 / 10
	doc := etree.NewDocument()
	doc.ReadFromBytes(buf.Bytes())
	// eliminamos las cabeceras xml
	for _, t := range doc.Child {
		if p, ok := t.(*etree.ProcInst); ok {
			doc.RemoveChild(p)
		}
		if p, ok := t.(*etree.Directive); ok {
			doc.RemoveChild(p)
		}
	}
	svg := doc.Root()
	svg.CreateAttr("width", strconv.Itoa(n))
	svg.CreateAttr("height", strconv.Itoa(n))
	doc.Indent(2)
	s, err := doc.WriteToString()
	if err != nil {
		return ""
	}
	return s
}

// Cambio automatico entre las variantes A, B y C
func c128auto(codigo []byte) (auto []byte) {
	variant := ' '
	i := 0
	for i < len(codigo) {
		// Averiguamos el número de digitos consecutivos a partir de la posición i-sima
		n := 0
		for i+n < len(codigo) && codigo[i+n] >= '0' && codigo[i+n] <= '9' {
			n++
		}
		// Condición para usar C: 6+ digitos en cualquier posición o 4+ digitos al principio o final
		if (n >= 6) || (n >= 4 && i == 0) || (n >= 4 && i+n == len(codigo)) {
			// Cambiamos a C si es preciso
			if variant != 'C' {
				auto = append(auto, '{', 'C')
				variant = 'C'
			}
			// Agregamos los runes en parejas
			for n > 1 {
				auto = append(auto, (codigo[i]-48)*10+codigo[i+1]-48)
				n -= 2
				i += 2
			}
		} else {
			// Cambiamos a B si es preciso
			if variant != 'B' {
				auto = append(auto, '{', 'B')
				variant = 'B'
			}
			// Agregamos los dígitos que no se hayan pasado como C
			for n > 0 {
				auto = append(auto, codigo[i])
				i++
				n--
			}
			// Agregamos hasta encontrar un dígito
			for i < len(codigo) && (codigo[i] < '0' || codigo[i] > '9') {
				r := codigo[i]
				if r == '{' {
					auto = append(auto, '{', '{')
				} else if r <= 31 {
					auto = append(auto, '{', 'S', r)
				} else if r <= 127 {
					auto = append(auto, r)
				} else {
					auto = append(auto, '?')
				}
				i++
			}
		}
	}
	return
}
