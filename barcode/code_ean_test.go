package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v2/barcode"
	"github.com/stretchr/testify/assert"
)

func TestEAN(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("12345670", barcode.EAN8)
	assert.NoError(t, err)
	assert.Equal(t, "1112221212214111132111111231111413123211111", bars)
	bars, err = barcode.GetBarcodeBARS("1234568", barcode.EAN8)
	assert.NoError(t, err)
	assert.Equal(t, "1112221212214111132111111231111412131312111", bars)
	bars, err = barcode.GetBarcodeBARS("1234567890128", barcode.EAN13)
	assert.NoError(t, err)
	assert.Equal(t, "11121221411231112314111213111111121331123211222121221213111", bars)
	bars, err = barcode.GetBarcodeBARS("123456789012", barcode.EAN13)
	assert.NoError(t, err)
	assert.Equal(t, "11121221411231112314111213111111121331123211222121221213111", bars)
	_, err = barcode.GetBarcodeBARS("", barcode.EAN13)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("", barcode.EAN8)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("12345671", barcode.EAN8)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("A1234567", barcode.EAN8)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("1234567890125", barcode.EAN13)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("123456789012A", barcode.EAN13)
	assert.Error(t, err)
}

func TestEAN_SVG(t *testing.T) {
	bars, err := barcode.GetBarcodeSVG("12345670", barcode.EAN8, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("EAN8.svg", []byte(bars), 0666)
	assert.NoError(t, err)
	bars, err = barcode.GetBarcodeSVG("1234567890128", barcode.EAN13, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("EAN13.svg", []byte(bars), 0666)
	assert.NoError(t, err)
}
