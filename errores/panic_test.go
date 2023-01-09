package errores

import (
	"fmt"
	"strconv"
)

func ExamplePanicIfError() {
	s := "42"
	k, err := strconv.Atoi(s)
	PanicIfError(err, "La cadena %q no se puede convertir a entero", s)
	fmt.Println(k)
	// Output: 42
}

func ExamplePanicIfTrue() {
	z := 2 * 3
	PanicIfTrue(z != 6, "Error en multiplicación: 2x3=%d", z)
	fmt.Println(6)
	// Output: 6
}
