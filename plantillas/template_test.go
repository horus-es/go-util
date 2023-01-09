package plantillas

import (
	"bytes"
	"hash/crc32"
	"horus-es/go-util/errores"
	"horus-es/go-util/parse"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/davrux/go-smtptester"
	"github.com/stretchr/testify/assert"
)

type tParte struct {
	VatID   string
	Name    string
	Address string
}
type tLinea struct {
	Service string
	Price   float64
}
type tFactura struct {
	ID        string
	IssueDate time.Time
	Subtotal  float64
	VatRate   float64
	Vat       float64
	Total     float64
	Currency  string
	Issuer    tParte
	Recipient tParte
	Lines     []tLinea
}

var factura = tFactura{
	ID:        "A000000123",
	IssueDate: time.Date(2022, 04, 02, 12, 45, 21, 0, time.UTC),
	Subtotal:  12.45,
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
		Address: "Calle Los Lirios, 44\n28220 - Getafe",
	},
	Lines: []tLinea{
		{Service: "Estancia 23/08/2022 12:43 (37 min)", Price: 3.00},
		{Service: "Lavado integral vehiculo", Price: 5.96},
	},
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

func TestMergeXhtmlTemplate(t *testing.T) {
	p, err := os.ReadFile("template_test.html")
	assert.Nil(t, err)
	f, err := MergeXhtmlTemplate("template_test.html", string(p), factura, "", parse.DMA, parse.EUR)
	assert.Nil(t, err)
	os.WriteFile("template_test_out.html", []byte(f), 0666)
	crc1 := crc(t, "template_test_expect.html", "", "")
	crc2 := crc(t, "template_test_out.html", "", "")
	assert.Equal(t, crc1, crc2)
}

func ExampleMergeXhtmlTemplate() {
	// Cargar plantilla
	plantilla, err := os.ReadFile("plantilla.html")
	errores.PanicIfError(err)
	// Fusionar plantilla con estructura factura
	f, err := MergeXhtmlTemplate(
		"html",
		string(plantilla),
		factura,
		"/assets",
		parse.DMA,
		parse.EUR,
	)
	errores.PanicIfError(err)
	// Guardar salida
	os.WriteFile("pagina.html", []byte(f), 0666)
}

func TestGenerateXhtmlPdf(t *testing.T) {
	plantilla, err := os.ReadFile("template_test.html")
	assert.Nil(t, err)
	wd, err := os.Getwd()
	assert.Nil(t, err)
	wd = "file:///" + strings.ReplaceAll(wd, "\\", "/")
	err = GenerateXhtmlPdf("template", string(plantilla), factura, wd, parse.DMA, parse.EUR, "template_test_out.pdf", "--no-outline")
	assert.Nil(t, err)
	crc1 := crc(t, "template_test_expect.pdf", "\n>>\n", "")
	crc2 := crc(t, "template_test_out.pdf", "\n>>\n", "")
	assert.Equal(t, crc1, crc2)
}

func ExampleGenerateXhtmlPdf() {
	// Carga plantilla HTML
	plantilla, err := os.ReadFile("plantilla.html")
	errores.PanicIfError(err)
	// Genera fichero PDF
	err = GenerateXhtmlPdf(
		"pdf",
		string(plantilla),
		factura,
		"file:///assets",
		parse.DMA,
		parse.EUR,
		"fichero.pdf",
		"--no-outline",
	)
	errores.PanicIfError(err)
}

func TestSendXhtmlMail(t *testing.T) {
	// servidor SMTP dummy
	s := smtptester.Standard()
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Printf("smtp server response %s", err)
		}
	}()
	defer s.Close()
	time.Sleep(time.Second * 10) // Tiempo para deshabilitar el fw de windows
	// cargamos plantilla
	plantilla, err := os.ReadFile("template_test.html")
	assert.Nil(t, err)
	// enviamos correo
	from := "automaticos@horus.es"
	to := "pablo.leon@horus.es"
	bcc := "pablo.leon100@gmail.com"
	err = SendXhtmlMail("template_test.html", string(plantilla), factura, "https://spark2.horus.es/assets", parse.DMA, parse.EUR, []string{"template_test_expect.pdf"},
		from, to, "Prueba de correo", []string{bcc}, []string{bcc},
		"localhost", 2525, "smtpuser", "c2VjcmV0bw==")
	assert.Nil(t, err)
	// comparamos
	eml, ok := smtptester.GetBackend(s).Load(from, []string{to, bcc})
	if !ok {
		t.Fatal("Correo no transmitido")
	}
	os.WriteFile("template_test_out.eml", eml.Data, 0666)
	crc1 := crc(t, "template_test_expect.eml", "<!DOCTYPE ", "</html>")
	crc2 := crc(t, "template_test_out.eml", "<!DOCTYPE ", "</html>")
	assert.Equal(t, crc1, crc2)
}

func ExampleSendXhtmlMail() {
	// Carga plantilla HTML
	plantilla, err := os.ReadFile("plantilla.html")
	errores.PanicIfError(err)
	// Enviar por correo
	err = SendXhtmlMail(
		"mail",
		string(plantilla),
		factura,
		"https://spark2.horus.es/assets",
		parse.DMA,
		parse.EUR,
		[]string{"adjunto.pdf"},
		"remitente@horus.es",
		"destinatario@horus.es",
		"Asunto",
		[]string{"bcc@horus.es"},
		[]string{"replyto@horus.es"},
		"smtp.horus.es",
		25,
		"automaticos@horus.es",
		"secreto",
	)
	errores.PanicIfError(err)
}
