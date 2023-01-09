package misc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSqlIn(t *testing.T) {
	obtiene := SqlIn("uno", "dos", "tres")
	assert.Equal(t, " in ('uno','dos','tres')", obtiene)
	obtiene = SqlIn("cuatro")
	assert.Equal(t, " = 'cuatro'", obtiene)
	obtiene = SqlIn()
	assert.Equal(t, " is null", obtiene)
}

func ExampleSqlIn() {
	codigos := []string{"cero", "uno", "dos", "tres"}
	sql := `select * from tabla where codigo` + SqlIn(codigos[1:]...)
	fmt.Println(sql)
	// Output: select * from tabla where codigo in ('uno','dos','tres')

}

func ExampleEscapeSQL() {
	sql := `select * from tabla where nombre=` + EscapeSQL("O'Brian")
	fmt.Println(sql)
	// Output: select * from tabla where nombre='O''Brian'
}
