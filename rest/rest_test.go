package rest

import (
	"fmt"
	"net/url"

	"github.com/horus-es/go-util/v2/errores"
)

func ExampleDoRestGet() {
	// Servicio REST publico que devuelve los festivos en una determinada región
	host := "https://kayaposoft.com"
	endpoint := "/enrico/json/v2.0"
	params := url.Values{}
	params.Set("action", "getHolidaysForYear")
	params.Set("year", "2023")
	params.Set("country", "esp")
	type tFestivos []struct {
		Date struct {
			Day       int
			Month     int
			Year      int
			DayofWeek int
		}
		Name []struct {
			Lang string
			Text string
		}
		HolidayType string
	}
	festivos, code, err := DoRestGet[tFestivos](host, endpoint, params)
	errores.PanicIfError(err, "Error de red")
	errores.PanicIfTrue(code < 200 || code > 299, "Recibido código de respuesta http: %d", code)
final:
	for _, festivo := range festivos {
		for _, name := range festivo.Name {
			if name.Lang == "es" {
				fmt.Printf("El primer festivo del año %d en España es %s (%02d/%02d)", festivo.Date.Year, name.Text, festivo.Date.Day, festivo.Date.Month)
				break final
			}
		}
	}
	// Output: El primer festivo del año 2023 en España es Año Nuevo (01/01)
}

func ExampleDoRestPost() {
	// Servicio REST publico de prueba que devuelve lo mismo que se le envía
	host := "https://httpbin.org"
	endpoint := "anything"
	type tEstructura struct {
		Uno  string
		Dos  string
		Tres string
	}
	request := tEstructura{Uno: "1", Dos: "2", Tres: "3"}
	response, code, err := DoRestPost[map[string]any](host, endpoint, request)
	errores.PanicIfError(err, "Error de red")
	errores.PanicIfTrue(code < 200 || code > 299, "Recibido código de respuesta http: %d", code)
	fmt.Println(response["data"])
	// Output: {"Uno":"1","Dos":"2","Tres":"3"}
}

func ExampleDoRestPut() {
	// Servicio REST publico de prueba que devuelve lo mismo que se le envía
	host := "https://httpbin.org"
	endpoint := "anything"
	type tEstructura struct {
		Uno  string
		Dos  string
		Tres string
	}
	request := tEstructura{Uno: "1", Dos: "2", Tres: "3"}
	response, code, err := DoRestPut[map[string]any](host, endpoint, request)
	errores.PanicIfError(err, "Error de red")
	errores.PanicIfTrue(code < 200 || code > 299, "Recibido código de respuesta http: %d", code)
	fmt.Println(response["data"])
	// Output: {"Uno":"1","Dos":"2","Tres":"3"}
}
