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
	// <svg width="156.00000" height="116.00000" version="1.1" xmlns="http://www.w3.org/2000/svg">
	// 	<g fill="#000" stroke="none">
	// 		<rect x="10.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 		<rect x="16.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="22.00000" y="0.00000" width="6.00000" height="100.00000" />
	// 		<rect x="32.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="36.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 		<rect x="44.00000" y="0.00000" width="6.00000" height="100.00000" />
	// 		<rect x="54.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="62.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="66.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 		<rect x="76.00000" y="0.00000" width="6.00000" height="100.00000" />
	// 		<rect x="88.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="92.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 		<rect x="98.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="106.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 		<rect x="112.00000" y="0.00000" width="6.00000" height="100.00000" />
	// 		<rect x="120.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 		<rect x="130.00000" y="0.00000" width="6.00000" height="100.00000" />
	// 		<rect x="138.00000" y="0.00000" width="2.00000" height="100.00000" />
	// 		<rect x="142.00000" y="0.00000" width="4.00000" height="100.00000" />
	// 	<text x="78.00000" text-anchor="middle" dominant-baseline="hanging" y="102.00000" fill="#000" font-size="15px">123456</text>
	// 	</g>
	// </svg>
}
