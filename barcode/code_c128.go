package barcode

import (
	"fmt"
)

func barcodeC128(code string, variant byte) (string, error) {

	patterns := []string{
		"212222",  /* 00 */
		"222122",  /* 01 */
		"222221",  /* 02 */
		"121223",  /* 03 */
		"121322",  /* 04 */
		"131222",  /* 05 */
		"122213",  /* 06 */
		"122312",  /* 07 */
		"132212",  /* 08 */
		"221213",  /* 09 */
		"221312",  /* 10 */
		"231212",  /* 11 */
		"112232",  /* 12 */
		"122132",  /* 13 */
		"122231",  /* 14 */
		"113222",  /* 15 */
		"123122",  /* 16 */
		"123221",  /* 17 */
		"223211",  /* 18 */
		"221132",  /* 19 */
		"221231",  /* 20 */
		"213212",  /* 21 */
		"223112",  /* 22 */
		"312131",  /* 23 */
		"311222",  /* 24 */
		"321122",  /* 25 */
		"321221",  /* 26 */
		"312212",  /* 27 */
		"322112",  /* 28 */
		"322211",  /* 29 */
		"212123",  /* 30 */
		"212321",  /* 31 */
		"232121",  /* 32 */
		"111323",  /* 33 */
		"131123",  /* 34 */
		"131321",  /* 35 */
		"112313",  /* 36 */
		"132113",  /* 37 */
		"132311",  /* 38 */
		"211313",  /* 39 */
		"231113",  /* 40 */
		"231311",  /* 41 */
		"112133",  /* 42 */
		"112331",  /* 43 */
		"132131",  /* 44 */
		"113123",  /* 45 */
		"113321",  /* 46 */
		"133121",  /* 47 */
		"313121",  /* 48 */
		"211331",  /* 49 */
		"231131",  /* 50 */
		"213113",  /* 51 */
		"213311",  /* 52 */
		"213131",  /* 53 */
		"311123",  /* 54 */
		"311321",  /* 55 */
		"331121",  /* 56 */
		"312113",  /* 57 */
		"312311",  /* 58 */
		"332111",  /* 59 */
		"314111",  /* 60 */
		"221411",  /* 61 */
		"431111",  /* 62 */
		"111224",  /* 63 */
		"111422",  /* 64 */
		"121124",  /* 65 */
		"121421",  /* 66 */
		"141122",  /* 67 */
		"141221",  /* 68 */
		"112214",  /* 69 */
		"112412",  /* 70 */
		"122114",  /* 71 */
		"122411",  /* 72 */
		"142112",  /* 73 */
		"142211",  /* 74 */
		"241211",  /* 75 */
		"221114",  /* 76 */
		"413111",  /* 77 */
		"241112",  /* 78 */
		"134111",  /* 79 */
		"111242",  /* 80 */
		"121142",  /* 81 */
		"121241",  /* 82 */
		"114212",  /* 83 */
		"124112",  /* 84 */
		"124211",  /* 85 */
		"411212",  /* 86 */
		"421112",  /* 87 */
		"421211",  /* 88 */
		"212141",  /* 89 */
		"214121",  /* 90 */
		"412121",  /* 91 */
		"111143",  /* 92 */
		"111341",  /* 93 */
		"131141",  /* 94 */
		"114113",  /* 95 */
		"114311",  /* 96 */
		"411113",  /* 97 */
		"411311",  /* 98 */
		"113141",  /* 99 */
		"114131",  /* 100 */
		"311141",  /* 101 */
		"411131",  /* 102 */
		"211412",  /* 103 START A */
		"211214",  /* 104 START B */
		"211232",  /* 105 START C */
		"2331112", /* STOP */
	}
	fncA := map[rune]int{
		'1': 102,
		'2': 97,
		'3': 96,
		'4': 101,
	}
	fncB := map[rune]int{
		'{': 91,
		'1': 102,
		'2': 97,
		'3': 96,
		'4': 100,
	}

	if code == "" {
		return "", fmt.Errorf("empty code")
	}
	codeRunes := []rune(code)
	codeData := []int{}

	switch variant {
	case 'A':
		codeData = append(codeData, 103)
		escape := false
		for _, r := range codeRunes {
			if escape {
				if fncA[r] > 0 {
					codeData = append(codeData, fncA[r])
				} else {
					return "", fmt.Errorf("invalid escape sequence {%c", r)
				}
				escape = false
			} else {
				if r == '{' {
					escape = true
				} else if r >= 32 && r <= 95 {
					codeData = append(codeData, int(r-32))
				} else if r >= 0 && r <= 31 {
					codeData = append(codeData, int(r+64))
				} else {
					return "", fmt.Errorf("invalid char %c", r)
				}
			}
		}
	case 'B':
		codeData = append(codeData, 104)
		escape := false
		for _, r := range codeRunes {
			if escape {
				if fncB[r] > 0 {
					codeData = append(codeData, fncB[r])
				} else {
					return "", fmt.Errorf("invalid escape sequence {%c", r)
				}
				escape = false
			} else {
				if r == '{' {
					escape = true
				} else if r >= 32 && r <= 127 {
					codeData = append(codeData, int(r-32))
				} else {
					return "", fmt.Errorf("invalid char %c", r)
				}
			}
		}
	case 'C':
		codeData = append(codeData, 105)
		for _, r := range codeRunes {
			if r < '0' || r > '9' {
				return "", fmt.Errorf("invalid char %c", r)
			}
		}
		if len(codeRunes)%2 != 0 {
			return "", fmt.Errorf("invalid odd length %d", len(codeRunes))
		}
		for i := 0; i < len(codeRunes); i += 2 {
			value := int((codeRunes[i]-48)*10 + codeRunes[i+1] - 48)
			codeData = append(codeData, value)
		}
	default: // auto
		i := 0
		for i < len(codeRunes) {
			// Averiguamos en número de digitos consecutivos a partir de la posición i-sima
			n := 0
			for i+n < len(codeRunes) && codeRunes[i+n] >= '0' && codeRunes[i+n] <= '9' {
				n++
			}
			// Condición para usar C: 6+ digitos en cualquier posición o 4+ digitos al principio o final
			if (n >= 6) || (n >= 4 && i == 0) || (n >= 4 && i+n == len(codeRunes)) {
				// Cambiamos a C si es preciso
				switch variant {
				case 'A', 'B':
					codeData = append(codeData, 99)
					variant = 'C'
				case 'C': // nada
				default:
					codeData = append(codeData, 105)
					variant = 'C'
				}
				// Agregamos los runes en parejas
				for n > 1 {
					value := int((codeRunes[i]-48)*10 + codeRunes[i+1] - 48)
					codeData = append(codeData, value)
					n -= 2
					i += 2
				}
			} else {
				// Cambiamos a B si es preciso
				switch variant {
				case 'A', 'C':
					codeData = append(codeData, 100)
					variant = 'B'
				case 'B': // nada
				default:
					codeData = append(codeData, 104)
					variant = 'B'
				}
				// Agregamos los dígitos que no se hayan pasado como C
				for n > 0 {
					value := codeRunes[i]
					codeData = append(codeData, int(value-32))
					i++
					n--
				}
				// Agregamos hasta encontrar un dígito
				escape := false
				for i < len(codeRunes) && (escape || codeRunes[i] < '0' || codeRunes[i] > '9') {
					r := codeRunes[i]
					if escape {
						if fncB[r] > 0 {
							codeData = append(codeData, fncB[r])
						} else {
							return "", fmt.Errorf("invalid escape sequence {%c", r)
						}
						escape = false
					} else {
						if r == '{' {
							escape = true
						} else if r >= 0 && r <= 31 {
							codeData = append(codeData, 98, int(r+64))
						} else if r >= 32 && r <= 127 {
							codeData = append(codeData, int(r-32))
						} else {
							return "", fmt.Errorf("invalid char %c", r)
						}
					}
					i++
				}
			}
		}
	}
	// add check character
	var sum int
	for key, value := range codeData {
		if key == 0 {
			sum = value
		} else {
			sum += value * key
		}
	}
	sum %= 103
	codeData = append(codeData, sum)
	// add stop
	codeData = append(codeData, 106)

	barpattern := ""
	for _, c := range codeData {
		barpattern += patterns[c]
	}
	return barpattern, nil
}
