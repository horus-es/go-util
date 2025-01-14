package barcode_test

import (
	"fmt"

	"github.com/horus-es/go-util/v2/barcode"
)

func ExampleGetBarcodeBARS() {
	bars, _ := barcode.GetBarcodeBARS("123456", barcode.C128)
	fmt.Print(bars)
	// Output: 2112321122321311233311211321312331112
}

func ExampleGetBarcodeSVG() {
	svg, _ := barcode.GetBarcodeSVG("123456", barcode.C128, 2, 100, "#000", barcode.Below, false)
	fmt.Print(svg)

	// Output:
	// <?xml version="1.0" standalone="no" ?>
	// <!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
	// <svg width="156" height="114" version="1.1" xmlns="http://www.w3.org/2000/svg">
	// 	<g fill="#000" stroke="none">
	// 		<rect x="10" y="0" width="4" height="100" />
	// 		<rect x="16" y="0" width="2" height="100" />
	// 		<rect x="22" y="0" width="6" height="100" />
	// 		<rect x="32" y="0" width="2" height="100" />
	// 		<rect x="36" y="0" width="4" height="100" />
	// 		<rect x="44" y="0" width="6" height="100" />
	// 		<rect x="54" y="0" width="2" height="100" />
	// 		<rect x="62" y="0" width="2" height="100" />
	// 		<rect x="66" y="0" width="4" height="100" />
	// 		<rect x="76" y="0" width="6" height="100" />
	// 		<rect x="88" y="0" width="2" height="100" />
	// 		<rect x="92" y="0" width="4" height="100" />
	// 		<rect x="98" y="0" width="2" height="100" />
	// 		<rect x="106" y="0" width="4" height="100" />
	// 		<rect x="112" y="0" width="6" height="100" />
	// 		<rect x="120" y="0" width="4" height="100" />
	// 		<rect x="130" y="0" width="6" height="100" />
	// 		<rect x="138" y="0" width="2" height="100" />
	// 		<rect x="142" y="0" width="4" height="100" />
	// 	<text x="78" text-anchor="middle" dominant-baseline="hanging" y="100" fill="#000" font-size="12px">123456</text>
	// 	</g>
	// </svg>
}
