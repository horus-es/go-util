package formato

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/horus-es/go-util/v3/misc"
	"github.com/jackc/pgx/v5/pgtype"
)

type IsoDuration struct {
	Years   float64
	Months  float64
	Weeks   float64
	Days    float64
	Hours   float64
	Minutes float64
	Seconds float64
}

// Parsea una duración en formato ISO8601
func ParseIsoDuration(s string) (*IsoDuration, error) {
	state := 'X'
	D := &IsoDuration{}
	num := ""
	for _, char := range s {
		switch char {
		case 'P':
			if state != 'X' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			state = 'P'
		case 'T':
			if state != 'P' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			state = 'T'
		case 'Y':
			if state != 'P' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			x, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			D.Years += x
			num = ""
		case 'M':
			switch state {
			case 'P':
				x, err := strconv.ParseFloat(num, 64)
				if err != nil {
					return nil, fmt.Errorf("error en duración: %s", s)
				}
				D.Months += x
				num = ""
			case 'T':
				x, err := strconv.ParseFloat(num, 64)
				if err != nil {
					return nil, fmt.Errorf("error en duración: %s", s)
				}
				D.Minutes += x
				num = ""
			default:
				return nil, fmt.Errorf("error en duración: %s", s)
			}
		case 'W':
			if state != 'P' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			x, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			D.Weeks += x
			num = ""
		case 'D':
			if state != 'P' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			x, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			D.Days += x
			num = ""
		case 'H':
			if state != 'T' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			x, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			D.Hours += x
			num = ""
		case 'S':
			if state != 'T' {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			x, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, fmt.Errorf("error en duración: %s", s)
			}
			D.Seconds += x
			num = ""
		default:
			if (char >= '0' && char <= '9') || char == '.' {
				num += string(char)
				continue
			}
			return nil, fmt.Errorf("error en duración: %s", s)
		}
	}
	if num != "" {
		return nil, fmt.Errorf("error en duración: %s", s)
	}
	return D, nil
}

var reIsoDuration = regexp.MustCompile(`([0-9\.]+)\s*([a-z]+)`)

// Parsea una duración especificada como 3 horas y 23 minutos, o 3h23min
func ParseHumanDuration(s string, m_is_month bool) (*IsoDuration, error) {
	D := &IsoDuration{}
	limpio := strings.TrimSpace(s)
	if limpio == "" {
		return D, nil
	}
	limpio = strings.ToLower(misc.QuitaAcentos(limpio))
	matches := reIsoDuration.FindAllStringSubmatch(limpio, -1)
	if matches == nil {
		return nil, fmt.Errorf("error en duración: %s", s)
	}
	for _, m := range matches {
		n, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return nil, fmt.Errorf("error en duración: %s", s)
		}
		switch m[2] {
		case "y", "year", "years", "a", "ano", "anos": // TODO: Toma 0,3 años como 3 años
			D.Years += n
		case "m":
			if m_is_month {
				D.Months += n
			} else {
				D.Minutes += n
			}
		case "month", "months", "mes", "meses":
			D.Months += n
		case "w", "week", "weeks", "s", "semana", "semanas":
			D.Weeks += n
		case "d", "day", "days", "dia", "dias":
			D.Days += n
		case "h", "hour", "hours", "hora", "horas":
			D.Hours += n
		case "min", "minute", "minutes", "minuto", "minutos":
			D.Minutes += n
		case "secs", "sec", "seconds", "second", "segundos", "segundo":
			D.Seconds += n
		default:
			return nil, fmt.Errorf("error en duración: %s", s)
		}
	}
	return D, nil
}

func (D IsoDuration) String() string {
	return PrintIsoDuration(&D)
}

func PrintIsoDuration(D *IsoDuration) string {
	fecha := "P"
	if D.Years > 0 {
		fecha += strconv.FormatFloat(D.Years, 'f', -1, 64) + "Y"
	}
	if D.Months > 0 {
		fecha += strconv.FormatFloat(D.Months, 'f', -1, 64) + "M"
	}
	if D.Weeks > 0 {
		fecha += strconv.FormatFloat(D.Weeks, 'f', -1, 64) + "W"
	}
	if D.Days > 0 {
		fecha += strconv.FormatFloat(D.Days, 'f', -1, 64) + "D"
	}
	hora := "T"
	if D.Hours > 0 {
		hora += strconv.FormatFloat(D.Hours, 'f', -1, 64) + "H"
	}
	if D.Minutes > 0 {
		hora += strconv.FormatFloat(D.Minutes, 'f', -1, 64) + "M"
	}
	if D.Seconds > 0 {
		hora += strconv.FormatFloat(D.Seconds, 'f', -1, 64) + "S"
	}

	if len(hora) > 1 {
		return fecha + hora
	}
	if len(fecha) > 1 {
		return fecha
	}
	return ""
}

func (D IsoDuration) Spanish() string {
	return PrintSpanishDuration(&D)
}

