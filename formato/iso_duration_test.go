package formato_test

import (
	"testing"
	"time"

	"github.com/horus-es/go-util/v3/formato"
	"github.com/stretchr/testify/assert"
)

func TestParseIsoDuration(t *testing.T) {
	tests := []struct {
		name        string
		isoDuration string
		spanish     string
		wantErr     bool
	}{
		{
			name:        "zero duration",
			isoDuration: "",
			spanish:     "",
		},
		{
			name:        "single year",
			isoDuration: "P1Y",
			spanish:     "1 año",
		},
		{
			name:        "multiple years",
			isoDuration: "P2.5Y",
			spanish:     "2.5 años",
		},
		{
			name:        "single month",
			isoDuration: "P1M",
			spanish:     "1 mes",
		},
		{
			name:        "multiple months",
			isoDuration: "P3M",
			spanish:     "3 meses",
		},
		{
			name:        "single week",
			isoDuration: "P1W",
			spanish:     "1 semana",
		},
		{
			name:        "multiple weeks",
			isoDuration: "P2W",
			spanish:     "2 semanas",
		},
		{
			name:        "single day",
			isoDuration: "P1D",
			spanish:     "1 día",
		},
		{
			name:        "multiple days",
			isoDuration: "P5D",
			spanish:     "5 dias",
		},
		{
			name:        "single hour",
			isoDuration: "PT1H",
			spanish:     "1 hora",
		},
		{
			name:        "multiple hours",
			isoDuration: "PT3H",
			spanish:     "3 horas",
		},
		{
			name:        "single minute",
			isoDuration: "PT1M",
			spanish:     "1 minuto",
		},
		{
			name:        "multiple minutes",
			isoDuration: "PT30M",
			spanish:     "30 minutos",
		},
		{
			name:        "single second",
			isoDuration: "PT1S",
			spanish:     "1 segundo",
		},
		{
			name:        "multiple seconds",
			isoDuration: "PT45S",
			spanish:     "45 segundos",
		},
		{
			name:        "combined duration",
			isoDuration: "P1Y2M3DT4H5M6S",
			spanish:     "1 año, 2 meses, 3 dias, 4 horas, 5 minutos y 6 segundos",
		},
		{
			name:        "years and hours",
			isoDuration: "P2YT12H",
			spanish:     "2 años y 12 horas",
		},
		{
			name:        "error",
			isoDuration: "PYT12H",
			spanish:     "12 ñordos",
			wantErr:     true,
		},
		{
			name:        "error T2",
			isoDuration: "T2",
			spanish:     "T2",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err := formato.ParseIsoDuration(tt.isoDuration)
				assert.Error(t, err)
			} else {
				got, err := formato.ParseIsoDuration(tt.isoDuration)
				assert.NoError(t, err)
				assert.Equal(t, tt.spanish, got.Spanish())
			}
			if tt.wantErr {
				_, err := formato.ParseHumanDuration(tt.spanish, false)
				assert.Error(t, err)
			} else {
				got, err := formato.ParseHumanDuration(tt.spanish, false)
				assert.NoError(t, err)
				assert.Equal(t, tt.isoDuration, got.String())
			}
		})
	}
}

func TestIsoDuration_AddToNative(t *testing.T) {
	ahora := time.Now()
	d, err := formato.ParseIsoDuration("PT24H")
	assert.NoError(t, err)
	mañana, err := d.AddToNative(ahora)
	assert.NoError(t, err)
	assert.Equal(t, ahora.AddDate(0, 0, 1), *mañana)
	ayer, err := d.SubtractFromNative(ahora)
	assert.NoError(t, err)
	assert.Equal(t, ahora.AddDate(0, 0, -1), *ayer)
}
