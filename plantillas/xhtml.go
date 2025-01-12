// Funciones para procesar plantillas y convertirlas en HTML, PDF o EMAIL
/*
Las plantillas soportan la sintaxis estándar de GO y además los siguientes attributos tomados de ThymeLeaf:
  * th:if="condition" => {{ if condition }} T {{ end }}
  * th:each="items"   => {{ range items }} T {{ end }}
  * th:with="item"    => {{ with item }} T {{end}}
  * th:remove="all"   => elimina el tag
  * th:text="content" => reemplaza el contenido del tag
  * th:attr="value"   => reemplaza el valor del atributo

Se soportan las funciones de formato DATETIME, DATE, TIME, PRICE y BR.

Ejemplo de plantillla en https://github.com/horus-es/go-util/blob/main/plantillas/plantilla.html

Ejemplos de plantillas preparadas para mailing en https://postmarkapp.com/transactional-email-templates

Ejemplos de uso de los atributos de Thymeleaf en https://www.thymeleaf.org/doc/tutorials/2.1/usingthymeleaf.html
*/
package plantillas

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/pkg/errors"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
)

// Fusiona una plantilla XHTML con un struct o map de datos.
//   - name: nombre arbitrario para la plantilla que aparece en los mensajes de error
//   - xhtml: plantilla en formato XHTML
//   - datos: estructura de datos para fusionar con la plantilla
//   - assets: ruta de imágenes u otros recursos (attributos src y href)
//   - ff: formato de las fechas para las funciones DATETIME y DATE
//   - fp: formato de los precios para la funcion PRICE
func MergeXhtmlTemplate(name, xhtml string, datos any, assets string, ff formato.Fecha, fp formato.Moneda) (string, error) {
	gotmpl, err := thTemplate(name, xhtml, assets)
	if err != nil {
		return "", err
	}
	var funciones = template.FuncMap{
		"DATETIME": func(x any) string {
			switch t := x.(type) {
			case time.Time:
				return formato.PrintFechaHora(t, ff)
			case string:
				t2, err := formato.ParseFechaHora(t, ff)
				if err == nil {
					return formato.PrintFechaHora(t2, ff)
				}
			}
			errores.PanicIfTrue(true, "fecha %q no soportada", x)
			return ""
		},
		"DATE": func(x any) string {
			switch t := x.(type) {
			case time.Time:
				return formato.PrintFecha(t, ff)
			case string:
				t2, err := formato.ParseFechaHora(t, ff)
				if err == nil {
					return formato.PrintFecha(t2, ff)
				}
			}
			errores.PanicIfTrue(true, "fecha %q no soportada", x)
			return ""
		},
		"TIME": func(x any) string {
			switch t := x.(type) {
			case time.Time:
				return formato.PrintHora(t, false)
			case string:
				t2, err := formato.ParseFechaHora(t, ff)
				if err == nil {
					return formato.PrintHora(t2, false)
				}
			}
			errores.PanicIfTrue(true, "fecha %q no soportada", x)
			return ""
		},
		"PRICE": func(f float64) string {
			return formato.PrintPrecio(f, fp)
		},
		"BR": func(s string) template.HTML {
			// Cambia los saltos de línea por <br/>
			lineas := strings.Split(s, "\n")
			for k := range lineas {
				lineas[k] = template.HTMLEscapeString(lineas[k])
			}
			return template.HTML(strings.Join(lineas, "<br/>"))
		},
	}
	var opt string
	if reflect.TypeOf(datos).Kind() == reflect.Map {
		// En los mapas se permite que falten campos
		opt = "missingkey=zero"
	} else {
		// En las estructuras se exige la existencia del dato
		opt = "missingkey=error"
	}
	tmpl, err := template.New(name).Funcs(funciones).Option(opt).Parse(gotmpl)
	if err != nil {
		return "", err
	}
	var marshaled bytes.Buffer
	err = tmpl.Execute(&marshaled, datos)
	if err != nil {
		return "", err
	}
	return marshaled.String(), nil
}

