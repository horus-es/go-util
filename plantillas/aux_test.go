package plantillas_test

// Datos y funciones auxiliares para los tests

import (
	"bytes"
	"hash/crc32"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type tParte struct {
	VatID   string `json:",omitempty"`
	Name    string `json:",omitempty"`
	Address string `json:",omitempty"`
}

type tLinea struct {
	Price   float64 `json:",omitempty"`
	Service string  `json:",omitempty"`
}

type tStay struct {
	Plate     string    `json:",omitempty"`
	Ticket    string    `json:",omitempty"`
	EntryDate time.Time `json:",omitempty"`
	Duration  string    `json:",omitempty"`
	Price     float64   `json:",omitempty"`
}

type tTransaction struct {
	Date   time.Time `json:",omitempty"`
	Method string    `json:",omitempty"`
}

type tFactura struct {
	ID          string    `json:",omitempty"`
	IssueDate   time.Time `json:",omitempty"`
	Reason      string    `json:",omitempty"`
	Issuer      tParte
	Recipient   tParte
	Currency    string  `json:",omitempty"`
	Base        float64 `json:",omitempty"`
	VatRate     float64 `json:",omitempty"`
	Vat         float64 `json:",omitempty"`
	Total       float64 `json:",omitempty"`
	Stay        tStay
	Lines       []tLinea
	Transaction tTransaction
	Additional  map[string]string
}

var factura = tFactura{
	ID:        "A000000123",
	IssueDate: time.Date(2022, 04, 02, 12, 45, 21, 0, time.UTC),
	Base:      12.45,
	VatRate:   21.00,
	Vat:       12.45 * 0.21,
	Total:     12.45 * 1.21,
	Currency:  "EUR",
	Issuer: tParte{
		VatID:   "A78022555",
		Name:    "Garajes Pérez S.L.",
		Address: "Calle del Pez, 4, 1º Izda\n28220 - Majadahonda",
	},
	Recipient: tParte{
		VatID:   "32774555K",
		Name:    "Ernesto Pérez",
		Address: "Calle Los Lirios, 44\n28220 - Getafe",
	},
	Stay: tStay{
		EntryDate: time.Date(2022, 04, 02, 12, 45, 21, 0, time.UTC),
		Duration:  "148 min",
		Plate:     "1708FVS",
		Ticket:    "1234-ABCD",
	},
	Lines: []tLinea{
		{Service: "Estancia 23/08/2022 12:43 (37 min)", Price: 3.00},
		{Service: "Lavado integral vehiculo", Price: 5.96},
	},
	Transaction: tTransaction{
		Date:   time.Date(2022, 04, 02, 15, 03, 29, 0, time.UTC),
		Method: "VISA 1234",
	},
	Additional: map[string]string{"CUFE": "dd83bc58ab454dd7b90dce6fe61da574e20d24b79aa340b38fc1d25567fc69baf50b5e366528413ea67fbc15599ec0e5"},
}

func crc(t *testing.T, fn string, start, end string) uint32 {
	data, err := os.ReadFile(fn)
	if err != nil {
		t.Fatal(err.Error())
	}
	s := 0
	e := len(data)
	if start > "" {
		k := bytes.Index(data, []byte(start))
		if k > 0 {
			s = k
		}
	}
	if end > "" {
		k := bytes.Index(data, []byte(end))
		if k > 0 {
			e = k + len(end)
		}
	}
	return crc32.ChecksumIEEE(data[s:e])
}

// Requiere Xpdf command line tools: https://www.xpdfreader.com/download.html
func readPdfText(t *testing.T, file string) string {
	cmd := exec.Command("pdftotext", "-raw", "-enc", "UTF-8", file)
	err := cmd.Run()
	assert.NoError(t, err)
	file = strings.TrimSuffix(file, ".pdf") + ".txt"
	content, err := os.ReadFile(file)
	assert.NoError(t, err)
	return string(content)
}
