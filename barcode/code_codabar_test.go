package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v3/barcode"
	"github.com/stretchr/testify/assert"
)

func TestCODABAR(t *testing.T) {
	bars, hri, err := barcode.GetBarcodeBARS("123456", barcode.CODABAR)
	assert.NoError(t, err)
	assert.Equal(t, "113313111111331111131131331111111131131131111311131111311113331", bars)
	assert.Equal(t, "123456", hri)
	bars, err = barcode.GetBarcodeSVG("a123456c", barcode.CODABAR, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("CODABAR.svg", []byte(bars), 0666)
	assert.NoError(t, err)
	_, _, err = barcode.GetBarcodeBARS("A1d23456A", barcode.CODABAR)
	assert.Error(t, err)
	_, _, err = barcode.GetBarcodeBARS("t123456C", barcode.CODABAR)
	assert.Error(t, err)
	_, _, err = barcode.GetBarcodeBARS("", barcode.CODABAR)
	assert.Error(t, err)
}
