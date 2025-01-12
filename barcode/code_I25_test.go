package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v2/barcode"
	"github.com/stretchr/testify/assert"
)

func TestI25(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("123456", barcode.I25)
	assert.NoError(t, err)
	assert.Equal(t, "1111211211112221211211122112221111211", bars)
	_, err = barcode.GetBarcodeBARS("12345", barcode.I25)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("12345A", barcode.I25)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("", barcode.I25)
	assert.Error(t, err)
}

func TestI25SVG(t *testing.T) {
	bars, err := barcode.GetBarcodeSVG("123456", barcode.I25, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("I25.svg", []byte(bars), 0666)
	assert.NoError(t, err)
}
