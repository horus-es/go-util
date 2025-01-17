package formato_test

import (
	"fmt"
	"testing"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/stretchr/testify/assert"
)

func TestParsePrecio(t *testing.T) {
	var p float64
	p, _ = formato.ParsePrecio(" 1 234,56 € ", formato.EUR)
	assert.Equal(t, 1234.56, p, formato.EUR)
	p, _ = formato.ParsePrecio(" EUR 1 234,56 € ", formato.EUR)
	assert.Equal(t, 1234.56, p, formato.EUR)
	p, _ = formato.ParsePrecio("1234", formato.EUR)
	assert.Equal(t, 1234.00, p, formato.EUR)
	p, _ = formato.ParsePrecio("$1,234.56$", formato.COP)
	assert.Equal(t, 1234.56, p, formato.EUR)
	p, _ = formato.ParsePrecio("$1.234,56$", formato.EUR)
	assert.Equal(t, 1234.56, p, formato.EUR)
	p, _ = formato.ParsePrecio("$1234.56$", "OTRO")
	assert.Equal(t, 1234.56, p, "OTRO")
	_, err := formato.ParsePrecio("1$234,56", formato.USD)
	assert.NotNil(t, err)
}

func ExampleParsePrecio() {
	p, err := formato.ParsePrecio(" 1 234,56 € ", formato.EUR)
	errores.PanicIfError(err)
	fmt.Println(p)
	// Output: 1234.56
}

func TestPrintPrecio(t *testing.T) {
	assert.Equal(t, "12.345,68 €", formato.PrintPrecio(12345.6789, formato.EUR), formato.EUR)
	assert.Equal(t, "$12,345.68", formato.PrintPrecio(12345.6789, formato.USD), formato.USD)
	assert.Equal(t, "$12.346", formato.PrintPrecio(12345.6789, formato.COP), formato.COP)
	assert.Equal(t, "$12.346", formato.PrintPrecio(12345.6789, formato.MXN), formato.MXN)
	assert.Equal(t, "-12.345,68 €", formato.PrintPrecio(-12345.6789, formato.EUR), formato.EUR)
	assert.Equal(t, "345,68 €", formato.PrintPrecio(345.6789, formato.EUR), formato.EUR)
	assert.Equal(t, "-345,68 €", formato.PrintPrecio(-345.6789, formato.EUR), formato.EUR)
	assert.Equal(t, "-123.345,68 €", formato.PrintPrecio(-123345.6789, formato.EUR), formato.EUR)
	assert.Equal(t, "123345.678900", formato.PrintPrecio(123345.6789, "OTRO"), "OTRO")
}

func ExamplePrintPrecio() {
	fmt.Println(formato.PrintPrecio(12345.6789, formato.EUR))
	// Output: 12.345,68 €
}

func TestRedondeaPrecio(t *testing.T) {
	assert.Equal(t, 12345.68, formato.RedondeaPrecio(12345.6789, 0.01), formato.EUR)
	assert.Equal(t, 12345.68, formato.RedondeaPrecio(12345.6789, 0.01), formato.USD)
	assert.Equal(t, 12350.00, formato.RedondeaPrecio(12345.6789, 50), formato.COP)
}

func ExampleRedondeaPrecio() {
	fmt.Println(formato.RedondeaPrecio(12345.6789, 50))
	// Output: 12350
}
