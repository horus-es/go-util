package errores_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/stretchr/testify/assert"
)

func TestPanicIfError(t *testing.T) {
	assert.Panics(t, func() {
		_, err := strconv.Atoi("4o")
		errores.PanicIfError(err)
	})
	assert.Panics(t, func() {
		_, err := strconv.Atoi("4l")
		errores.PanicIfError(err, "Error de conversión")
	})
	assert.Panics(t, func() {
		_, err := strconv.Atoi("4z")
		errores.PanicIfError(err, "Error de conversión %d", 42)
	})
}

func ExamplePanicIfError() {
	s := "42"
	k, err := strconv.Atoi(s)
	errores.PanicIfError(err, "La cadena %q no se puede convertir a entero", s)
	fmt.Println(k)
	// Output: 42
}

func TestPanicIfTrue(t *testing.T) {
	assert.Panics(t, func() {
		errores.PanicIfTrue(true, "Ciertamente")
	})
	assert.Panics(t, func() {
		errores.PanicIfTrue(true, "Ciertamente %s", "es un error")
	})
}

func ExamplePanicIfTrue() {
	z := 2 * 3
	errores.PanicIfTrue(z != 6, "Error en multiplicación: 2x3=%d", z)
	fmt.Println(z)
	// Output: 6
}
