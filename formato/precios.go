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

const (
	ALZA  byte = '+' // Redondeo al alza
	JUSTO byte = '~' // Redondeo al valor mas cercano
	BAJA  byte = '-' // Redondeo a la baja
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
		result = PrintNumero(v, 2, ",", ".") + " €"
	case USD:
		if v >= 0 {
			result = "$" + PrintNumero(v, 2, ".", ",")
		} else {
			result = "-$" + PrintNumero(-v, 2, ".", ",")
		}
	case COP, MXN:
		if v >= 0 {
			result = "$" + PrintNumero(v, 0, ",", ".")
		} else {
			result = "-$" + PrintNumero(-v, 0, ",", ".")
		}
	default:
		result = fmt.Sprintf("%f", v)
	}
	return result
}

// Redondea un precio (p) usando la unidad monetaria mínima (umm). El redondeo puede ser al ALZA, a la BAJA, o al valor mas cercano(JUSTO). Los negativos se redondean igual que los positivos.
func RedondeaPrecio(p, umm float64, tipo byte) float64 {
	const epsilon = 0.00001
	if p < 0 {
		return -RedondeaPrecio(-p, umm, tipo)
	}
	mmu := 1 / umm
	switch tipo {
	case ALZA:
		return math.Ceil(mmu*p-epsilon) / mmu
	case BAJA:
		return math.Floor(mmu*p+epsilon) / mmu
	case JUSTO:
		return math.Round(mmu*p+epsilon) / mmu
	}
	panic("tipo no soportado: " + string(tipo))
}
