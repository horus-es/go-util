package formato_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/stretchr/testify/assert"
)

func TestPrintNumero(t *testing.T) {
	assert.Equal(t, "12 345,68", formato.PrintNumero(12345.6789, 2, ",", " "))
	assert.Equal(t, "-12 300", formato.PrintNumero(-12345.6789, -2, ",", " "))
	assert.Equal(t, "123 456 790", formato.PrintNumero(123456789, -1, ",", " "))
	assert.Equal(t, "1.234.567.890,100", formato.PrintNumero(1234567890.1, 3, ",", "."))
	assert.Equal(t, "34.567.890,100", formato.PrintNumero(34567890.1, 3, ",", "."))
}

func ExamplePrintNumero() {
	fmt.Println(formato.PrintNumero(12345.6789, 2, ",", "."))
	fmt.Println(formato.PrintNumero(12345.6789, 0, ",", "."))
	fmt.Println(formato.PrintNumero(12345.6789, -2, ",", "."))
	// Output:
	// 12.345,68
	// 12.346
	// 12.300
}

func ExampleParseNumero() {
	fmt.Println(formato.PrintNumero(12345.6789, 2, ",", "."))
	fmt.Println(formato.PrintNumero(12345.6789, 0, ",", "."))
	fmt.Println(formato.PrintNumero(12345.6789, -2, ",", "."))
	// Output:
	// 12.345,68
	// 12.346
	// 12.300
}

func compruebaParseLogica(t *testing.T, espera bool, valor string) {
	obtiene, err := formato.ParseLogica(valor)
	if valor == "--FAIL--" {
		assert.NotNil(t, err)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, espera, obtiene)
}

func TestParseLogica(t *testing.T) {
	compruebaParseLogica(t, true, "1")
	compruebaParseLogica(t, true, "TRUE")
	compruebaParseLogica(t, true, "S")
	compruebaParseLogica(t, false, "N")
	compruebaParseLogica(t, false, "False")
	compruebaParseLogica(t, false, "--FAIL--")
}

func ExampleParseLogica() {
	f, err := formato.ParseLogica("Si")
	errores.PanicIfError(err)
	fmt.Println(f)
	f, err = formato.ParseLogica("F")
	errores.PanicIfError(err)
	fmt.Println(f)
	_, err = formato.ParseLogica("ZERO")
	fmt.Println(err)
	// Output:
	// true
	// false
	// valor "ZERO" no reconocido
}

func compruebaParseOpcion(t *testing.T, espera string, opcion string, admitidas ...string) {
	obtiene, err := formato.ParseOpcion(opcion, admitidas...)
	if espera == "--FAIL--" {
		assert.NotNil(t, err)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, espera, obtiene)
}

func TestParseOpciones(t *testing.T) {
	compruebaParseOpcion(t, "Dos", "dos", "Uno", "Dos", "Tres")
	compruebaParseOpcion(t, "ERROR", "erroRes", "avisos->WARN", "errores->ERROR")
	compruebaParseOpcion(t, "3", "tres", "Uno->1", "Dos->2", "Tres->3")
	compruebaParseOpcion(t, "--FAIL--", "trez", "Uno->1", "Dos->2", "Tres->3")
}

func ExampleParseOpcion() {
	opt, err := formato.ParseOpcion("tres", "Uno", "Dos", "Tres")
	errores.PanicIfError(err)
	fmt.Println(opt)
	opt, err = formato.ParseOpcion("tres", "Uno->1", "Dos->2", "Tres->3")
	errores.PanicIfError(err)
	fmt.Println(opt)
	_, err = formato.ParseOpcion("cuatro", "Uno", "Dos", "Tres")
	fmt.Println(err)
	// Output:
	// Tres
	// 3
	// opción "cuatro" no reconocida
}

func TestParseUUID(t *testing.T) {
	suuid := "8D2C5C10-62D6-4B90-b4Af-8C006883C648"
	uuid, err := formato.ParseUUID(suuid)
	assert.NoError(t, err)
	assert.Equal(t, strings.ToLower(suuid), formato.PrintUUID(uuid))
	uuid, err = formato.ParseUUID("")
	assert.NoError(t, err)
	assert.False(t, uuid.Valid)
	assert.Equal(t, "", formato.PrintUUID(uuid))
	suuid = "error"
	_, err = formato.ParseUUID(suuid)
	assert.NotNil(t, err)
}

func ExampleParseUUID() {
	uuid, err := formato.ParseUUID("b55e7cec-7126-4a19-8ab1-86481ead2803")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintUUID(uuid))
	// Output: b55e7cec-7126-4a19-8ab1-86481ead2803
}

func ExampleMustParseUUID() {
	uuid := formato.MustParseUUID("7e2a8034-c319-4dfa-a846-e2c176aba2e4")
	fmt.Println(formato.PrintUUID(uuid))
	// Output: 7e2a8034-c319-4dfa-a846-e2c176aba2e4
}

func ExamplePrintUUID() {
	uuid, err := formato.ParseUUID("7ca4d9eb-9b81-4aae-86d9-682d73f4138f")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintUUID(uuid))
	// Output: 7ca4d9eb-9b81-4aae-86d9-682d73f4138f
}
