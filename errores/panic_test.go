package errores

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicIfError(t *testing.T) {
	assert.Panics(t, func() {
		_, err := strconv.Atoi("4o")
		PanicIfError(err)
	})
	assert.Panics(t, func() {
		_, err := strconv.Atoi("4l")
		PanicIfError(err, "Error de conversión")
	})
	assert.Panics(t, func() {
		_, err := strconv.Atoi("4z")
		PanicIfError(err, "Error de conversión %d", 42)
	})
}

func ExamplePanicIfError() {
	s := "42"
	k, err := strconv.Atoi(s)
	PanicIfError(err, "La cadena %q no se puede convertir a entero", s)
	fmt.Println(k)
	// Output: 42
}

func TestPanicIfTrue(t *testing.T) {
	assert.Panics(t, func() {
		PanicIfTrue(true, "Ciertamente")
	})
	assert.Panics(t, func() {
		PanicIfTrue(true, "Ciertamente %s", "es un error")
	})
}

func ExamplePanicIfTrue() {
	z := 2 * 3
	PanicIfTrue(z != 6, "Error en multiplicación: 2x3=%d", z)
	fmt.Println(z)
	// Output: 6
}