func PrintSpanishDuration(D *IsoDuration) string {
	var partes []string
	if D.Years > 0 && D.Years != 1 {
		partes = append(partes, strconv.FormatFloat(D.Years, 'f', -1, 64)+" años")
	}
	if D.Years == 1 {
		partes = append(partes, "1 año")
	}
	if D.Months > 0 && D.Months != 1 {
		partes = append(partes, strconv.FormatFloat(D.Months, 'f', -1, 64)+" meses")
	}
	if D.Months == 1 {
		partes = append(partes, "1 mes")
	}
	if D.Weeks > 0 && D.Weeks != 1 {
		partes = append(partes, strconv.FormatFloat(D.Weeks, 'f', -1, 64)+" semanas")
	}
	if D.Weeks == 1 {
		partes = append(partes, "1 semana")
	}
	if D.Days > 0 && D.Days != 1 {
		partes = append(partes, strconv.FormatFloat(D.Days, 'f', -1, 64)+" dias")
	}
	if D.Days == 1 {
		partes = append(partes, "1 día")
	}
	if D.Hours > 0 && D.Hours != 1 {
		partes = append(partes, strconv.FormatFloat(D.Hours, 'f', -1, 64)+" horas")
	}
	if D.Hours == 1 {
		partes = append(partes, "1 hora")
	}
	if D.Minutes > 0 && D.Minutes != 1 {
		partes = append(partes, strconv.FormatFloat(D.Minutes, 'f', -1, 64)+" minutos")
	}
	if D.Minutes == 1 {
		partes = append(partes, "1 minuto")
	}
	if D.Seconds > 0 && D.Seconds != 1 {
		partes = append(partes, strconv.FormatFloat(D.Seconds, 'f', -1, 64)+" segundos")
	}
	if D.Seconds == 1 {
		partes = append(partes, "1 segundo")
	}
	result := ""
	for k, p := range partes {
		result += p
		switch {
		case k == len(partes)-2:
			result += " y "
		case k < len(partes)-2:
			result += ", "
		}
	}
	return result
}

// Atención: Se ignora el hecho de que el cambio de hora puede hacer los dias mas largos o cortos, siempre se consideran de 24H.
func (D IsoDuration) ToNative() (*time.Duration, error) {
	return IsoDurationToNative(&D)
}

// Atención: Se ignora el hecho de que el cambio de hora puede hacer los dias mas largos o cortos, siempre se consideran de 24H.
func IsoDurationToNative(D *IsoDuration) (*time.Duration, error) {
	if D.Years > 0 || D.Months > 0 {
		return nil, fmt.Errorf("no se soportan duraciones de años o meses")
	}
	segundos := D.Weeks*7*24*60*60 + D.Days*24*60*60 + D.Hours*60*60 + D.Minutes*60 + D.Seconds
	native := time.Duration(segundos * float64(time.Second))
	return &native, nil
}

func IsoDurationToInterval(D *IsoDuration) pgtype.Interval {
	var interval pgtype.Interval
	interval.Months = int32(D.Years*12 + D.Months)
	interval.Days = int32(D.Weeks*7 + D.Days)
	interval.Microseconds = int64(1000000 * (D.Seconds + D.Minutes*60 + D.Hours*3600))
	interval.Valid = true
	return interval
}

func (D IsoDuration) AddToNative(t time.Time) (*time.Time, error) {
	return AddIsoDurationToNative(&D, t)
}

func AddIsoDurationToNative(D *IsoDuration, t time.Time) (*time.Time, error) {
	if D.Years != math.Trunc(D.Years) || D.Months != math.Trunc(D.Months) || D.Weeks != math.Trunc(D.Weeks) || D.Days != math.Trunc(D.Days) {
		return nil, fmt.Errorf("no se soportan duraciones fraccionarias")
	}
	dias := D.Weeks*7 + D.Days
	segundos := D.Hours*60*60 + D.Minutes*60 + D.Seconds
	t = t.AddDate(int(D.Years), int(D.Months), int(dias)).Add(time.Duration(segundos * float64(time.Second)))
	return &t, nil
}

func (D IsoDuration) SubtractFromNative(t time.Time) (*time.Time, error) {
	return SubtractIsoDurationFromNative(&D, t)
}

func SubtractIsoDurationFromNative(D *IsoDuration, t time.Time) (*time.Time, error) {
	if D.Years != math.Trunc(D.Years) || D.Months != math.Trunc(D.Months) || D.Weeks != math.Trunc(D.Weeks) || D.Days != math.Trunc(D.Days) {
		return nil, fmt.Errorf("no se soportan duraciones fraccionarias")
	}
	dias := D.Weeks*7 + D.Days
	segundos := D.Hours*60*60 + D.Minutes*60 + D.Seconds
	t = t.AddDate(int(-D.Years), int(-D.Months), int(-dias)).Add(time.Duration(-segundos) * time.Second)
	return &t, nil
}

func (D IsoDuration) IsZero() bool {
	return IsZeroDuration(&D)
}

func IsZeroDuration(D *IsoDuration) bool {
	return D.Years == 0 && D.Months == 0 && D.Weeks == 0 && D.Days == 0 && D.Hours == 0 && D.Minutes == 0 && D.Seconds == 0
}
