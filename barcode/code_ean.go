package barcode

import (
	"fmt"
)

func barcodeEAN13(code string) (bars, hri string, err error) {
	const startend = "111"            // patrón inicial y final
	const middle = "11111"            // patrón medio
	var patterns = map[byte][]string{ // patrones[paridad]
		'O': {
			"3211", "2221", "2122", "1411", "1132",
			"1231", "1114", "1312", "1213", "3112",
		},
		'E': {
			"1123", "1222", "2212", "1141", "2311",
			"1321", "4111", "2131", "3121", "2113",
		},
	}
	var parities = []string{ // paridades
		"OOOOOOOOOOOO",
		"OOEOEEOOOOOO",
		"OOEEOEOOOOOO",
		"OOEEEOOOOOOO",
		"OEOOEEOOOOOO",
		"OEEOOEOOOOOO",
		"OEEEOOOOOOOO",
		"OEOEOEOOOOOO",
		"OEOEEOOOOOOO",
		"OEEOEOOOOOOO",
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return "", "", fmt.Errorf("invalid char 0x%02X", r)
		}
	}
	switch len(code) {
	case 12:
		// Añadimos el dígito de control
		sum := ean13Sum(code)
		code = code + string(rune(sum+48))
	case 13:
		// Validamos el dígito de control
		sum := ean13Sum(code)
		chk := int(code[12]) - 48
		if chk != sum {
			return "", "", fmt.Errorf("invalid check digit %d", chk)
		}
	default:
		return "", "", fmt.Errorf("invalid length %d", len(code))
	}
	// Hallamos la secuencia
	nsc := int(code[0]) - 48
	if nsc > 0 {
		hri += string(code[0]) + " "
	}
	par := parities[nsc]
	bars = startend
	for i := 1; i < 7; i++ {
		n := int(code[i]) - 48
		p := par[i-1]
		bars += patterns[p][n]
		hri += string(code[i])
	}
	bars += middle
	hri += " "
	for i := 7; i < 13; i++ {
		n := int(code[i]) - 48
		p := par[i-1]
		bars += patterns[p][n]
		hri += string(code[i])
	}
	bars += startend
	return
}

// Halla la suma de control
func ean13Sum(code string) int {
	sum := 0
	for i := 1; i < 12; i += 2 {
		sum += int(code[i]) - 48
	}
	sum *= 3
	for i := 0; i < 11; i += 2 {
		sum += int(code[i]) - 48
	}
	sum %= 10
	if sum > 0 {
		sum = 10 - sum
	}
	return sum
}

func barcodeEAN8(code string) (bars, hri string, err error) {
	const startend = "111"   // patrón inicial y final
	const middle = "11111"   // patrón medio
	var patterns = []string{ // patrones
		"3211", "2221", "2122", "1411", "1132",
		"1231", "1114", "1312", "1213", "3112",
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return "", "", fmt.Errorf("invalid char 0x%02X", r)
		}
	}
	switch len(code) {
	case 7:
		// Añadimos el dígito de control
		sum := ean8Sum(code)
		code = code + string(rune(sum+48))
	case 8:
		// Validamos el dígito de control
		sum := ean8Sum(code)
		chk := int(code[7]) - 48
		if chk != sum {
			return "", "", fmt.Errorf("invalid check digit %d", chk)
		}
	default:
		return "", "", fmt.Errorf("invalid length %d", len(code))
	}
	// Hallamos la secuencia
	bars = startend
	for i := 0; i < 4; i++ {
		n := int(code[i]) - 48
		bars += patterns[n]
		hri += string(code[i])
	}
	bars += middle
	hri += " "
	for i := 4; i < 8; i++ {
		n := int(code[i]) - 48
		bars += patterns[n]
		hri += string(code[i])
	}
	bars += startend
	return
}

// Halla la suma de control
func ean8Sum(code string) int {
	sum := 0
	for i := 0; i < 7; i += 2 {
		sum += int(code[i]) - 48
	}
	sum *= 3
	for i := 1; i < 6; i += 2 {
		sum += int(code[i]) - 48
	}
	sum %= 10
	if sum > 0 {
		sum = 10 - sum
	}
	return sum
}
