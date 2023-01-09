// Funciones de utilidad para consumir servicios REST
package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func DoRestGet[T any](host, endpoint string, params url.Values) (response T, code int, err error) {
	fullURL, err := url.JoinPath(host, endpoint)
	if err != nil {
		return
	}
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}
	httpResponse, err := http.Get(fullURL)
	if err != nil {
		return
	}
	defer httpResponse.Body.Close()
	code = httpResponse.StatusCode
	if code >= 500 {
		return
	}
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}

func DoRestPost[T any](host, endpoint string, request any) (response T, code int, err error) {
	fullURL, err := url.JoinPath(host, endpoint)
	if err != nil {
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return
	}
	httpResponse, err := http.Post(fullURL, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return
	}
	defer httpResponse.Body.Close()
	code = httpResponse.StatusCode
	if code >= 500 {
		return
	}
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}
