package formato

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// https://en.wikipedia.org/wiki/ISO_4217
type Moneda string

const (
	EUR Moneda = "EUR" // Euros
	USD Moneda = "USD" // Dólares USA
	COP Moneda = "COP" // Pesos colombianos
	MXN Moneda = "MXN" // Pesos mexicanos
)

// Parsea un precio
func ParsePrecio(p string, fp Moneda) (result float64, err error) {
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
	s = strings.TrimSuffix(s, string(MXN))
	s = strings.TrimPrefix(s, string(MXN))
	switch fp {
	case EUR:
		result, err = ParseNumero(s, ",")
	case USD:
		result, err = ParseNumero(s, ".")
	case COP, MXN:
		result, err = ParseNumero(s, ",")
	default:
		result, err = strconv.ParseFloat(s, 64)
	}
	if err != nil {
		err = fmt.Errorf("precio %q no reconocido", p)
	}
	return
}

// Imprime un precio
func PrintPrecio(v float64, fp Moneda) string {
	var result string
	switch fp {
	case EUR:
		result = PrintNumero(v, 2, ",", ".") + "€"
	case USD:
		result = "$" + PrintNumero(v, 2, ".", ",")
	case COP, MXN:
		result = "$" + PrintNumero(v, 0, ",", ".")
	default:
		result = fmt.Sprintf("%f", v)
	}
	return result
}

// Redondea un precio con el número de decimales de la moneda
func RedondeaPrecio(v float64, fp Moneda) float64 {
	var result float64
	switch fp {
	case EUR, USD:
		result = 100
	case COP, MXN:
		result = 1
	default:
		result = 1000000
	}
	result = math.Round(v*result) / result
	return result
}
