package formato

import (
	"fmt"
	"testing"
	"time"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/stretchr/testify/assert"
)

func compruebaParseFechaHora(t *testing.T, s string, ff Fecha, espera time.Time) {
	obtiene, err := ParseFechaHora(s, ff)
	assert.NoError(t, err)
	assert.Equal(t, espera, obtiene)
}

func TestParseFechaHora(t *testing.T) {
	cinconoviembre := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	cinconoviembre0 := time.Date(2021, 11, 5, 11, 12, 0, 0, time.UTC)
	cinconoviembre00 := time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC)
	cincoagosto00 := time.Date(2021, 8, 5, 0, 0, 0, 0, time.UTC)
	compruebaParseFechaHora(t, "2021-11-05T11:12:13Z", ISO, cinconoviembre)
	compruebaParseFechaHora(t, "2021-11-05T11:12:13", ISO, cinconoviembre)
	compruebaParseFechaHora(t, "2021-11-05 11:12:13", AMD, cinconoviembre)
	compruebaParseFechaHora(t, "2021-11-05", ISO, cinconoviembre00)
	compruebaParseFechaHora(t, "05/11/2021 11:12:13", DMA, cinconoviembre)
	compruebaParseFechaHora(t, "05/11/2021 11:12", DMA, cinconoviembre0)
	compruebaParseFechaHora(t, "05/11/2021", DMA, cinconoviembre00)
	compruebaParseFechaHora(t, "11/05/2021 11:12:13", MDA, cinconoviembre)
	compruebaParseFechaHora(t, "11/05/2021 11:12", MDA, cinconoviembre0)
	compruebaParseFechaHora(t, "2021-11-05 11:12", AMD, cinconoviembre0)
	compruebaParseFechaHora(t, "11/05/2021", MDA, cinconoviembre00)
	compruebaParseFechaHora(t, "11-05-2021", MDA, cinconoviembre00)
	compruebaParseFechaHora(t, "11.05.2021", MDA, cinconoviembre00)
	compruebaParseFechaHora(t, "5.8-2021", DMA, cincoagosto00)
	compruebaParseFechaHora(t, "2021-8-5", AMD, cincoagosto00)
	_, err := ParseFechaHora("05-11-2021 15:16:17", ISO)
	assert.NotNil(t, err)
}

func TestPrintFechaHora(t *testing.T) {
	cinconoviembre := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	assert.Equal(t, "2021-11-05T11:12:13", PrintFechaHora(cinconoviembre, ISO))
	assert.Equal(t, "05/11/2021 11:12", PrintFechaHora(cinconoviembre, DMA))
	assert.Equal(t, "11/05/2021 11:12", PrintFechaHora(cinconoviembre, MDA))
	assert.Equal(t, "2021-11-05 11:12", PrintFechaHora(cinconoviembre, AMD))
	assert.Equal(t, "2021-11-05", PrintFecha(cinconoviembre, ISO))
	assert.Equal(t, "05/11/2021", PrintFecha(cinconoviembre, DMA))
	assert.Equal(t, "11/05/2021", PrintFecha(cinconoviembre, MDA))
	assert.Equal(t, "11:12", PrintHora(cinconoviembre, false))
}

func ExampleParseFechaHora() {
	fh, err := ParseFechaHora("12/10/2023 12:45", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintFechaHora(fh, ISO))
	// Output: 2023-10-12T12:45:00
}

func ExamplePrintFechaHora() {
	fh, err := ParseFechaHora("12/10/2023 12:45", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintFechaHora(fh, ISO))
	// Output: 2023-10-12T12:45:00
}

func ExamplePrintFecha() {
	fh, err := ParseFechaHora("12/10/2023 12:45", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintFecha(fh, DMA))
	// Output: 12/10/2023
}

func TestMustParseFechaHora(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	obtiene := MustParseFechaHora("2021-11-05T11:12:13")
	assert.Equal(t, fecha, obtiene)
	fecha = time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC)
	obtiene = MustParseFechaHora("2021-11-05")
	assert.Equal(t, fecha, obtiene)
}

