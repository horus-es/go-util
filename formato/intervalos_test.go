package formato_test

import (
	"fmt"
	"testing"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestDuracion(t *testing.T) {
	d, err := formato.ParseDuracion("48h")
	assert.Equal(t, "48h", formato.PrintDuracion(d))
	assert.Equal(t, "PT172800S", formato.PrintDuracionIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseDuracion("30m")
	assert.Equal(t, "30m", formato.PrintDuracion(d))
	assert.Equal(t, "PT1800S", formato.PrintDuracionIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseDuracion("120s")
	assert.Equal(t, "2m", formato.PrintDuracion(d))
	assert.Equal(t, "PT120S", formato.PrintDuracionIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseDuracion("  120M   ")
	assert.Equal(t, "2h", formato.PrintDuracion(d))
	assert.NoError(t, err)
	d, err = formato.ParseDuracion("121s")
	assert.Equal(t, "121s", formato.PrintDuracion(d))
	assert.NoError(t, err)
	d, err = formato.ParseDuracion("0s")
	assert.Equal(t, "", formato.PrintDuracion(d))
	assert.NoError(t, err)
	_, err = formato.ParseDuracion("59M01S")
	assert.NotNil(t, err)
	_, err = formato.ParseDuracion("59m01z")
	assert.NotNil(t, err)
	_, err = formato.ParseDuracion("00")
	assert.NotNil(t, err)
	_, err = formato.ParseDuracion("  ")
	assert.NotNil(t, err)
}

func ExamplePrintDuracion() {
	d, err := formato.ParseDuracion("48h")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintDuracion(d))
	// Output: 48h
}

func ExamplePrintDuracionIso() {
	d, err := formato.ParseDuracion("48h")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintDuracionIso(d))
	// Output: PT172800S
}

func ExampleParseDuracion() {
	d, err := formato.ParseDuracion("48h")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintDuracion(d))
	// Output: 48h
}

func TestInterval(t *testing.T) {
	d, err := formato.ParseIntervalAMSD("48d")
	assert.Equal(t, "48d", formato.PrintInterval(d))
	assert.Equal(t, "P48D", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalAMSD("3 meses")
	assert.Equal(t, "3m", formato.PrintInterval(d))
	assert.Equal(t, "P3M", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalAMSD("2 semanas")
	assert.Equal(t, "14d", formato.PrintInterval(d))
	assert.Equal(t, "P14D", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalAMSD("1año")
	assert.Equal(t, "12m", formato.PrintInterval(d))
	assert.Equal(t, "P12M", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalAMSD("  ")
	assert.NoError(t, err)
	assert.False(t, d.Valid)
	d, err = formato.ParseIntervalHMS("33 horas")
	assert.Equal(t, "33h", formato.PrintInterval(d))
	assert.Equal(t, "PT118800S", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalHMS("61 minutos")
	assert.Equal(t, "61m", formato.PrintInterval(d))
	assert.Equal(t, "PT3660S", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalHMS("61 segundos")
	assert.Equal(t, "61s", formato.PrintInterval(d))
	assert.Equal(t, "PT61S", formato.PrintIntervalIso(d))
	assert.NoError(t, err)
	d, err = formato.ParseIntervalHMS("  ")
	assert.NoError(t, err)
	assert.False(t, d.Valid)
	_, err = formato.ParseIntervalAMSD("1m2m")
	assert.NotNil(t, err)
	_, err = formato.ParseIntervalHMS("1d")
	assert.NotNil(t, err)
	_, err = formato.ParseIntervalAMSD("1h")
	assert.NotNil(t, err)
	assert.Equal(t, "", formato.PrintInterval(pgtype.Interval{}))
	assert.Equal(t, "", formato.PrintIntervalIso(pgtype.Interval{}))
	assert.Equal(t, "", formato.PrintInterval(pgtype.Interval{Valid: true}))
	assert.Equal(t, "", formato.PrintIntervalIso(pgtype.Interval{Valid: true}))
	assert.Equal(t, "1m", formato.PrintInterval(pgtype.Interval{Microseconds: 60000000, Valid: true}))
	assert.Equal(t, "PT60S", formato.PrintIntervalIso(pgtype.Interval{Microseconds: 60000000, Valid: true}))
}

func ExampleParseIntervalAMSD() {
	i, err := formato.ParseIntervalAMSD("30d")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintInterval(i))
	// Output: 30d
}

func ExampleParseIntervalHMS() {
	i, err := formato.ParseIntervalHMS("48h")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintInterval(i))
	// Output: 48h
}

func ExamplePrintInterval() {
	i, err := formato.ParseIntervalAMSD("30d")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintInterval(i))
	// Output: 30d
}

func ExamplePrintIntervalIso() {
	i, err := formato.ParseIntervalAMSD("30d")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintIntervalIso(i))
	// Output: P30D
}
