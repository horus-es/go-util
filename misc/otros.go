// Funciones misceláneas
package misc

import (
	"bytes"
	"runtime"
	"strconv"
	"unicode"

	"github.com/horus-es/go-util/v3/errores"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

// Transformer para quitar acentos
var quitaAcentosTransformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// QuitaAcentos quita acentos y diacríticos de un string, p.e. "canción" → "cancion"
func QuitaAcentos(s string) string {
	s, _, _ = transform.String(quitaAcentosTransformer, s)
	return s
}
