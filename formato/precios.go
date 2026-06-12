package formato

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/horus-es/go-util/v3/errores"
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
	ALZA     string = "ALZA"     // Redondeo al alza
	ESTANDAR string = "ESTANDAR" // Redondeo al valor mas cercano
	BAJA     string = "BAJA"     // Redondeo a la baja
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

const DECIMALES_DEFECTO = 1000

// Imprime un precio. Si decimales=DEFAULT use sa el valor por defecto de la moneda.
func PrintPrecio(v float64, fp Moneda, decimales int) string {
	var result string
	switch fp {
	case EUR:
		if decimales == DECIMALES_DEFECTO {
			decimales = 2
		}
		result = PrintNumero(v, decimales, ",", ".") + " €"
	case USD:
		if decimales == DECIMALES_DEFECTO {
			decimales = 2
		}
		if v >= 0 {
			result = "$" + PrintNumero(v, decimales, ".", ",")
		} else {
			result = "-$" + PrintNumero(-v, decimales, ".", ",")
		}
	case COP, MXN:
		if decimales == DECIMALES_DEFECTO {
			decimales = 0
		}
		if v >= 0 {
			result = "$" + PrintNumero(v, decimales, ",", ".")
		} else {
			result = "-$" + PrintNumero(-v, decimales, ",", ".")
		}
	default:
		result = fmt.Sprintf("%f", v)
	}
	return result
}

// Redondea un precio (p) usando la unidad monetaria mínima (umm).
// El redondeo puede ser al ALZA, a la BAJA, o al valor mas cercano(ESTANDAR).
// Los precios negativos se redondean igual que los positivos (simetrico respecto a 0).
// Si umm=0 no se redondea
func RedondeaPrecio(p, um float64, tipo string) float64 {
	if um == 0 {
		return p
	}
	if um < 0 {
		errores.PanicIfTrue(um < 0, "unidad monetaria negativa: %f", um)
	}
	if p < 0 {
		return -RedondeaPrecio(-p, um, tipo)
	}
	const epsilon = 0.00001
	mu := 1 / um
	switch tipo {
	case ALZA:
		return math.Ceil(mu*p-epsilon) / mu
	case BAJA:
		return math.Floor(mu*p+epsilon) / mu
	case ESTANDAR:
		return math.Round(mu*p+epsilon) / mu
	default:
		errores.PanicIfTrue(true, "tipo no soportado: %s", tipo)
		return p
	}
}