// Traduce una plantilla de estilo thymeleaf al estilo GO.
// Tambien añade assets a las URLs de los atributos src o href
func thTemplate(name, template, assets string) (string, error) {
	doc := etree.NewDocument()
	doc.ReadSettings.Entity = map[string]string{"nbsp": "\u00A0"}
	err := doc.ReadFromString(template)
	if err != nil {
		return "", errors.Wrap(err, name)
	}
	var base *url.URL
	if assets > "" {
		base, err = url.Parse(assets)
		if err != nil {
			return "", errors.Wrap(err, name)
		}
	}
	root := doc.Root()
	//	for _, tag := range root.ChildElements() {
	for _, tag := range root.FindElements("//") {
		parent := tag.Parent()
		if parent == nil {
			// tag eliminado
			continue
		}
		th_remove := selectOneAttr(tag, "th:remove")
		if th_remove != nil {
			// de momento solo se soporta th:remove="all"
			if th_remove.Value != "all" {
				return "", fmt.Errorf("%s: th:remove=%q no soportado", name, th_remove.Value)
			}
			parent.RemoveChildAt(tag.Index())
			continue
		}
		procesaTag(tag, "th:if", "if")
		procesaTag(tag, "th:each", "range")
		procesaTag(tag, "th:with", "with")
		th_text := selectOneAttr(tag, "th:text")
		if th_text != nil {
			for _, c := range tag.FindElements("*") {
				tag.RemoveChild(c)
			}
			tag.SetText(th_text.Value)
			tag.RemoveAttr("th:text")
		}
		th_otros := selectSpaceAttr(tag, "th")
		for _, otro := range th_otros {
			tag.CreateAttr(otro.Key, otro.Value)
			tag.RemoveAttr("th:" + otro.Key)
		}
		err = procesaURL(tag, base, "src")
		if err != nil {
			return "", errors.Wrap(err, name)
		}
		err = procesaURL(tag, base, "href")
		if err != nil {
			return "", errors.Wrap(err, name)
		}
	}
	doc.WriteSettings.CanonicalText = true
	doc.WriteSettings.CanonicalAttrVal = true
	return doc.WriteToString()
}

// Funcion axiliar que cambia el attributo th:attr="item" por {{accion item}}T{{end}}
func procesaTag(tag *etree.Element, attribute, action string) {
	th := selectOneAttr(tag, attribute)
	if th == nil {
		return
	}
	parent := tag.Parent()
	parent.InsertChildAt(tag.Index(), etree.NewText("{{"+action+" "+th.Value+"}}"))
	parent.InsertChildAt(tag.Index()+1, etree.NewText("{{end}}"))
	tag.RemoveAttr(attribute)
}

// Funcion axiliar que añade base a los attributos href y/o src del tag
func procesaURL(tag *etree.Element, base *url.URL, attribute string) error {
	if base == nil {
		return nil
	}
	ref := selectOneAttr(tag, attribute)
	if ref == nil {
		return nil
	}
	dest := *base
	dest.Path = path.Join(dest.Path, ref.Value)
	ref.Value = dest.String()
	return nil
}

// Funcion similar a etree.SelectAttr, pero produce pánico si existe mas de 1 atributo
func selectOneAttr(tag *etree.Element, key string) *etree.Attr {
	space := ""
	l, r, ok := strings.Cut(key, ":")
	if ok {
		space = l
		key = r
	}
	var z, n int
	for i, a := range tag.Attr {
		if a.Space == space && a.Key == key {
			z = i
			n++
		}
	}
	if n == 1 {
		return &tag.Attr[z]
	}
	errores.PanicIfTrue(n > 1, "Atributo %q duplicado", key)
	return nil
}

// Obtiene todos los atributos de un determinado espacio de nombres
func selectSpaceAttr(tag *etree.Element, space string) []*etree.Attr {
	result := []*etree.Attr{}
	for _, a := range tag.Attr {
		if a.Space == space {
			result = append(result, &a)
		}
	}
	return result
}

