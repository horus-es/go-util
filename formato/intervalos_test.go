package formato

import (
	"fmt"
	"testing"
	"time"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestDuracion(t *testing.T) {
	d, err := ParseDuracion("48h")
	assert.Equal(t, "48h", PrintDuracion(d))
	assert.Equal(t, "PT172800S", PrintDuracionIso(d))
	assert.Nil(t, err)
	d, err = ParseDuracion("30m")
	assert.Equal(t, "30m", PrintDuracion(d))
	assert.Equal(t, "PT1800S", PrintDuracionIso(d))
	assert.Nil(t, err)
	d, err = ParseDuracion("120s")
	assert.Equal(t, "2m", PrintDuracion(d))
	assert.Equal(t, "PT120S", PrintDuracionIso(d))
	assert.Nil(t, err)
	d, err = ParseDuracion("  120M   ")
	assert.Equal(t, "2h", PrintDuracion(d))
	assert.Nil(t, err)
	d, err = ParseDuracion("121s")
	assert.Equal(t, "121s", PrintDuracion(d))
	assert.Nil(t, err)
	d, err = ParseDuracion("0s")
	assert.Equal(t, "", PrintDuracion(d))
	assert.Nil(t, err)
	_, err = ParseDuracion("59M01S")
	assert.NotNil(t, err)
	_, err = ParseDuracion("59m01z")
	assert.NotNil(t, err)
	_, err = ParseDuracion("00")
	assert.NotNil(t, err)
	_, err = ParseDuracion("  ")
	assert.NotNil(t, err)
}

func ExamplePrintDuracion() {
	d, err := ParseDuracion("48h")
	errores.PanicIfError(err)
	fmt.Println(PrintDuracion(d))
	// Output: 48h
}

func ExamplePrintDuracionIso() {
	d, err := ParseDuracion("48h")
	errores.PanicIfError(err)
	fmt.Println(PrintDuracionIso(d))
	// Output: PT172800S
}

func ExampleParseDuracion() {
	d, err := ParseDuracion("48h")
	errores.PanicIfError(err)
	fmt.Println(PrintDuracion(d))
	// Output: 48h
}

func TestInterval(t *testing.T) {
	d, err := ParseIntervalAMSD("48d")
	assert.Equal(t, "48d", PrintInterval(d))
	assert.Equal(t, "P48D", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalAMSD("3 meses")
	assert.Equal(t, "3m", PrintInterval(d))
	assert.Equal(t, "P3M", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalAMSD("2 semanas")
	assert.Equal(t, "14d", PrintInterval(d))
	assert.Equal(t, "P14D", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalAMSD("1año")
	assert.Equal(t, "12m", PrintInterval(d))
	assert.Equal(t, "P12M", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalAMSD("  ")
	assert.Nil(t, err)
	assert.Equal(t, pgtype.Null, d.Status)
	d, err = ParseIntervalHMS("33 horas")
	assert.Equal(t, "33h", PrintInterval(d))
	assert.Equal(t, "PT118800S", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalHMS("61 minutos")
	assert.Equal(t, "61m", PrintInterval(d))
	assert.Equal(t, "PT3660S", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalHMS("61 segundos")
	assert.Equal(t, "61s", PrintInterval(d))
	assert.Equal(t, "PT61S", PrintIntervalIso(d))
	assert.Nil(t, err)
	d, err = ParseIntervalHMS("  ")
	assert.Nil(t, err)
	assert.Equal(t, pgtype.Null, d.Status)
	_, err = ParseIntervalAMSD("1m2m")
	assert.NotNil(t, err)
	_, err = ParseIntervalHMS("1d")
	assert.NotNil(t, err)
	_, err = ParseIntervalAMSD("1h")
	assert.NotNil(t, err)
	assert.Equal(t, "", PrintInterval(pgtype.Interval{}))
	assert.Equal(t, "", PrintIntervalIso(pgtype.Interval{}))
	assert.Equal(t, "", PrintInterval(pgtype.Interval{Status: pgtype.Present}))
	assert.Equal(t, "", PrintIntervalIso(pgtype.Interval{Status: pgtype.Present}))
	assert.Equal(t, "1m", PrintInterval(pgtype.Interval{Microseconds: 60000000, Status: pgtype.Present}))
	assert.Equal(t, "PT60S", PrintIntervalIso(pgtype.Interval{Microseconds: 60000000, Status: pgtype.Present}))
}

func ExampleParseIntervalAMSD() {
	i, err := ParseIntervalAMSD("30d")
	errores.PanicIfError(err)
	fmt.Println(PrintInterval(i))
	// Output: 30d
}

func ExampleParseIntervalHMS() {
	i, err := ParseIntervalHMS("48h")
	errores.PanicIfError(err)
	fmt.Println(PrintInterval(i))
	// Output: 48h
}

func ExamplePrintInterval() {
	i, err := ParseIntervalAMSD("30d")
	errores.PanicIfError(err)
	fmt.Println(PrintInterval(i))
	// Output: 30d
}

func ExamplePrintIntervalIso() {
	i, err := ParseIntervalAMSD("30d")
	errores.PanicIfError(err)
	fmt.Println(PrintIntervalIso(i))
	// Output: P30D
}

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
	fechaInicial, err := ParseTimestamp("4/7/2023 12:00", DMA)
	errores.PanicIfError(err)
	fechaSiguiente := AddInterval(fechaInicial, unDia)
	fmt.Println(PrintTimestamp(fechaSiguiente, DMA))
	// Output: 05/07/2023 12:00
}

func ExampleSubInterval() {
	unDia := pgtype.Interval{Days: 1, Status: pgtype.Present}
	fechaInicial, err := ParseTimestamp("4/7/2023 12:00", DMA)
	errores.PanicIfError(err)
	fechaAnterior := SubInterval(fechaInicial, unDia)
	fmt.Println(PrintTimestamp(fechaAnterior, DMA))
	// Output: 03/07/2023 12:00
}
