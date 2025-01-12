package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v2/barcode"
	"github.com/stretchr/testify/assert"
)

func TestUPC(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("04210000526", barcode.UPCA)
	assert.NoError(t, err)
	assert.Equal(t, "11132111132212222213211321111111321132111231212211141132111", bars)
	bars, err = barcode.GetBarcodeBARS("042100005264", barcode.UPCA)
	assert.NoError(t, err)
	assert.Equal(t, "11132111132212222213211321111111321132111231212211141132111", bars)
	bars, err = barcode.GetBarcodeBARS("425261", barcode.UPCE)
	assert.NoError(t, err)
	assert.Equal(t, "111231121221321221211142221111111", bars)
	bars, err = barcode.GetBarcodeBARS("0425261", barcode.UPCE)
	assert.NoError(t, err)
	assert.Equal(t, "111231121221321221211142221111111", bars)
	bars, err = barcode.GetBarcodeBARS("04252614", barcode.UPCE)
	assert.NoError(t, err)
	assert.Equal(t, "111231121221321221211142221111111", bars)
	bars, err = barcode.GetBarcodeBARS("04210000526", barcode.UPCE)
	assert.NoError(t, err)
	assert.Equal(t, "111231121221321221211142221111111", bars)
	bars, err = barcode.GetBarcodeBARS("042100005264", barcode.UPCE)
	assert.NoError(t, err)
	assert.Equal(t, "111231121221321221211142221111111", bars)
	_, err = barcode.GetBarcodeBARS("0421000264", barcode.UPCA)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("0421000264", barcode.UPCE)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("042100X05264", barcode.UPCA)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("042100X05264", barcode.UPCE)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("042100005260", barcode.UPCA)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("042100005260", barcode.UPCE)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("04252611", barcode.UPCE)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("24210000526", barcode.UPCE)
	assert.Error(t, err)
}

func TestUPC_SVG(t *testing.T) {
	bars, err := barcode.GetBarcodeSVG("042100005264", barcode.UPCA, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("UPCA.svg", []byte(bars), 0666)
	assert.NoError(t, err)
	bars, err = barcode.GetBarcodeSVG("425261", barcode.UPCE, 2, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("UPCE.svg", []byte(bars), 0666)
	assert.NoError(t, err)
}