func ExampleMustParseFechaHora() {
	fh := MustParseFechaHora("2023-10-12T12:45:00")
	fmt.Println(PrintFechaHora(fh, DMA))
	// Output: 12/10/2023 12:45
}

func TestParseHora(t *testing.T) {
	h, err := ParseHora("15:16:17")
	assert.Equal(t, "15:16", PrintHora(h, false))
	assert.NoError(t, err)
	_, err = ParseHora("05-11-2021 15:16:17")
	assert.NotNil(t, err)
}

func ExampleParseHora() {
	h, err := ParseHora("12:45")
	errores.PanicIfError(err)
	fmt.Println(PrintHora(h, false))
	// Output: 12:45
}

func ExamplePrintHora() {
	h, err := ParseHora("12:34:56")
	errores.PanicIfError(err)
	fmt.Println(PrintHora(h, false))
	fmt.Println(PrintHora(h, true))
	// Output:
	// 12:34
	// 12:34:56
}

func ExampleMustParseHora() {
	h := MustParseHora("15:16:17")
	fmt.Println(PrintHora(h, true))
	// Output: 15:16:17
}

func TestTimestamp(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	iso := PrintFechaHora(fecha, ISO)
	obtiene, err := ParseTimestamp("", ISO)
	assert.NoError(t, err)
	assert.False(t, obtiene.Valid)
	assert.Equal(t, "", PrintTimestamp(obtiene, ISO))
	obtiene, err = ParseTimestamp(iso, ISO)
	assert.NoError(t, err)
	assert.True(t, obtiene.Valid)
	assert.Equal(t, fecha, obtiene.Time)
	assert.Equal(t, iso, PrintTimestamp(obtiene, ISO))
}

func ExamplePrintTimestamp() {
	ts, err := ParseTimestamp("25/4/1974 22:55", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintTimestamp(ts, AMD))
	// Output: 1974-04-25 22:55
}

func ExampleParseTimestamp() {
	ts, err := ParseTimestamp("25/4/1974 22:55", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintTimestamp(ts, AMD))
	// Output: 1974-04-25 22:55
}

func TestDate(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	dia := time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC)
	iso := PrintFecha(fecha, ISO)
	obtiene, err := ParseDate("", ISO)
	assert.NoError(t, err)
	assert.False(t, obtiene.Valid)
	assert.Equal(t, "", PrintDate(obtiene, ISO))
	obtiene, err = ParseDate(iso, ISO)
	assert.NoError(t, err)
	assert.True(t, obtiene.Valid)
	assert.Equal(t, dia, obtiene.Time)
	assert.Equal(t, iso, PrintDate(obtiene, ISO))
}

func ExamplePrintDate() {
	d, err := ParseDate("25/4/1974", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintDate(d, AMD))
	// Output: 1974-04-25
}

func ExampleParseDate() {
	d, err := ParseDate("25/4/1974", DMA)
	errores.PanicIfError(err)
	fmt.Println(PrintDate(d, AMD))
	// Output: 1974-04-25
}

func TestTime(t *testing.T) {
	fecha := time.Date(2021, 11, 5, 11, 12, 13, 0, time.UTC)
	hora := PrintHora(fecha, false)
	obtiene, err := ParseTime("")
	assert.NoError(t, err)
	assert.False(t, obtiene.Valid)
	assert.Equal(t, "", PrintTime(obtiene, false))
	usec := int64((fecha.Hour()*3600 + fecha.Minute()*60) * 1000000)
	obtiene, err = ParseTime(hora)
	assert.NoError(t, err)
	assert.True(t, obtiene.Valid)
	assert.Equal(t, usec, obtiene.Microseconds)
	assert.Equal(t, hora, PrintTime(obtiene, false))
}

func ExamplePrintTime() {
	t, err := ParseTime("22:55")
	errores.PanicIfError(err)
	fmt.Println(PrintTime(t, false))
	// Output: 22:55
}

func ExampleParseTime() {
	t, err := ParseTime("22:55")
	errores.PanicIfError(err)
	fmt.Println(PrintTime(t, false))
	// Output: 22:55
}
