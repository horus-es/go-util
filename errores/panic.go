// Gestión de panics
package errores

import (
	"fmt"

	"github.com/pkg/errors"
)

// Panic si error
func PanicIfError(err error, params ...any) {
	if err == nil {
		return
	}
	if len(params) == 0 {
		panic(err)
	}
	msg := fmt.Sprint(params[0])
	if len(params) == 1 {
		panic(errors.WithMessage(err, msg))
	}
	panic(errors.WithMessagef(err, msg, params[1:]...))
}

// Panic si condición
func PanicIfTrue(condition bool, msg string, params ...any) {
	if !condition {
		return
	}
	if len(params) > 0 {
		panic(errors.Errorf(msg, params...))
	} else {
		panic(errors.New(msg))
	}
}
