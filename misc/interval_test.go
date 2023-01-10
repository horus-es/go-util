package misc

import (
	"fmt"
	"testing"
	"time"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAddInterval(t *testing.T) {
	cinco11 := pgtype.Timestamp{Time: time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC), Status: pgtype.Present}
	seis11 := pgtype.Timestamp{Time: time.Date(2021, 11, 6, 11, 12, 13, 0, time.UTC), Status: pgtype.Present}
	cuatro11 := pgtype.Timestamp{Time: time.Date(2021, 11, 4, 11, 12, 13, 0, time.UTC), Status: pgtype.Present}
	dia := pgtype.Interval{Days: 1, Status: pgtype.Present}
	assert.Equal(t, seis11, AddInterval(cinco11, dia))
	assert.Equal(t, cuatro11, SubInterval(cinco11, dia))
	cinco11a := pgtype.Timestamp{Time: time.Date(2021, 11, 5, 11, 13, 13, 0, time.UTC), Status: pgtype.Present}
	cinco11s := pgtype.Timestamp{Time: time.Date(2021, 11, 5, 11, 11, 13, 0, time.UTC), Status: pgtype.Present}
	minuto := pgtype.Interval{Microseconds: 60000000, Status: pgtype.Present}
	assert.Equal(t, cinco11a, AddInterval(cinco11, minuto))
	assert.Equal(t, cinco11s, SubInterval(cinco11, minuto))
	assert.Equal(t, pgtype.Null, AddInterval(cinco11, pgtype.Interval{}).Status)
	assert.Equal(t, pgtype.Null, SubInterval(cinco11, pgtype.Interval{}).Status)
}

func ExampleAddInterval() {
	unDia := pgtype.Interval{Days: 1, Status: pgtype.Present}
	fechaInicial, err := formato.ParseTimestamp("4/7/2023 12:00", formato.DMA)
	errores.PanicIfError(err)
	fechaSiguiente := AddInterval(fechaInicial, unDia)
	fmt.Println(formato.PrintTimestamp(fechaSiguiente, formato.DMA))
	// Output: 05/07/2023 12:00
}

func ExampleSubInterval() {
	unDia := pgtype.Interval{Days: 1, Status: pgtype.Present}
	fechaInicial, err := formato.ParseTimestamp("4/7/2023 12:00", formato.DMA)
	errores.PanicIfError(err)
	fechaAnterior := SubInterval(fechaInicial, unDia)
	fmt.Println(formato.PrintTimestamp(fechaAnterior, formato.DMA))
	// Output: 03/07/2023 12:00
}
