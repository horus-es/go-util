// Funciones de conversión de fechas, precios, intervalos, números, uuids, lógicas y opciones a textos y viceversa
package formato

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/jackc/pgtype"
)

// Parsea un numero, considerando que puede incluir separadores
func ParseNumero(v, sepDecimal string) (result float64, err error) {
	s := strings.ReplaceAll(v, " ", "")
	if sepDecimal == "." {
		s = strings.ReplaceAll(s, ",", "")
	} else {
		s = strings.ReplaceAll(s, ".", "")
		s = strings.Replace(s, sepDecimal, ".", 1)
	}
	result, err = strconv.ParseFloat(s, 64)
	if err != nil {
		err = fmt.Errorf("número %q no reconocido", v)
	}
	return
}

// Imprime un número con decimales, usando los separadores especificados.
// Si los decimales son negativos, se ponen a 0 las cifras menos significativas.
// PrintNumero(12345.6789,2,","," ") = "12 345,68"
// PrintNumero(-12345.6789,-2,","," ") = "-12 300"
// Ver https://en.wikipedia.org/wiki/Decimal_separator.
func PrintNumero(v float64, decimales int, sepDecimal, sepMiles string) string {
	var s string
	if decimales >= 0 {
		s = strconv.FormatFloat(v, 'f', decimales, 64)
	} else {
		var n int64 = 1
		for ; decimales < 0; decimales++ {
			n *= 10
		}
		fn := float64(n)
		s = strconv.FormatFloat(math.Round(v/fn)*fn, 'f', 0, 64)
	}
	l := len(s)
	if decimales > 0 {
		l -= decimales + 1
	}
	var result strings.Builder
	z0 := 0
	if s[0] == '-' {
		z0 = 1
	}
	for z := 0; z < l; z++ {
		if z > z0 && (l-z)%3 == 0 {
			result.WriteString(sepMiles)
		}
		result.WriteByte(s[z])
	}
	if decimales > 0 {
		result.WriteString(sepDecimal)
		result.WriteString(s[l+1:])
	}
	return result.String()
}

// Parsea una variable lógica, soporta: true, false, yes, si, no, 1 y 0
func ParseLogica(s string) (bool, error) {
	if strings.EqualFold(s, "true") ||
		strings.EqualFold(s, "t") ||
		strings.EqualFold(s, "yes") ||
		strings.EqualFold(s, "y") ||
		strings.EqualFold(s, "si") ||
		strings.EqualFold(s, "s") ||
		strings.EqualFold(s, "1") {
		return true, nil
	}
	if strings.EqualFold(s, "false") ||
		strings.EqualFold(s, "f") ||
		strings.EqualFold(s, "no") ||
		strings.EqualFold(s, "n") ||
		strings.EqualFold(s, "0") {
		return false, nil
	}
	return false, fmt.Errorf("valor %q no reconocido", s)
}

// Se asegura de que la opción sea alguna de las admitidas, en mayúsculas o minúsculas.
// Opcionalmente se puede hacer una conversion si admitido tiene el formato "opcion->convertido"
// La opción por defecto se puede indicar con el formato "->opcion"
func ParseOpcion(opcion string, admitidas ...string) (string, error) {
	for _, admitido := range admitidas {
		k, v, ok := strings.Cut(admitido, "->")
		if ok {
			// con conversion
			if strings.EqualFold(opcion, k) {
				return v, nil
			}
		} else {
			if strings.EqualFold(opcion, admitido) {
				return admitido, nil
			}
		}
	}
	return "", fmt.Errorf("opción %q no reconocida", opcion)
}

// Parsea un objeto a UUID, los vacíos se consideran NULL
func ParseUUID(uuid string) (result pgtype.UUID, err error) {
	if uuid == "" {
		result.Status = pgtype.Null
		return
	}
	err = result.Set(uuid)
	return
}

// Parsea un objeto a UUID, los vacíos se consideran NULL, panic si error
func MustParseUUID(uuid string) pgtype.UUID {
	result, err := ParseUUID(uuid)
	errores.PanicIfError(err)
	return result
}

// Imprime un UUID, los NULL se imprimen como vacíos
func PrintUUID(uuid pgtype.UUID) string {
	if uuid.Status == pgtype.Null {
		return ""
	}
	var result string
	err := uuid.AssignTo(&result)
	errores.PanicIfError(err)
	return result
}
