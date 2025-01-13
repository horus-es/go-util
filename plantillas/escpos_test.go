package plantillas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/stretchr/testify/assert"
)

func TestMergeEscPosTemplate(t *testing.T) {
	p, err := os.ReadFile("plantilla.escpos")
	assert.NoError(t, err)
	f, err := MergeEscPosTemplate("plantilla.escpos", string(p), factura, "", formato.DMA, formato.EUR)
	assert.NoError(t, err)
	os.WriteFile("escpos_test_out.escpos", []byte(f), 0666)
	//os.WriteFile("escpos_test_expect.escpos", []byte(f), 0666)
	crc1 := crc(t, "escpos_test_expect.escpos", "", "")
	crc2 := crc(t, "escpos_test_out.escpos", "", "")
	assert.Equal(t, crc1, crc2)
}

func ExampleMergeEscPosTemplate() {
	// Cargar plantilla
	plantilla, err := os.ReadFile("plantilla.escpos")
	errores.PanicIfError(err)
	// Fusionar plantilla con estructura factura
	f, err := MergeEscPosTemplate(
		"escpos",
		string(plantilla),
		factura,
		"/assets",
		formato.DMA,
		formato.EUR,
	)
	errores.PanicIfError(err)
	// Guardar salida
	os.WriteFile("recibo.escpos", []byte(f), 0666)
	fmt.Println("Generado fichero recibo.escpos")
	// Output: Generado fichero recibo.escpos
}

func TestGenerateEscPos(t *testing.T) {
	p, err := os.ReadFile("plantilla.escpos")
	assert.NoError(t, err)
	f, err := MergeEscPosTemplate("plantilla.escpos", string(p), factura, "", formato.DMA, formato.EUR)
	assert.NoError(t, err)
	escpos, err := GenerateEscPos(f)
	assert.NoError(t, err)
	//os.WriteFile("escpos_test_expect.prn", escpos, 0644)
	os.WriteFile("escpos_test_out.prn", escpos, 0644)
	expect, err := os.ReadFile("escpos_test_expect.prn")
	assert.NoError(t, err)
	assert.Equal(t, expect, escpos)
}

func ExampleGenerateEscPos() {
	// Cargar plantilla
	plantilla, err := os.ReadFile("plantilla.escpos")
	errores.PanicIfError(err)
	// Fusionar plantilla con estructura factura
	f, err := MergeEscPosTemplate(
		"escpos",
		string(plantilla),
		factura,
		"/assets",
		formato.DMA,
		formato.EUR,
	)
	errores.PanicIfError(err)
	// Convertir la plantilla fusionada a fichero esc/pos binario
	g, err := GenerateEscPos(f)
	errores.PanicIfError(err)
	// Guardar salida binaria (o enviar a impresora)
	os.WriteFile("recibo.prn", []byte(g), 0666)
	fmt.Println("Generado fichero recibo.prn")
	// Output: Generado fichero recibo.prn
}

func TestGenerateEscPosPdf(t *testing.T) {
	plantilla, err := os.ReadFile("plantilla.escpos")
	assert.NoError(t, err)
	wd, err := os.Getwd()
	assert.NoError(t, err)
	wd = "file:///" + strings.ReplaceAll(wd, "\\", "/")
	err = GenerateEscPosPdf("template", string(plantilla), factura, wd, formato.DMA, formato.EUR, "escpos_test_out.pdf", 80)
	assert.NoError(t, err)
	t1 := readPdfText(t, "escpos_test_expect.pdf")
	t2 := readPdfText(t, "escpos_test_out.pdf")
	assert.Equal(t, t1, t2)
}

func ExampleGenerateEscPosPdf() {
	// Carga plantilla HTML
	plantilla, err := os.ReadFile("plantilla.escpos")
	errores.PanicIfError(err)
	// Genera fichero PDF
	err = GenerateEscPosPdf(
		"pdf",
		string(plantilla),
		factura,
		"file:///assets",
		formato.DMA,
		formato.EUR,
		"recibo.pdf",
		80,
	)
	errores.PanicIfError(err)
	fmt.Println("Generado fichero recibo.pdf")
	// Output: Generado fichero recibo.pdf
}
