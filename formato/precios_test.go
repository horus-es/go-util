package formato

import (
	"fmt"
	"testing"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/stretchr/testify/assert"
)

func TestParsePrecio(t *testing.T) {
	var p float64
	p, _ = ParsePrecio(" 1 234,56 € ", EUR)
	assert.Equal(t, 1234.56, p, EUR)
	p, _ = ParsePrecio(" EUR 1 234,56 € ", EUR)
	assert.Equal(t, 1234.56, p, EUR)
	p, _ = ParsePrecio("1234", EUR)
	assert.Equal(t, 1234.00, p, EUR)
	p, _ = ParsePrecio("$1,234.56$", COP)
	assert.Equal(t, 1234.56, p, EUR)
	p, _ = ParsePrecio("$1.234,56$", EUR)
	assert.Equal(t, 1234.56, p, EUR)
	p, _ = ParsePrecio("$1234.56$", "OTRO")
	assert.Equal(t, 1234.56, p, "OTRO")
	_, err := ParsePrecio("1$234,56", USD)
	assert.NotNil(t, err)
}

func ExampleParsePrecio() {
	p, err := ParsePrecio(" 1 234,56 € ", EUR)
	errores.PanicIfError(err)
	fmt.Println(p)
	// Output: 1234.56
}

func TestPrintPrecio(t *testing.T) {
	assert.Equal(t, "12.345,68€", PrintPrecio(12345.6789, EUR), EUR)
	assert.Equal(t, "$12,345.68", PrintPrecio(12345.6789, USD), USD)
	assert.Equal(t, "$12.346", PrintPrecio(12345.6789, COP), COP)
	assert.Equal(t, "-12.345,68€", PrintPrecio(-12345.6789, EUR), EUR)
	assert.Equal(t, "345,68€", PrintPrecio(345.6789, EUR), EUR)
	assert.Equal(t, "-345,68€", PrintPrecio(-345.6789, EUR), EUR)
	assert.Equal(t, "-123.345,68€", PrintPrecio(-123345.6789, EUR), EUR)
	assert.Equal(t, "123345.678900", PrintPrecio(123345.6789, "OTRO"), "OTRO")
}

func ExamplePrintPrecio() {
	fmt.Println(PrintPrecio(12345.6789, EUR))
	// Output: 12.345,68€
}

func TestRedondeaPrecio(t *testing.T) {
	assert.Equal(t, 12345.68, RedondeaPrecio(12345.6789, EUR))
	assert.Equal(t, 12345.68, RedondeaPrecio(12345.6789, USD), USD)
	assert.Equal(t, 12346.00, RedondeaPrecio(12345.6789, COP), COP)
	assert.Equal(t, 12345.6789, RedondeaPrecio(12345.6789, "OTRO"), "OTRO")
}

func ExampleRedondeaPrecio() {
	fmt.Println(RedondeaPrecio(12345.6789, EUR))
	// Output: 12345.68
}
