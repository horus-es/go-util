package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v2/barcode"
	"github.com/stretchr/testify/assert"
)

func TestCode93(t *testing.T) {
	bars, hri, err := barcode.GetBarcodeBARS("1234ABCD", barcode.C93)
	assert.NoError(t, err)
	assert.Equal(t, "1111411112131113121114111211132111132112122113112211122112211321111111411", bars)
	assert.Equal(t, "1234ABCD", hri)
	_, _, err = barcode.GetBarcodeBARS("1234:ABCD", barcode.C93)
	assert.Error(t, err)
	_, _, err = barcode.GetBarcodeBARS("123ABC*", barcode.C93)
	assert.NoError(t, err)
	_, _, err = barcode.GetBarcodeBARS("123abc", barcode.C93)
	assert.Error(t, err)
	_, _, err = barcode.GetBarcodeBARS("**", barcode.C93)
	assert.Error(t, err)
}

func TestCode93SVG(t *testing.T) {
	bars, err := barcode.GetBarcodeSVG("1234ABCD", barcode.C93, 2, 100, "#000", barcode.Above, false)
	assert.NoError(t, err)
	err = os.WriteFile("C93.svg", []byte(bars), 0666)
	assert.NoError(t, err)
}
