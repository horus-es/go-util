package formato_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/stretchr/testify/assert"
)

func compruebaParseFechaHora(t *testing.T, s string, ff formato.Fecha, espera time.Time) {
	obtiene, err := formato.ParseFechaHora(s, ff)
	assert.NoError(t, err)
	assert.Equal(t, espera, obtiene)
}

func TestParseFechaHora(t *testing.T) {
	cinconoviembre := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	cinconoviembre0 := time.Date(2021, 11, 5, 11, 12, 0, 0, time.UTC)
	cinconoviembre00 := time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC)
	cincoagosto00 := time.Date(2021, 8, 5, 0, 0, 0, 0, time.UTC)
	compruebaParseFechaHora(t, "2021-11-05T11:12:13Z", formato.ISO, cinconoviembre)
	compruebaParseFechaHora(t, "2021-11-05T11:12:13", formato.ISO, cinconoviembre)
	compruebaParseFechaHora(t, "2021-11-05 11:12:13", formato.AMD, cinconoviembre)
	compruebaParseFechaHora(t, "2021-11-05", formato.ISO, cinconoviembre00)
	compruebaParseFechaHora(t, "05/11/2021 11:12:13", formato.DMA, cinconoviembre)
	compruebaParseFechaHora(t, "05/11/2021 11:12", formato.DMA, cinconoviembre0)
	compruebaParseFechaHora(t, "05/11/2021", formato.DMA, cinconoviembre00)
	compruebaParseFechaHora(t, "11/05/2021 11:12:13", formato.MDA, cinconoviembre)
	compruebaParseFechaHora(t, "11/05/2021 11:12", formato.MDA, cinconoviembre0)
	compruebaParseFechaHora(t, "2021-11-05 11:12", formato.AMD, cinconoviembre0)
	compruebaParseFechaHora(t, "11/05/2021", formato.MDA, cinconoviembre00)
	compruebaParseFechaHora(t, "11-05-2021", formato.MDA, cinconoviembre00)
	compruebaParseFechaHora(t, "11.05.2021", formato.MDA, cinconoviembre00)
	compruebaParseFechaHora(t, "5.8-2021", formato.DMA, cincoagosto00)
	compruebaParseFechaHora(t, "2021-8-5", formato.AMD, cincoagosto00)
	_, err := formato.ParseFechaHora("05-11-2021 15:16:17", formato.ISO)
	assert.NotNil(t, err)
}

func TestPrintFechaHora(t *testing.T) {
	cinconoviembre := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	assert.Equal(t, "2021-11-05T11:12:13", formato.PrintFechaHora(cinconoviembre, formato.ISO))
	assert.Equal(t, "05/11/2021 11:12", formato.PrintFechaHora(cinconoviembre, formato.DMA))
	assert.Equal(t, "11/05/2021 11:12", formato.PrintFechaHora(cinconoviembre, formato.MDA))
	assert.Equal(t, "2021-11-05 11:12", formato.PrintFechaHora(cinconoviembre, formato.AMD))
	assert.Equal(t, "2021-11-05", formato.PrintFecha(cinconoviembre, formato.ISO))
	assert.Equal(t, "05/11/2021", formato.PrintFecha(cinconoviembre, formato.DMA))
	assert.Equal(t, "11/05/2021", formato.PrintFecha(cinconoviembre, formato.MDA))
	assert.Equal(t, "11:12", formato.PrintHora(cinconoviembre, false))
}

func ExampleParseFechaHora() {
	fh, err := formato.ParseFechaHora("12/10/2023 12:45", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintFechaHora(fh, formato.ISO))
	// Output: 2023-10-12T12:45:00
}

func ExamplePrintFechaHora() {
	fh, err := formato.ParseFechaHora("12/10/2023 12:45", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintFechaHora(fh, formato.ISO))
	// Output: 2023-10-12T12:45:00
}

func ExamplePrintFecha() {
	fh, err := formato.ParseFechaHora("12/10/2023 12:45", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintFecha(fh, formato.DMA))
	// Output: 12/10/2023
}

func TestMustParseFechaHora(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	obtiene := formato.MustParseFechaHora("2021-11-05T11:12:13")
	assert.Equal(t, fecha, obtiene)
	fecha = time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC)
	obtiene = formato.MustParseFechaHora("2021-11-05")
	assert.Equal(t, fecha, obtiene)
}