// Envia un correo a partir de una plantilla XHTML. Parámetros:
//   - name: nombre arbitrario para la plantilla que aparece en los mensajes de error
//   - xhtml: plantilla en formato XHTML
//   - datos: estructura de datos para fusionar con la plantilla
//   - assets: URL de imágenes u otros recursos (attributos src y href). Debe ser una ruta públicamente accesible por internet.
//   - ff: formato de las fechas para las funciones DATETIME y DATE
//   - fp: formato de los precios para la funcion PRICE
//   - adjuntos: ficheros a adjuntar
//   - to,form,subject,bcc,replyto: parámetros MIME
//   - host,port,username,password: parámtros SMTP. La contraseña debe ir codificada en base64.
func SendXhtmlMail(name, xhtml string, datos any, assets string, ff formato.Fecha, fp formato.Moneda, adjuntos []string,
	from, to, subject string, bcc, replyto []string,
	host string, port int, username, password string) error {

	// procesa la plantilla XHTML
	body, err := MergeXhtmlTemplate(name, xhtml, datos, assets, ff, fp)
	if err != nil {
		return errors.Wrap(err, name)
	}

	// css-inline: mejora la compatibilidad de los clientes de email
	bodycss, err := premailer.NewPremailerFromBytes([]byte(body), premailer.NewOptions())
	if err != nil {
		return errors.Wrap(err, name)
	}
	html, err := bodycss.Transform()
	if err != nil {
		return errors.Wrap(err, name)
	}

	// cabeceras
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	if len(bcc) > 0 {
		m.SetHeader("Bcc", bcc...)
	}
	if len(replyto) > 0 {
		m.SetHeader("Reply-To", replyto...)
	}

	// cuerpo XHTML
	m.SetBody("text/html", html)

	// ficheros adjuntos
	for _, f := range adjuntos {
		_, err = os.Stat(f)
		errores.PanicIfError(err)
		m.Attach(f)
	}

	// parámetros SMTP
	if len(password) > 0 {
		p, err := base64.StdEncoding.DecodeString(password)
		if err != nil {
			return errors.Wrap(err, name)
		}
		password = string(p)
	}
	d := gomail.NewDialer(host, port, username, password)
	err = d.DialAndSend(m)
	if err != nil {
		return errors.Wrap(err, name)
	}

	// OK
	return nil
}

// Genera un fichero PDF a partir de una plantilla XHTML
// usando la utilidad wkhtmltopdf, que debe estar previamente instalada. Parámetros:
//   - name: nombre arbitrario para la plantilla que aparece en los mensajes de error.
//   - xhtml: plantilla en formato XHTML
//   - datos: estructura de datos para fusionar con la plantilla
//   - assets: ruta de imágenes u otros recursos (attributos src y href). Si es una ruta local, debe estar precedida por file://
//   - ff: formato de las fechas para las funciones DATETIME y DATE
//   - fp: formato de los precios para la funcion PRICE
//   - out: fichero PDF de salida
//   - opciones: opciones adicionales utilidad wkhtmltopdf (ver https://wkhtmltopdf.org/usage/wkhtmltopdf.txt)
func GenerateXhtmlPdf(name, xhtml string, datos any, assets string, ff formato.Fecha, fp formato.Moneda, out string, opciones ...string) error {
	// Procesa la plantilla XHTML
	body, err := MergeXhtmlTemplate(name, xhtml, datos, assets, ff, fp)
	if err != nil {
		return err
	}
	// Fichero temporal
	tmp, err := os.CreateTemp("", "horus-*.html")
	if err != nil {
		return errors.Wrap(err, name)
	}
	defer os.Remove(tmp.Name())
	_, err = tmp.WriteString(body)
	if err != nil {
		return errors.Wrap(err, name)
	}
	err = tmp.Close()
	if err != nil {
		return errors.Wrap(err, name)
	}
	// Ejecución wkhtmltopdf
	args := append([]string{"-q", "--enable-local-file-access"}, opciones...)
	args = append(args, tmp.Name(), out)
	cmd := exec.Command("wkhtmltopdf", args...)
	var log bytes.Buffer
	cmd.Stdout = &log
	cmd.Stderr = &log
	err = cmd.Run()
	if err != nil {
		if log.Len() > 0 {
			return fmt.Errorf("%s: %s", name, log.String())
		}
		return errors.Wrap(err, name)
	}
	return nil
}
