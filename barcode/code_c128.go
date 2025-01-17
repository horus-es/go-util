package barcode

import (
	"fmt"
)

// Genera un CODE128. Si automático, cambia de entre las variantes B y C de forma autómatica. si manual hay que emplear las siguientes secuencias de escape (robadas del estándar esc/pos):
//   - {A: Cambiar a variante A
//   - {B: Cambiar a variante B
//   - {C: Cambiar a variante C (estricto esc/pos)
//   - {N: Cambiar a variante C (usando pares de caracteres '0' a '9')
//   - {S: Cambiar el siguiente caracter de variante A a B o viceversa
//   - {1: Caracter de función 1
//   - {2: Caracter de función 2 (solo variantes A y B)
//   - {3: Caracter de función 3 (solo variantes A y B)
//   - {4: Caracter de función 4 (solo variantes A y B)
//   - {{: Caracter '{' en la variante B
func barcodeC128(code string, automatico bool) (bars, hri string, err error) {

	patterns := []string{
		"212222", "222122", "222221", "121223", "121322", "131222", "122213", "122312", "132212", "221213", // 0-9
		"221312", "231212", "112232", "122132", "122231", "113222", "123122", "123221", "223211", "221132", // 10-19
		"221231", "213212", "223112", "312131", "311222", "321122", "321221", "312212", "322112", "322211", // 20-29
		"212123", "212321", "232121", "111323", "131123", "131321", "112313", "132113", "132311", "211313", // 30-39
		"231113", "231311", "112133", "112331", "132131", "113123", "113321", "133121", "313121", "211331", // 40-49
		"231131", "213113", "213311", "213131", "311123", "311321", "331121", "312113", "312311", "332111", // 50-59
		"314111", "221411", "431111", "111224", "111422", "121124", "121421", "141122", "141221", "112214", // 60-69
		"112412", "122114", "122411", "142112", "142211", "241211", "221114", "413111", "241112", "134111", // 70-79
		"111242", "121142", "121241", "114212", "124112", "124211", "411212", "421112", "421211", "212141", // 80-89
		"214121", "412121", "111143", "111341", "131141", "114113", "114311", "411113", "411311", "113141", // 90-99
		"114131", "311141", "411131", "211412", "211214", "211232", "2331112", // 100-106
	}

	const (
		STARTA = 103
		STARTB = 104
		STARTC = 105
		MODEA  = 101
		MODEB  = 100
		MODEC  = 99
		STOP   = 106
		FNC1   = 102
		FNC2   = 97
		FNC3   = 96
		FNC4A  = 101
		FNC4B  = 100
		SHIFT  = 98
		ESCAPE = 91
	)
	if code == "" {
		return "", "", fmt.Errorf("empty code")
	}
	codeRunes := []rune(code)
	codeData := []int{}
	var variante rune
	if automatico {
		i := 0
		for i < len(codeRunes) {
			// Averiguamos el número de dígitos consecutivos a partir de la posición i-sima
			n := 0
			for i+n < len(codeRunes) && codeRunes[i+n] >= '0' && codeRunes[i+n] <= '9' {
				n++
			}
			// Condición para usar C: 6+ digitos en cualquier posición o 4+ digitos al principio o final
			if (n >= 6) || (n >= 4 && i == 0) || (n >= 4 && i+n == len(codeRunes)) {
				// Cambiamos a C si es preciso
				switch variante {
				case 'B':
					codeData = append(codeData, MODEC)
					variante = 'C'
				case 'C': // nada
				default:
					codeData = append(codeData, STARTC)
					variante = 'C'
				}
				// Agregamos los runes en parejas
				for n > 1 {
					r := int((codeRunes[i]-48)*10 + codeRunes[i+1] - 48)
					codeData = append(codeData, r)
					hri += fmt.Sprintf("%02d", r)
					n -= 2
					i += 2
				}
			} else {
				// Cambiamos a B si es preciso
				switch variante {
				case 'C':
					codeData = append(codeData, MODEB)
					variante = 'B'
				case 'B': // nada
				default:
					codeData = append(codeData, STARTB)
					variante = 'B'
				}
				// Agregamos los dígitos que no se hayan pasado como C
				for n > 0 {
					r := codeRunes[i]
					codeData = append(codeData, int(r-32))
					hri += string(r)
					i++
					n--
				}
				// Agregamos hasta encontrar un dígito
				escape := false
				for i < len(codeRunes) && (escape || codeRunes[i] < '0' || codeRunes[i] > '9') {
					r := codeRunes[i]
					if escape {
						switch r {
						case '1':
							codeData = append(codeData, FNC1)
						case '2':
							codeData = append(codeData, FNC2)
						case '3':
							codeData = append(codeData, FNC3)
						case '4':
							codeData = append(codeData, FNC4B)
						case '{':
							codeData = append(codeData, ESCAPE)
							hri += string(r)
						default:
							return "", "", fmt.Errorf("invalid escape sequence {%c", r)
						}
						escape = false
					} else {
						if r == '{' {
							escape = true
						} else if r >= 0 && r <= 31 {
							codeData = append(codeData, SHIFT, int(r+64))
							hri += " "
						} else if r >= 32 && r <= 127 {
							codeData = append(codeData, int(r-32))
							hri += string(r)
						} else {
							return "", "", fmt.Errorf("invalid char 0x%02X", r)
						}
					}
					i++
				}
			}
		}
	} else {
		escape := false
		shift := false
		for i := 0; i < len(codeRunes); i++ {
			r := codeRunes[i]
			if escape {
				switch r {
				case 'A':
					if variante == 0 {
						codeData = append(codeData, STARTA)
					} else {
						codeData = append(codeData, MODEA)
					}
					variante = r
				case 'B':
					if variante == 0 {
						codeData = append(codeData, STARTB)
					} else {
						codeData = append(codeData, MODEB)
					}
					variante = r
				case 'C', 'N':
					if variante == 0 {
						codeData = append(codeData, STARTC)
					} else {
						codeData = append(codeData, MODEC)
					}
					variante = r
				case 'S':
					if variante == 'A' || variante == 'B' {
						shift = true
						codeData = append(codeData, SHIFT)
					} else {
						return "", "", fmt.Errorf("invalid escape sequence {%c", r)
					}
				case '1':
					codeData = append(codeData, FNC1)
				case '2':
					if variante == 'A' || variante == 'B' {
						codeData = append(codeData, FNC2)
					} else {
						return "", "", fmt.Errorf("invalid escape sequence {%c", r)
					}
				case '3':
					if variante == 'A' || variante == 'B' {
						codeData = append(codeData, FNC3)
					} else {
						return "", "", fmt.Errorf("invalid escape sequence {%c", r)
					}
				case '4':
					switch variante {
					case 'A':
						codeData = append(codeData, FNC4A)
					case 'B':
						codeData = append(codeData, FNC4B)
					default:
						return "", "", fmt.Errorf("invalid escape sequence {%c", r)
					}
				case '{':
					if variante == 'B' {
						codeData = append(codeData, ESCAPE)
						hri += string(r)
					} else {
						return "", "", fmt.Errorf("invalid escape sequence {%c", r)
					}
				default:
					return "", "", fmt.Errorf("invalid escape sequence {%c", r)
				}
				escape = false
			} else if r == '{' {
				escape = true
			} else {
				if shift {
					switch variante {
					case 'A':
						variante = 'B'
					case 'B':
						variante = 'A'
					}
				}
				switch variante {
				case 'A':
					if r >= 32 && r <= 95 {
						codeData = append(codeData, int(r-32))
						hri += string(r)
					} else if r >= 0 && r <= 31 {
						codeData = append(codeData, int(r+64))
						hri += " "
					} else {
						return "", "", fmt.Errorf("invalid char 0x%02X", r)
					}
				case 'B':
					if r >= 32 && r <= 127 {
						codeData = append(codeData, int(r-32))
						hri += string(r)
					} else {
						return "", "", fmt.Errorf("invalid char 0x%02X", r)
					}
				case 'C':
					if r >= 0 && r <= 99 {
						codeData = append(codeData, int(r))
						hri += fmt.Sprintf("%02d", r)
					} else {
						return "", "", fmt.Errorf("invalid number %d", r)
					}
				case 'N':
					i++
					var r2 rune
					if i < len(codeRunes) {
						r2 = codeRunes[i]
					}
					if r >= '0' && r <= '9' && r2 >= '0' && r2 <= '9' {
						r = (r-48)*10 + r2 - 48
						codeData = append(codeData, int(r))
						hri += fmt.Sprintf("%02d", r)
					} else {
						return "", "", fmt.Errorf("invalid numeric sequence")
					}
				default:
					return "", "", fmt.Errorf("missing start sequence")
				}
				if shift {
					switch variante {
					case 'A':
						variante = 'B'
					case 'B':
						variante = 'A'
					}
					shift = false
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
	codeData = append(codeData, STOP)
	// Pasamos a bars
	for _, c := range codeData {
		bars += patterns[c]
	}
	return
}
