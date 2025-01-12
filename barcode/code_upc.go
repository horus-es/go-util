package barcode

import (
	"errors"
	"fmt"
)

func barcodeUPCA(code string) (string, error) {
	const startEnd = "111"
	const middle = "11111"
	var patterns = []string{
		"3211", "2221", "2122", "1411", "1132",
		"1231", "1114", "1312", "1213", "3112"}
	for _, r := range code {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("invalid char %c", r)
		}
	}
	// Calculamos el dígito de control
	switch len(code) {
	case 11:
		// Añadimos el dígito de control
		sum := upcSum(code)
		code = code + string(rune(sum+48))
	case 12:
		// Validamos el dígito de control
		sum := upcSum(code)
		chk := int(code[11]) - 48
		if chk != sum {
			return "", fmt.Errorf("invalid check digit %d", chk)
		}
	default:
		return "", fmt.Errorf("invalid length %d", len(code))
	}
	// Hallamos la secuencia
	barpattern := startEnd
	for i := 0; i < 6; i++ {
		barpattern += patterns[code[i]-48]
	}
	barpattern += middle
	for i := 6; i < 12; i++ {
		barpattern += patterns[code[i]-48]
	}
	barpattern += startEnd
	return barpattern, nil
}

func barcodeUPCE(code string) (string, error) {
	const start = "111"               // patrón inicial
	const end = "111111"              // patrón final
	var patterns = map[byte][]string{ // patrones[paridad]
		'E': {
			"1123", "1222", "2212", "1141", "2311",
			"1321", "4111", "2131", "3121", "2113",
		},
		'O': {
			"3211", "2221", "2122", "1411", "1132",
			"1231", "1114", "1312", "1213", "3112",
		},
	}
	var parities = map[int][]string{ // paridades[nsc]
		0: {
			"EEEOOO", "EEOEOO", "EEOOEO", "EEOOOE", "EOEEOO",
			"EOOEEO", "EOOOEE", "EOEOEO", "EOEOOE", "EOOEOE",
		},
		1: {
			"OOOEEE", "OOEOEE", "OOEEOE", "OOEEEO", "OEOOEE",
			"OEEOOE", "OEEEOO", "OEOEOE", "OEOEEO", "OEEOEO",
		},
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("invalid char %c", r)
		}
	}
	switch len(code) {
	case 6:
		// Agregamos nsc=0 y el dígito de control
		code = "0" + code + "X"
		upca, err := upcExpand(code)
		if err != nil {
			return "", err
		}
		sum := upcSum(upca)
		code = code[:7] + string(rune(sum+48))
	case 7:
		// Agregamos el dígito de control
		code = code + "X"
		upca, err := upcExpand(code)
		if err != nil {
			return "", err
		}
		sum := upcSum(upca)
		code = code[:7] + string(rune(sum+48))
	case 8:
		// Validamos el dígito de control
		upca, err := upcExpand(code)
		if err != nil {
			return "", err
		}
		sum := upcSum(upca)
		chk := int(code[7]) - 48
		if chk != sum {
			return "", fmt.Errorf("invalid check digit %d", chk)
		}
	case 11:
		// Agregamos el dígito de control y comprimimos
		code = code + "X"
		sum := upcSum(code)
		var err error
		code, err = upcCompress(code)
		if err != nil {
			return "", err
		}
		code = code[:7] + string(rune(sum+48))
	case 12:
		// Validamos el dígito de control y comprimimos
		sum := upcSum(code)
		chk := int(code[11]) - 48
		if chk != sum {
			return "", fmt.Errorf("invalid check digit %d", chk)
		}
		var err error
		code, err = upcCompress(code)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid length %d", len(code))
	}
	// Hallamos la secuencia
	nsc := int(code[0]) - 48
	if nsc != 0 && nsc != 1 {
		return "", fmt.Errorf("invalid nsc %d", nsc)
	}
	sum := int(code[7]) - 48
	par := parities[nsc][sum]
	barpattern := start
	for i := 1; i < 7; i++ {
		n := int(code[i]) - 48
		p := par[i-1]
		barpattern += patterns[p][n]
	}
	barpattern += end
	return barpattern, nil
}

// Halla la suma de control
func upcSum(code string) int {
	sum := 0
	for i := 0; i < 11; i += 2 {
		sum += int(code[i]) - 48
	}
	sum *= 3
	for i := 1; i < 10; i += 2 {
		sum += int(code[i]) - 48
	}
	sum %= 10
	if sum > 0 {
		sum = 10 - sum
	}
	return sum
}

// Convierte UPC-E a UPC-A
func upcExpand(code string) (string, error) {
	switch code[6] {
	case '0', '1', '2':
		return code[0:3] + code[6:7] + "0000" + code[3:6] + code[7:8], nil
	case '3':
		return code[0:4] + "00000" + code[4:6] + code[7:8], nil
	case '4':
		return code[0:5] + "00000" + code[5:6] + code[7:8], nil
	case '5', '6', '7', '8', '9':
		return code[0:6] + "0000" + code[6:8], nil
	}
	return "", errors.New("expand error")
}

// Convierte UPC-A a UPC-E
func upcCompress(code string) (string, error) {
	if code[4] == '0' && code[5] == '0' && code[6] == '0' && code[7] == '0' && code[3] >= '0' && code[3] <= '2' {
		return code[0:3] + code[8:11] + code[3:4] + code[11:12], nil
	}
	if code[4] == '0' && code[5] == '0' && code[6] == '0' && code[7] == '0' && code[8] == '0' {
		return code[0:4] + code[9:11] + "3" + code[11:12], nil
	}
	if code[5] == '0' && code[6] == '0' && code[7] == '0' && code[8] == '0' && code[9] == '0' {
		return code[0:5] + code[10:11] + "4" + code[11:12], nil
	}
	if code[6] == '0' && code[7] == '0' && code[8] == '0' && code[9] == '0' && code[10] >= '5' && code[10] <= '9' {
		return code[0:6] + code[10:12], nil
	}
	return "", errors.New("compress error")
}
