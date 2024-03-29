package plantillas

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/davrux/go-smtptester"
	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
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

// Requiere Xpdf command line tools: https://www.xpdfreader.com/download.html
func readPdfText(t *testing.T, file string) string {
	cmd := exec.Command("pdftotext", "-raw", "-enc", "UTF-8", file)
	err := cmd.Run()
	assert.Nil(t, err)
	file = strings.TrimSuffix(file, ".pdf") + ".txt"
	content, err := os.ReadFile(file)
	assert.Nil(t, err)
	return string(content)
}

func TestMergeXhtmlTemplate(t *testing.T) {
	p, err := os.ReadFile("template_test.html")
	assert.Nil(t, err)
	f, err := MergeXhtmlTemplate("template_test.html", string(p), factura, "", formato.DMA, formato.EUR)
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
		formato.DMA,
		formato.EUR,
	)
	errores.PanicIfError(err)
	// Guardar salida
	os.WriteFile("pagina.html", []byte(f), 0666)
	fmt.Println("Ok")
	// Output: Ok
}

func TestGenerateXhtmlPdf(t *testing.T) {
	plantilla, err := os.ReadFile("template_test.html")
	assert.Nil(t, err)
	wd, err := os.Getwd()
	assert.Nil(t, err)
	wd = "file:///" + strings.ReplaceAll(wd, "\\", "/")
	err = GenerateXhtmlPdf("template", string(plantilla), factura, wd, formato.DMA, formato.EUR, "template_test_out.pdf", "--no-outline")
	assert.Nil(t, err)
	t1 := readPdfText(t, "template_test_expect.pdf")
	t2 := readPdfText(t, "template_test_out.pdf")
	assert.Equal(t, t1, t2)
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
		formato.DMA,
		formato.EUR,
		"fichero.pdf",
		"--no-outline",
	)
	errores.PanicIfError(err)
	fmt.Println("Ok")
	// Output: Ok
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
	err = SendXhtmlMail("template_test.html", string(plantilla), factura, "https://spark2.horus.es/assets", formato.DMA, formato.EUR, []string{"template_test_expect.pdf"},
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
		formato.DMA,
		formato.EUR,
		[]string{"adjunto.pdf"},
		"remitente@horus.es",
		"destinatario@horus.es",
		"Asunto",
		[]string{"bcc@horus.es"},
		[]string{"replyto@horus.es"},
		"smtp.horus.es",
		25,
		"automaticos@horus.es",
		"c2VjcmV0bw==",
	)
	fmt.Println(err)
	// Output: mail: dial tcp: lookup smtp.horus.es: no such host
}
