package parse

import (
	"fmt"
	"horus-es/go-util/errores"
	"strings"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestPrintNumero(t *testing.T) {
	assert.Equal(t, "12 345,68", PrintNumero(12345.6789, 2, ",", " "))
	assert.Equal(t, "-12 300", PrintNumero(-12345.6789, -2, ",", " "))
	assert.Equal(t, "123 456 790", PrintNumero(123456789, -1, ",", " "))
	assert.Equal(t, "1.234.567.890,100", PrintNumero(1234567890.1, 3, ",", "."))
	assert.Equal(t, "34.567.890,100", PrintNumero(34567890.1, 3, ",", "."))
}

func ExamplePrintNumero() {
	fmt.Println(PrintNumero(12345.6789, 2, ",", "."))
	fmt.Println(PrintNumero(12345.6789, 0, ",", "."))
	fmt.Println(PrintNumero(12345.6789, -2, ",", "."))
	// Output:
	// 12.345,68
	// 12.346
	// 12.300
}

func ExampleParseNumero() {
	fmt.Println(PrintNumero(12345.6789, 2, ",", "."))
	fmt.Println(PrintNumero(12345.6789, 0, ",", "."))
	fmt.Println(PrintNumero(12345.6789, -2, ",", "."))
	// Output:
	// 12.345,68
	// 12.346
	// 12.300
}

func compruebaParseLogica(t *testing.T, espera bool, valor string) {
	obtiene, err := ParseLogica(valor)
	if valor == "--FAIL--" {
		assert.NotNil(t, err)
		return
	}
	assert.Nil(t, err)
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
	f, err := ParseLogica("Si")
	errores.PanicIfError(err)
	fmt.Println(f)
	f, err = ParseLogica("F")
	errores.PanicIfError(err)
	fmt.Println(f)
	_, err = ParseLogica("ZERO")
	fmt.Println(err)
	// Output:
	// true
	// false
	// valor "ZERO" no reconocido
}

func compruebaParseOpcion(t *testing.T, espera string, opcion string, admitidas ...string) {
	obtiene, err := ParseOpcion(opcion, admitidas...)
	if espera == "--FAIL--" {
		assert.NotNil(t, err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, espera, obtiene)
}

func TestParseOpciones(t *testing.T) {
	compruebaParseOpcion(t, "Dos", "dos", "Uno", "Dos", "Tres")
	compruebaParseOpcion(t, "ERROR", "erroRes", "avisos->WARN", "errores->ERROR")
	compruebaParseOpcion(t, "3", "tres", "Uno->1", "Dos->2", "Tres->3")
	compruebaParseOpcion(t, "--FAIL--", "trez", "Uno->1", "Dos->2", "Tres->3")
}

func ExampleParseOpcion() {
	opt, err := ParseOpcion("tres", "Uno", "Dos", "Tres")
	errores.PanicIfError(err)
	fmt.Println(opt)
	opt, err = ParseOpcion("tres", "Uno->1", "Dos->2", "Tres->3")
	errores.PanicIfError(err)
	fmt.Println(opt)
	_, err = ParseOpcion("cuatro", "Uno", "Dos", "Tres")
	fmt.Println(err)
	// Output:
	// Tres
	// 3
	// opción "cuatro" no reconocida
}

func TestParseUUID(t *testing.T) {
	suuid := "8D2C5C10-62D6-4B90-b4Af-8C006883C648"
	uuid, err := ParseUUID(suuid)
	assert.Nil(t, err)
	assert.Equal(t, strings.ToLower(suuid), PrintUUID(uuid))
	uuid, err = ParseUUID("")
	assert.Nil(t, err)
	assert.Equal(t, pgtype.Null, uuid.Status)
	assert.Equal(t, "", PrintUUID(uuid))
	suuid = "error"
	_, err = ParseUUID(suuid)
	assert.NotNil(t, err)
}

func ExampleParseUUID() {
	uuid, err := ParseUUID("b55e7cec-7126-4a19-8ab1-86481ead2803")
	errores.PanicIfError(err)
	fmt.Println(PrintUUID(uuid))
	// Output: b55e7cec-7126-4a19-8ab1-86481ead2803
}

func ExampleMustParseUUID() {
	uuid := MustParseUUID("7e2a8034-c319-4dfa-a846-e2c176aba2e4")
	fmt.Println(PrintUUID(uuid))
	// Output: 7e2a8034-c319-4dfa-a846-e2c176aba2e4
}

func ExamplePrintUUID() {
	uuid, err := ParseUUID("7ca4d9eb-9b81-4aae-86d9-682d73f4138f")
	errores.PanicIfError(err)
	fmt.Println(PrintUUID(uuid))
	// Output: 7ca4d9eb-9b81-4aae-86d9-682d73f4138f
}
