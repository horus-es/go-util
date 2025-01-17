package barcode

import (
	"fmt"
)

/* barcodeI25 Interleaved 2 of 5 barcodes.
 * Compact numeric code, widely used in industry, air cargo
 * Contains digits (0 to 9) and encodes the data in the width of both bars and spaces.
 */
func barcodeI25(code string) (bars, hri string, err error) {
	patterns := []string{
		"11221", "21112", "12112", "22111", "11212",
		"21211", "12211", "11122", "21121", "12121",
	}
	start := "1111"
	end := "211"

	if code == "" {
		return "", "", fmt.Errorf("empty code")
	}
	runes := []rune(code)
	for _, r := range runes {
		if r < '0' || r > '9' {
			return "", "", fmt.Errorf("invalid char 0x%02X", r)
		}
	}
	if len(runes)%2 != 0 {
		return "", "", fmt.Errorf("invalid odd length %d", len(runes))
	}
	bars = start
	for i := 0; i < len(runes); i = i + 2 {
		pBar := patterns[runes[i]-48]
		pSpace := patterns[runes[i+1]-48]
		for j := 0; j < 5; j++ {
			bars += string(pBar[j])
			bars += string(pSpace[j])
		}
	}
	bars += end
	hri = code
	return
}
