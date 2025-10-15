package misc_test

import (
	"fmt"

	"github.com/horus-es/go-util/v3/misc"
)

func ExampleTitle() {
	fmt.Println(misc.Title("hola mundo otra vEZ"))
	// output: Hola Mundo Otra VEZ
}
