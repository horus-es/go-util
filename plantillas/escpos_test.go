package plantillas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/horus-es/go-util/v3/errores"
	"github.com/horus-es/go-util/v3/formato"
	"github.com/horus-es/go-util/v3/plantillas"
	"github.com/stretchr/testify/assert"
)

func TestMergeEscPosTemplate(t *testing.T) {
	p, err := os.ReadFile("plantilla.escpos")
	assert.NoError(t, err)
	f, err := plantillas.MergeEscPosTemplate("plantilla.escpos", string(p), factura, "", formato.DMA, formato.EUR)
	assert.NoError(t, err)
	err = os.WriteFile("escpos_test_out.escpos", []byte(f), 0666)
	assert.NoError(t, err)
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
	f, err := plantillas.MergeEscPosTemplate(
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

func TestGenerateEpson(t *testing.T) {
	p, err := os.ReadFile("plantilla.escpos")
	assert.NoError(t, err)
	f, err := plantillas.MergeEscPosTemplate("plantilla.escpos", string(p), factura, "", formato.DMA, formato.EUR)
	assert.NoError(t, err)
	escpos, _, err := plantillas.GenerateEscPos(f, plantillas.EPSON)
	assert.NoError(t, err)
	err = os.WriteFile("epson_test_out.prn", escpos, 0644)
	assert.NoError(t, err)
	//os.WriteFile("escpos_test_expect.prn", escpos, 0644)
	crc1 := crc(t, "epson_test_expect.prn", "", "")
	crc2 := crc(t, "epson_test_out.prn", "", "")
	assert.Equal(t, crc1, crc2)
}

func TestGenerateSeiko(t *testing.T) {
	p, err := os.ReadFile("plantilla.escpos")
	assert.NoError(t, err)
	f, err := plantillas.MergeEscPosTemplate("plantilla.escpos", string(p), factura, "", formato.DMA, formato.EUR)
	assert.NoError(t, err)
	escpos, _, err := plantillas.GenerateEscPos(f, plantillas.SEIKO)
	assert.NoError(t, err)
	err = os.WriteFile("seiko_test_out.prn", escpos, 0644)
	assert.NoError(t, err)
	//os.WriteFile("escpos_test_expect.prn", escpos, 0644)
	crc1 := crc(t, "seiko_test_expect.prn", "", "")
	crc2 := crc(t, "seiko_test_out.prn", "", "")
	assert.Equal(t, crc1, crc2)
}

func ExampleGenerateEscPos() {
	// Cargar plantilla
	plantilla, err := os.ReadFile("plantilla.escpos")
	errores.PanicIfError(err)
	// Fusionar plantilla con estructura factura
	f, err := plantillas.MergeEscPosTemplate(
		"escpos",
		string(plantilla),
		factura,
		"/assets",
		formato.DMA,
		formato.EUR,
	)
	errores.PanicIfError(err)
	// Convertir la plantilla fusionada a fichero esc/pos sabor EPSON (pos) o sabor SEIKO (cajero)
	prn, _, err := plantillas.GenerateEscPos(f, plantillas.EPSON)
	errores.PanicIfError(err)
	// Guardar salida binaria (o enviar a impresora)
	os.WriteFile("recibo.prn", []byte(prn), 0666)
	fmt.Println("Generado fichero recibo.prn")
	// Output: Generado fichero recibo.prn
}

func TestGenerateEscPosPdf(t *testing.T) {
	// Cargar plantilla
	plantilla, err := os.ReadFile("plantilla.escpos")
	errores.PanicIfError(err)
	// Fusionar plantilla con estructura factura
	f, err := plantillas.MergeEscPosTemplate(
		"escpos",
		string(plantilla),
		factura,
		"/assets",
		formato.DMA,
		formato.EUR,
	)
	errores.PanicIfError(err)
	// Convertir la plantilla fusionada a fichero esc/pos binario
	prn, mm, err := plantillas.GenerateEscPos(f, plantillas.SEIKO)
	errores.PanicIfError(err)
	err = plantillas.GenerateEscPosPdf(prn, "escpos_test_out.pdf", mm)
	assert.NoError(t, err)
	t1 := readPdfText(t, "escpos_test_expect.pdf")
	t2 := readPdfText(t, "escpos_test_out.pdf")
	assert.Equal(t, t1, t2)
}

func ExampleGenerateEscPosPdf() {
	// Cargar plantilla
	plantilla, err := os.ReadFile("plantilla.escpos")
	errores.PanicIfError(err)
	// Fusionar plantilla con estructura factura
	f, err := plantillas.MergeEscPosTemplate(
		"escpos",
		string(plantilla),
		factura,
		"/assets",
		formato.DMA,
		formato.EUR,
	)
	errores.PanicIfError(err)
	// Convertir la plantilla fusionada a fichero esc/pos sabor EPSON (pos) o sabor SEIKO (cajero)
	prn, mm, err := plantillas.GenerateEscPos(f, plantillas.EPSON)
	errores.PanicIfError(err)
	// Genera fichero PDF, debe ser igual si el esc/pos es sabor EPSON o sabor SEIKO
	err = plantillas.GenerateEscPosPdf(prn, "recibo.pdf", mm)
	errores.PanicIfError(err)
	fmt.Println("Generado fichero recibo.pdf")
	// Output: Generado fichero recibo.pdf
}
