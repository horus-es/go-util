package formato

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// https://en.wikipedia.org/wiki/ISO_4217
type FormatoPrecio string

const (
	EUR FormatoPrecio = "EUR" // Euros
	USD FormatoPrecio = "USD" // Dólares USA
	COP FormatoPrecio = "COP" // Pesos colombianos
)

// Parsea un precio
func ParsePrecio(p string, fp FormatoPrecio) (result float64, err error) {
	s := strings.TrimSpace(p)
	s = strings.TrimSuffix(s, "€")
	s = strings.TrimPrefix(s, "€")
	s = strings.TrimSuffix(s, "$")
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimSuffix(s, string(EUR))
	s = strings.TrimPrefix(s, string(EUR))
	s = strings.TrimSuffix(s, string(USD))
	s = strings.TrimPrefix(s, string(USD))
	s = strings.TrimSuffix(s, string(COP))
	s = strings.TrimPrefix(s, string(COP))
	switch fp {
	case EUR:
		result, err = ParseNumero(s, ",")
	case USD:
		result, err = ParseNumero(s, ".")
	case COP:
		result, err = ParseNumero(s, ".")
	default:
		result, err = strconv.ParseFloat(s, 64)
	}
	if err != nil {
		err = fmt.Errorf("precio %q no reconocido", p)
	}
	return
}

// Imprime un precio
func PrintPrecio(v float64, fp FormatoPrecio) string {
	var result string
	switch fp {
	case EUR:
		result = PrintNumero(v, 2, ",", ".") + "€"
	case USD:
		result = "$" + PrintNumero(v, 2, ".", ",")
	case COP:
		result = "$" + PrintNumero(v, 0, ",", ".")
	default:
		result = fmt.Sprintf("%f", v)
	}
	return result
}

// Redondea un precio con el número de decimales de la moneda
func RedondeaPrecio(v float64, fp FormatoPrecio) float64 {
	var result float64
	switch fp {
	case EUR, USD:
		result = 100
	case COP:
		result = 1
	default:
		result = 1000000
	}
	result = math.Round(v*result) / result
	return result
}
