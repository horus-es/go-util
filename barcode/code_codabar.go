package barcode

import (
	"fmt"
)

func barcodeCODABAR(code string) (string, error) {
	patterns := map[rune]string{
		'0': "1111133",
		'1': "1111331",
		'2': "1113113",
		'3': "3311111",
		'4': "1131131",
		'5': "3111131",
		'6': "1311113",
		'7': "1311311",
		'8': "1331111",
		'9': "3113111",
		'-': "1113311",
		'$': "1133111",
		':': "3111313",
		'/': "3131113",
		'.': "3131311",
		'+': "1131313",
		'a': "1133131",
		'b': "1313113",
		'c': "1113133",
		'd': "1113331",
	}

	// Ajustes de start/stop
	runes := []rune(code)
	if len(runes) == 0 {
		return "", fmt.Errorf("empty code")
	}
	if runes[0] >= 'A' && runes[0] <= 'D' {
		runes[0] += 32
	}
	if runes[len(runes)-1] >= 'A' && runes[len(runes)-1] <= 'D' {
		runes[len(runes)-1] += 32
	}
	if runes[0] < 'a' || runes[0] > 'd' {
		runes = append([]rune{'a'}, runes...)
	}
	if runes[len(runes)-1] < 'a' || runes[len(runes)-1] > 'd' {
		runes = append(runes, 'd')
	}
	barpattern := ""
	for k, r := range runes {
		p := patterns[r]
		if p == "" {
			return "", fmt.Errorf("invalid char %c", r)
		}
		if r >= 'a' && r <= 'd' && k > 0 && k < len(runes)-1 {
			return "", fmt.Errorf("invalid char %c", r)
		}
		if k > 0 {
			barpattern += "1"
		}
		barpattern += p
	}
	return barpattern, nil
}
