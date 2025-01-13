// Funciones misceláneas
package misc

import (
	"bytes"
	"runtime"
	"strconv"

	"github.com/horus-es/go-util/v2/errores"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Hace lo mismo que la deprecada strings.Title: pasar la primera letra a mayúsculas
func Title(s string) string {
	return cases.Title(language.Und, cases.NoLower).String(s)
}

// Obtiene el ID de la goroutine actual
func GetGID() int64 {
	var s [64]byte
	b := s[:runtime.Stack(s[:], false)]
	b = b[len("goroutine "):]
	b = b[:bytes.IndexByte(b, ' ')]
	gid, err := strconv.ParseInt(string(b), 10, 64)
	errores.PanicIfError(err)
	return gid
}
