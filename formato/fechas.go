package formato

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/horus-es/go-util/errores"
	"github.com/jackc/pgtype"
)

// https://en.wikipedia.org/wiki/Date_format_by_country
type FormatoFecha string

const (
	ISO FormatoFecha = "ISO" // RFC3339
	DMA FormatoFecha = "DMA" // Dia/Mes/Año
	MDA FormatoFecha = "MDA" // Mes/Dia/Año
	AMD FormatoFecha = "AMD" // Año-Mes-Dia
)

var reParseFecha1 = regexp.MustCompile(`^([1-9])/`)
var reParseFecha2 = regexp.MustCompile(`/([1-9])/`)
var reParseFecha3 = regexp.MustCompile(`/([1-9])( |$)`)

// Parsea una fecha con hora opcional
// Admite varios tipos de fecha: ff=ISO (ISO A-M-D), ff=MDA (USA M/D/A), ff=AMD (internacional A-M-D), ff=DMA (resto D/M/A)
// Para MDA, AMD y DMA también se admiten . y - como separador de fecha, el separador de hora siempre es :
func ParseFechaHora(s string, ff FormatoFecha) (result time.Time, err error) {
	switch ff {
	case DMA:
		t := strings.ReplaceAll(s, "-", "/")
		t = strings.ReplaceAll(t, ".", "/")
		t = reParseFecha1.ReplaceAllString(t, "0$1/")
		t = reParseFecha2.ReplaceAllString(t, "/0$1/")
		result, err = time.Parse("02/01/2006 15:04:05", t)
		if err == nil {
			return
		}
		result, err = time.Parse("02/01/2006 15:04", t)
		if err == nil {
			return
		}
		result, err = time.Parse("02/01/2006", t)
		if err == nil {
			return
		}
	case MDA:
		t := strings.ReplaceAll(s, "-", "/")
		t = strings.ReplaceAll(t, ".", "/")
		t = reParseFecha1.ReplaceAllString(t, "0$1/")
		t = reParseFecha2.ReplaceAllString(t, "/0$1/")
		result, err = time.Parse("01/02/2006 15:04:05", t)
		if err == nil {
			return
		}
		result, err = time.Parse("01/02/2006 15:04", t)
		if err == nil {
			return
		}
		result, err = time.Parse("01/02/2006", t)
		if err == nil {
			return
		}
	case AMD:
		t := strings.ReplaceAll(s, "-", "/")
		t = strings.ReplaceAll(t, ".", "/")
		t = reParseFecha3.ReplaceAllString(t, "/0$1")
		t = reParseFecha2.ReplaceAllString(t, "/0$1/")
		result, err = time.Parse("2006/01/02 15:04:05", t)
		if err == nil {
			return
		}
		result, err = time.Parse("2006/01/02 15:04", t)
		if err == nil {
			return
		}
		result, err = time.Parse("2006/01/02", t)
		if err == nil {
			return
		}
	}
	result, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return
	}
	result, err = time.Parse(time.RFC3339, s+"Z")
	if err == nil {
		return
	}
	result, err = time.Parse(time.RFC3339, s+"T00:00:00Z")
	if err == nil {
		return
	}
	err = fmt.Errorf("fecha %q no reconocida", s)
	return
}

// Parsea una fecha con hora opcional en formato ISO
func MustParseFechaHora(s string) time.Time {
	result, err := ParseFechaHora(s, ISO)
	errores.PanicIfError(err)
	return result
}

// Parsea una hora
func ParseHora(s string) (result time.Time, err error) {
	result, err = time.Parse("15:04:05", s)
	if err == nil {
		return
	}
	result, err = time.Parse("15:04", s)
	if err == nil {
		return
	}
	err = fmt.Errorf("hora %q no reconocida", s)
	return
}

// Parsea una hora
func MustParseHora(s string) time.Time {
	result, err := ParseHora(s)
	errores.PanicIfError(err)
	return result
}

// Parsea una fecha+hora a pgtype.Timestamp. Mismas consideraciones que ParseFechaHora. Los vacios se consideran null.
func ParseTimestamp(s string, ff FormatoFecha) (result pgtype.Timestamp, err error) {
	if s == "" {
		result.Status = pgtype.Null
		return
	}
	fh, err := ParseFechaHora(s, ff)
	if err == nil {
		result.Time = fh
		result.Status = pgtype.Present
	}
	return
}

// Parsea una fecha a pgtype.Date. Mismas consideraciones que ParseFechaHora. Los vacios se consideran null.
func ParseDate(s string, ff FormatoFecha) (result pgtype.Date, err error) {
	if s == "" {
		result.Status = pgtype.Null
		return
	}
	fh, err := ParseFechaHora(s, ff)
	if err == nil {
		result.Time = fh
		result.Status = pgtype.Present
	}
	return
}

// Parsea una hora a pgtype.Time. Mismas consideraciones que ParseFechaHora. Los vacios se consideran null.
func ParseTime(s string) (result pgtype.Time, err error) {
	if s == "" {
		result.Status = pgtype.Null
		return
	}
	fh, err := ParseHora(s)
	if err == nil {
		medianoche := time.Date(fh.Year(), fh.Month(), fh.Day(), 0, 0, 0, 0, fh.Location())
		result.Microseconds = fh.Sub(medianoche).Microseconds()
		result.Status = pgtype.Present
	}
	return
}

// Imprime una fecha+hora
func PrintFechaHora(fh time.Time, ff FormatoFecha) string {
	switch ff {
	case DMA:
		return fh.Format("02/01/2006 15:04")
	case MDA:
		return fh.Format("01/02/2006 15:04")
	case AMD:
		return fh.Format("2006-01-02 15:04")
	default:
		return fh.Format("2006-01-02T15:04:05")
	}
}

// Imprime una fecha
func PrintFecha(fh time.Time, ff FormatoFecha) string {
	switch ff {
	case DMA:
		return fh.Format("02/01/2006")
	case MDA:
		return fh.Format("01/02/2006")
	default:
		return fh.Format("2006-01-02")
	}
}

// Imprime una hora
func PrintHora(fh time.Time, secs bool) string {
	if secs {
		return fh.Format("15:04:05")
	} else {
		return fh.Format("15:04")
	}
}

// Imprime un pgtype.Timestamp
func PrintTimestamp(fh pgtype.Timestamp, ff FormatoFecha) string {
	if fh.Status != pgtype.Present {
		return ""
	}
	return PrintFechaHora(fh.Time, ff)
}

// Imprime un pgtype.Date
func PrintDate(fh pgtype.Date, ff FormatoFecha) string {
	if fh.Status != pgtype.Present {
		return ""
	}
	return PrintFecha(fh.Time, ff)
}

// Imprime un pgtype.Time
func PrintTime(fh pgtype.Time, secs bool) string {
	if fh.Status != pgtype.Present {
		return ""
	}
	f := time.Date(2000, 01, 01, 0, 0, 0, int(fh.Microseconds)*1000, time.UTC)
	return PrintHora(f, secs)
}
