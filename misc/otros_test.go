package misc

import (
	"fmt"
)

func ExampleMax() {
	fmt.Println(Max(-2, 3, 1))
	// Output: 3
}

func ExampleMin() {
	fmt.Println(Min(3, -2, 3, 1))
	// Output: -2
}

func ExampleTitle() {
	fmt.Println(Title("hola mundo otra vEZ"))
	// output: Hola Mundo Otra VEZ
}