func ExampleMustParseFechaHora() {
	fh := formato.MustParseFechaHora("2023-10-12T12:45:00")
	fmt.Println(formato.PrintFechaHora(fh, formato.DMA))
	// Output: 12/10/2023 12:45
}

func TestParseHora(t *testing.T) {
	h, err := formato.ParseHora("15:16:17")
	assert.Equal(t, "15:16", formato.PrintHora(h, false))
	assert.NoError(t, err)
	_, err = formato.ParseHora("05-11-2021 15:16:17")
	assert.NotNil(t, err)
}

func ExampleParseHora() {
	h, err := formato.ParseHora("12:45")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintHora(h, false))
	// Output: 12:45
}

func ExamplePrintHora() {
	h, err := formato.ParseHora("12:34:56")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintHora(h, false))
	fmt.Println(formato.PrintHora(h, true))
	// Output:
	// 12:34
	// 12:34:56
}

func ExampleMustParseHora() {
	h := formato.MustParseHora("15:16:17")
	fmt.Println(formato.PrintHora(h, true))
	// Output: 15:16:17
}

func TestTimestamp(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	ISO := formato.PrintFechaHora(fecha, formato.ISO)
	obtiene, err := formato.ParseTimestamp("", formato.ISO)
	assert.NoError(t, err)
	assert.False(t, obtiene.Valid)
	assert.Equal(t, "", formato.PrintTimestamp(obtiene, formato.ISO))
	obtiene, err = formato.ParseTimestamp(ISO, formato.ISO)
	assert.NoError(t, err)
	assert.True(t, obtiene.Valid)
	assert.Equal(t, fecha, obtiene.Time)
	assert.Equal(t, ISO, formato.PrintTimestamp(obtiene, formato.ISO))
}

func ExamplePrintTimestamp() {
	ts, err := formato.ParseTimestamp("25/4/1974 22:55", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintTimestamp(ts, formato.AMD))
	// Output: 1974-04-25 22:55
}

func ExampleParseTimestamp() {
	ts, err := formato.ParseTimestamp("25/4/1974 22:55", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintTimestamp(ts, formato.AMD))
	// Output: 1974-04-25 22:55
}

func TestDate(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	dia := time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC)
	ISO := formato.PrintFecha(fecha, formato.ISO)
	obtiene, err := formato.ParseDate("", formato.ISO)
	assert.NoError(t, err)
	assert.False(t, obtiene.Valid)
	assert.Equal(t, "", formato.PrintDate(obtiene, formato.ISO))
	obtiene, err = formato.ParseDate(ISO, formato.ISO)
	assert.NoError(t, err)
	assert.True(t, obtiene.Valid)
	assert.Equal(t, dia, obtiene.Time)
	assert.Equal(t, ISO, formato.PrintDate(obtiene, formato.ISO))
}

func ExamplePrintDate() {
	d, err := formato.ParseDate("25/4/1974", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintDate(d, formato.AMD))
	// Output: 1974-04-25
}

func ExampleParseDate() {
	d, err := formato.ParseDate("25/4/1974", formato.DMA)
	errores.PanicIfError(err)
	fmt.Println(formato.PrintDate(d, formato.AMD))
	// Output: 1974-04-25
}

func TestTime(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	hora := formato.PrintHora(fecha, false)
	obtiene, err := formato.ParseTime("")
	assert.NoError(t, err)
	assert.False(t, obtiene.Valid)
	assert.Equal(t, "", formato.PrintTime(obtiene, false))
	usec := int64((fecha.Hour()*3600 + fecha.Minute()*60) * 1000000)
	obtiene, err = formato.ParseTime(hora)
	assert.NoError(t, err)
	assert.True(t, obtiene.Valid)
	assert.Equal(t, usec, obtiene.Microseconds)
	assert.Equal(t, hora, formato.PrintTime(obtiene, false))
}

func ExamplePrintTime() {
	t, err := formato.ParseTime("22:55")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintTime(t, false))
	// Output: 22:55
}

func ExampleParseTime() {
	t, err := formato.ParseTime("22:55")
	errores.PanicIfError(err)
	fmt.Println(formato.PrintTime(t, false))
	// Output: 22:55
}
