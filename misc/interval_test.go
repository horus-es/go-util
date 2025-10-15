package misc_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/horus-es/go-util/v3/errores"
	"github.com/horus-es/go-util/v3/formato"
	"github.com/horus-es/go-util/v3/misc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAddInterval(t *testing.T) {
	cinco11 := pgtype.Timestamp{Time: time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC), Valid: true}
	seis11 := pgtype.Timestamp{Time: time.Date(2021, 11, 6, 11, 12, 13, 0, time.UTC), Valid: true}
	cuatro11 := pgtype.Timestamp{Time: time.Date(2021, 11, 4, 11, 12, 13, 0, time.UTC), Valid: true}
	dia := pgtype.Interval{Days: 1, Valid: true}
	assert.Equal(t, seis11, misc.AddInterval(cinco11, dia))
	assert.Equal(t, cuatro11, misc.SubInterval(cinco11, dia))
	cinco11a := pgtype.Timestamp{Time: time.Date(2021, 11, 5, 11, 13, 13, 0, time.UTC), Valid: true}
	cinco11s := pgtype.Timestamp{Time: time.Date(2021, 11, 5, 11, 11, 13, 0, time.UTC), Valid: true}
	minuto := pgtype.Interval{Microseconds: 60000000, Valid: true}
	assert.Equal(t, cinco11a, misc.AddInterval(cinco11, minuto))
	assert.Equal(t, cinco11s, misc.SubInterval(cinco11, minuto))
	assert.False(t, misc.AddInterval(cinco11, pgtype.Interval{}).Valid)
	assert.False(t, misc.SubInterval(cinco11, pgtype.Interval{}).Valid)
}

func ExampleAddInterval() {
	unDia := pgtype.Interval{Days: 1, Valid: true}
	fechaInicial, err := formato.ParseTimestamp("4/7/2023 12:00", formato.DMA)
	errores.PanicIfError(err)
	fechaSiguiente := misc.AddInterval(fechaInicial, unDia)
	fmt.Println(formato.PrintTimestamp(fechaSiguiente, formato.DMA))
	// Output: 05/07/2023 12:00
}

func ExampleSubInterval() {
	unDia := pgtype.Interval{Days: 1, Valid: true}
	fechaInicial, err := formato.ParseTimestamp("4/7/2023 12:00", formato.DMA)
	errores.PanicIfError(err)
	fechaAnterior := misc.SubInterval(fechaInicial, unDia)
	fmt.Println(formato.PrintTimestamp(fechaAnterior, formato.DMA))
	// Output: 03/07/2023 12:00
}
