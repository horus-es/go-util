package formato

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgtype"
)

var reDuration = regexp.MustCompile(`^\s*([0-9]+)\s*([hms])[^0-9]*$`)

// Parsea una duración a time.Duration. Soporta H horas, M minutos o S segundos.
func ParseDuracion(s string) (result time.Duration, err error) {
	partes := reDuration.FindStringSubmatch(strings.ToLower(s))
	if len(partes) != 3 {
		err = fmt.Errorf("duración %q no soportada", s)
		return
	}
	n, _ := strconv.Atoi(partes[1])
	switch partes[2] {
	case "h":
		result = time.Duration(n) * time.Hour
	case "m":
		result = time.Duration(n) * time.Minute
	case "s":
		result = time.Duration(n) * time.Second
	}
	return
}

var reInterval = regexp.MustCompile(`^\s*([0-9]+)\s*([yamwsd])[^0-9]*$`)

// Parsea una duración a pgtype.Interval. Soporta A años, M meses, S semanas o D dias.
func ParseIntervalAMSD(s string) (result pgtype.Interval, err error) {
	if strings.TrimSpace(s) == "" {
		result.Status = pgtype.Null
		return
	}
	partes := reInterval.FindStringSubmatch(strings.ToLower(s))
	if len(partes) != 3 {
		err = fmt.Errorf("duración %q no soportada", s)
		return
	}
	n, _ := strconv.Atoi(partes[1])
	switch partes[2] {
	case "a", "y":
		result.Months = int32(n * 12)
	case "m":
		result.Months = int32(n)
	case "s", "w":
		result.Days = int32(n * 7)
	case "d":
		result.Days = int32(n)
	}
	result.Status = pgtype.Present
	return
}

// Parsea una duración a pgtype.Interval. Soporta H horas, M minutos o S segundos.
func ParseIntervalHMS(s string) (result pgtype.Interval, err error) {
	if strings.TrimSpace(s) == "" {
		result.Status = pgtype.Null
		return
	}
	partes := reDuration.FindStringSubmatch(strings.ToLower(s))
	if len(partes) != 3 {
		err = fmt.Errorf("duración %q no soportada", s)
		return
	}
	n, _ := strconv.Atoi(partes[1])
	switch partes[2] {
	case "h":
		result.Microseconds = int64(n) * 3600000000
	case "m":
		result.Microseconds = int64(n) * 60000000
	case "s":
		result.Microseconds = int64(n) * 1000000
	}
	result.Status = pgtype.Present
	return
}

// Imprime una duración en formato X horas, X minutos o X segundos.
func PrintDuracion(z time.Duration) string {
	z = z / time.Second
	if z%60 != 0 {
		return fmt.Sprint(int(z), "s")
	}
	z = z / 60
	if z%60 != 0 {
		return fmt.Sprint(int(z), "m")
	}
	z = z / 60
	if z != 0 {
		return fmt.Sprint(int(z), "h")
	}
	return ""
}

// Imprime una duración en formato ISO8601, util para órdenes SQL
func PrintDuracionIso(z time.Duration) string {
	s := int64(z / time.Second)
	return fmt.Sprint("PT", s, "S")
}

// Imprime una duración en formato A años, M meses, S semanas, D dias, H horas, M minutos, S segundos.
func PrintInterval(z pgtype.Interval) string {
	if z.Status != pgtype.Present {
		return ""
	}
	partes := []string{}
	if z.Months > 0 {
		partes = append(partes, fmt.Sprint(z.Months, "m"))
	}
	if z.Days > 0 {
		partes = append(partes, fmt.Sprint(z.Days, "d"))
	}
	s := PrintDuracion(time.Duration(z.Microseconds) * time.Microsecond)
	if s != "" {
		partes = append(partes, s)
	}
	if len(partes) == 0 {
		return ""
	}
	return strings.Join(partes, " ")
}

// Imprime una duración en formato ISO8601, util en órdenes SQL
func PrintIntervalIso(z pgtype.Interval) string {
	if z.Status != pgtype.Present {
		return ""
	}
	result := "P"
	if z.Months > 0 {
		result += fmt.Sprint(z.Months, "M")
	}
	if z.Days > 0 {
		result += fmt.Sprint(z.Days, "D")
	}
	s := z.Microseconds / 1000000
	if s > 0 {
		result += fmt.Sprint("T", s, "S")
	}
	if len(result) == 1 {
		return ""
	}
	return result
}

// Añade un pgtype.Interval a un pgtype.Timestamp
func AddInterval(t pgtype.Timestamp, i pgtype.Interval) (result pgtype.Timestamp) {
	if t.Status != pgtype.Present || i.Status != pgtype.Present {
		result.Status = pgtype.Null
		return
	}
	result.Time = t.Time.AddDate(0, int(i.Months), int(i.Days)).Add(time.Duration(i.Microseconds) * time.Microsecond)
	result.Status = pgtype.Present
	return
}

// Sustrae un pgtype.Interval de un pgtype.Timestamp
func SubInterval(t pgtype.Timestamp, i pgtype.Interval) (result pgtype.Timestamp) {
	if t.Status != pgtype.Present || i.Status != pgtype.Present {
		result.Status = pgtype.Null
		return
	}
	result.Time = t.Time.AddDate(0, -int(i.Months), -int(i.Days)).Add(-time.Duration(i.Microseconds) * time.Microsecond)
	result.Status = pgtype.Present
	return
}
