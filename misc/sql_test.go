package misc_test

import (
	"fmt"
	"testing"

	"github.com/horus-es/go-util/v3/misc"
	"github.com/stretchr/testify/assert"
)

func TestSqlIn(t *testing.T) {
	obtiene := misc.SqlIn("uno", "dos", "tres")
	assert.Equal(t, " in ('uno','dos','tres')", obtiene)
	obtiene = misc.SqlIn("cuatro")
	assert.Equal(t, " = 'cuatro'", obtiene)
	obtiene = misc.SqlIn()
	assert.Equal(t, " is null", obtiene)
}

func ExampleSqlIn() {
	codigos := []string{"cero", "uno", "dos", "tres"}
	sql := `select * from tabla where codigo` + misc.SqlIn(codigos[1:]...)
	fmt.Println(sql)
	// Output: select * from tabla where codigo in ('uno','dos','tres')

}

func ExampleEscapeSQL() {
	sql := `select * from tabla where nombre=` + misc.EscapeSQL("O'Brian")
	fmt.Println(sql)
	// Output: select * from tabla where nombre='O''Brian'
}
