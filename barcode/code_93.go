package barcode

import (
	"fmt"
	"strings"
)

func barcodeCode93(code string) (string, error) {

	patterns := map[rune]string{
		'0': "131112", '1': "111213", '2': "111312", '3': "111411", '4': "121113",
		'5': "121212", '6': "121311", '7': "111114", '8': "131211", '9': "141111",
		'A': "211113", 'B': "211212", 'C': "211311", 'D': "221112", 'E': "221211",
		'F': "231111", 'G': "112113", 'H': "112212", 'I': "112311", 'J': "122112",
		'K': "132111", 'L': "111123", 'M': "111222", 'N': "111321", 'O': "121122",
		'P': "131121", 'Q': "212112", 'R': "212211", 'S': "211122", 'T': "211221",
		'U': "221121", 'V': "222111", 'W': "112122", 'X': "112221", 'Y': "122121",
		'Z': "123111", '-': "121131", '.': "311112", ' ': "311211", '$': "321111",
		'/': "112131", '+': "113121", '%': "211131",
		128: "121221", // ($)
		129: "311121", // (/)
		130: "122211", // (+)
		131: "312111", // (%)
		42:  "111141", // start-stop
	}

	code = strings.TrimPrefix(code, "*")
	code = strings.TrimSuffix(code, "*")
	if code == "" {
		return "", fmt.Errorf("empty code")
	}
	code += checksumCode93([]rune(code))
	code = "*" + code + "*"
	barpattern := ""
	for _, r := range code {
		p := patterns[r]
		if len(p) == 0 {
			return "", fmt.Errorf("invalid char %c", r)
		}
		barpattern += p
	}
	barpattern += "1"
	return barpattern, nil
}

// checksumCode93 calculates checksum of
// code digit C & K
func checksumCode93(code []rune) string {
	chars := map[rune]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19, 'K': 20,
		'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29, 'U': 30, 'V': 31,
		'W': 32, 'X': 33, 'Y': 34, 'Z': 35, '-': 36, '.': 37, ' ': 38, '$': 39, '/': 40, '+': 41, '%': 42,
		128: 43, 129: 44, 130: 45, 131: 46,
	}

	// calc check digit C
	weight := 1
	check := 0
	for i := len(code) - 1; i >= 0; i-- {
		char := code[i]
		k := chars[char]
		check = check + (k * weight)
		if weight++; weight > 20 {
			weight = 1
		}
	}
	check = check % 47
	var c rune
	for key, value := range chars {
		if value == check {
			c = key
		}
	}

	// calc check digit K
	code = append(code, c)
	weight = 1
	check = 0
	for i := len(code) - 1; i >= 0; i-- {
		char := []rune(code)[i]
		k := chars[char]
		check = check + (k * weight)
		if weight++; weight > 15 {
			weight = 1
		}
	}
	check = check % 47
	var k rune
	for key, value := range chars {
		if value == check {
			k = key
		}
	}
	return string(c) + string(k)
}
