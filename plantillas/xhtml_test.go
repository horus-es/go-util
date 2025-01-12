package plantillas

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/davrux/go-smtptester"
	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/stretchr/testify/assert"
)

func TestMergeXhtmlTemplate(t *testing.T) {
	p, err := os.ReadFile("plantilla.html")
	assert.NoError(t, err)
	f, err := MergeXhtmlTemplate("plantilla.html", string(p), factura, "", formato.DMA, formato.EUR)
	assert.NoError(t, err)
	os.WriteFile("xhtml_test_out.html", []byte(f), 0666)
	//os.WriteFile("xhtml_test_expect.html", []byte(f), 0666)
	crc1 := crc(t, "xhtml_test_expect.html", "", "")
	crc2 := crc(t, "xhtml_test_out.html", "", "")
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
	plantilla, err := os.ReadFile("plantilla.html")
	assert.NoError(t, err)
	wd, err := os.Getwd()
	assert.NoError(t, err)
	wd = "file:///" + strings.ReplaceAll(wd, "\\", "/")
	err = GenerateXhtmlPdf("template", string(plantilla), factura, wd, formato.DMA, formato.EUR, "xhtml_test_out.pdf", "--no-outline")
	assert.NoError(t, err)
	t1 := readPdfText(t, "xhtml_test_expect.pdf")
	t2 := readPdfText(t, "xhtml_test_out.pdf")
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
	plantilla, err := os.ReadFile("plantilla.html")
	assert.NoError(t, err)
	// enviamos correo
	from := "automaticos@horus.es"
	to := "pablo.leon@horus.es"
	bcc := "pablo.leon100@gmail.com"
	err = SendXhtmlMail("plantilla.html", string(plantilla), factura, "https://spark2.horus.es/assets", formato.DMA, formato.EUR, []string{"xhtml_test_expect.pdf"},
		from, to, "Prueba de correo", []string{bcc}, []string{bcc},
		"localhost", 2525, "smtpuser", "c2VjcmV0bw==")
	assert.NoError(t, err)
	// comparamos
	eml, ok := smtptester.GetBackend(s).Load(from, []string{to, bcc})
	if !ok {
		t.Fatal("Correo no transmitido")
	}
	os.WriteFile("xhtml_test_out.eml", eml.Data, 0666)
	//os.WriteFile("xhtml_test_expect.eml", eml.Data, 0666)
	crc1 := crc(t, "xhtml_test_expect.eml", "<!DOCTYPE ", "</html>")
	crc2 := crc(t, "xhtml_test_out.eml", "<!DOCTYPE ", "</html>")
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
