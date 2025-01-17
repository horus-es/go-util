package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v2/barcode"
	"github.com/stretchr/testify/assert"
)

func TestCode39(t *testing.T) {
	bars, hri, err := barcode.GetBarcodeBARS("1234ABCD", barcode.C39)
	assert.NoError(t, err)
	assert.Equal(t, "131131311131131111311133111131313311111111133111313111131131113113113131311311111111331131131131311", bars)
	assert.Equal(t, "*1234ABCD*", hri)
	_, _, err = barcode.GetBarcodeBARS("1234:ABCD", barcode.C39)
	assert.Error(t, err)
	_, _, err = barcode.GetBarcodeBARS("", barcode.C39)
	assert.Error(t, err)
}

func TestCode39SVG(t *testing.T) {
	bars, err := barcode.GetBarcodeSVG("1234ABCD", barcode.C39, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("C39.svg", []byte(bars), 0666)
	assert.NoError(t, err)
}
