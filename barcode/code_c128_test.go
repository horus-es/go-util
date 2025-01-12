package barcode_test

import (
	"os"
	"testing"

	"github.com/horus-es/go-util/v2/barcode"
	"github.com/stretchr/testify/assert"
)

func TestC128A(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("123456", barcode.C128A)
	assert.NoError(t, err)
	assert.Equal(t, "2114121232212232112211322212312132122231121132222331112", bars)
	bars, err = barcode.GetBarcodeBARS("A1Z9", barcode.C128A)
	assert.NoError(t, err)
	assert.Equal(t, "2114121113231232213123113211222321212331112", bars)
	_, err = barcode.GetBarcodeBARS("a1Z9", barcode.C128A)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("A1{{", barcode.C128A)
	assert.Error(t, err)
	bars, err = barcode.GetBarcodeBARS("A1{2Z9", barcode.C128A)
	assert.NoError(t, err)
	assert.Equal(t, "2114121113231232214111133123113211224111132331112", bars)
	bars, err = barcode.GetBarcodeBARS("DOS\nLINEAS", barcode.C128A)
	assert.NoError(t, err)
	assert.Equal(t, "2114121123131331212131131422111321312313111133211321131113232131133111412331112", bars)
}

func TestC128B(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("A1Z9", barcode.C128B)
	assert.NoError(t, err)
	assert.Equal(t, "2112141113231232213123113211221113232331112", bars)
	bars, err = barcode.GetBarcodeBARS("a1z9", barcode.C128B)
	assert.NoError(t, err)
	assert.Equal(t, "2112141211241232212141213211223123112331112", bars)
	_, err = barcode.GetBarcodeBARS("a1ñz9", barcode.C128B)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("a1{z9", barcode.C128B)
	assert.Error(t, err)
	bars, err = barcode.GetBarcodeBARS("a1{{}z9", barcode.C128B)
	assert.NoError(t, err)
	assert.Equal(t, "2112141211241232214121211113412141213211221222132331112", bars)
}

func TestC128C(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("123456", barcode.C128C)
	assert.NoError(t, err)
	assert.Equal(t, "2112321122321311233311211321312331112", bars)
	_, err = barcode.GetBarcodeBARS("12345", barcode.C128C)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("12345Q", barcode.C128C)
	assert.Error(t, err)
}

func TestC128AUTO(t *testing.T) {
	bars, err := barcode.GetBarcodeBARS("123456", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112321122321311233311211321312331112", bars)
	bars, err = barcode.GetBarcodeBARS("12345", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112321122321311231141312132123111232331112", bars)
	bars, err = barcode.GetBarcodeBARS("12345Q", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112321122321311231141312132122113311113412331112", bars)
	bars, err = barcode.GetBarcodeBARS("A1Z9", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141113231232213123113211221113232331112", bars)
	bars, err = barcode.GetBarcodeBARS("1234a", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112321122321311231141311211243112222331112", bars)
	bars, err = barcode.GetBarcodeBARS("123ab", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141232212232112211321211241214211142122331112", bars)
	bars, err = barcode.GetBarcodeBARS("a1234", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141211241131411122321311233112222331112", bars)
	bars, err = barcode.GetBarcodeBARS("ab123", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141211241214211232212232112211321213222331112", bars)
	bars, err = barcode.GetBarcodeBARS("a123456z", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141211241131411122321311233311211141312141211111432331112", bars)
	bars, err = barcode.GetBarcodeBARS("a12345z", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141211241232212232112211322212312132122141213311212331112", bars)
	bars, err = barcode.GetBarcodeBARS("{{123456}", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112144121211131411122321311233311211141311113411123132331112", bars)
	bars, err = barcode.GetBarcodeBARS("dos\nlineas", barcode.C128)
	assert.NoError(t, err)
	assert.Equal(t, "2112141412211341111142124113111422112211141421122411121122141211241142121122142331112", bars)
	_, err = barcode.GetBarcodeBARS("a12345ñ", barcode.C128)
	assert.Error(t, err)
	_, err = barcode.GetBarcodeBARS("{ADEU}", barcode.C128)
	assert.Error(t, err)
}

func TestC128SVG(t *testing.T) {
	bars, err := barcode.GetBarcodeSVG("1234ABCD", barcode.C128A, 3, 100, "#000", barcode.Both, false)
	assert.NoError(t, err)
	err = os.WriteFile("C128A.svg", []byte(bars), 0666)
	assert.NoError(t, err)
	bars, err = barcode.GetBarcodeSVG("1234abcd", barcode.C128B, 3, 100, "#000", barcode.Above, false)
	assert.NoError(t, err)
	err = os.WriteFile("C128B.svg", []byte(bars), 0666)
	assert.NoError(t, err)
	bars, err = barcode.GetBarcodeSVG("12345678", barcode.C128C, 3, 100, "#000", barcode.Below, false)
	assert.NoError(t, err)
	err = os.WriteFile("C128C.svg", []byte(bars), 0666)
	assert.NoError(t, err)
	bars, err = barcode.GetBarcodeSVG("1234ABCD", barcode.C128, 3, 100, "#000", barcode.None, false)
	assert.NoError(t, err)
	err = os.WriteFile("C128.svg", []byte(bars), 0666)
	assert.NoError(t, err)
}
