// Funciones de utilidad para consumir servicios REST
package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func DoRestGet[T any](host, endpoint string, params url.Values, headers ...string) (response T, code int, err error) {
	fullURL, err := url.JoinPath(host, endpoint)
	if err != nil {
		return
	}
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	for _, h := range headers {
		k, v, ok := strings.Cut(h, ":")
		if ok {
			req.Header.Add(k, v)
		}
	}
	client := &http.Client{}
	httpResponse, err := client.Do(req)
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

func DoRestPost[T any](host, endpoint string, request any, headers ...string) (response T, code int, err error) {
	fullURL, err := url.JoinPath(host, endpoint)
	if err != nil {
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewReader(requestBytes))
	for _, h := range headers {
		k, v, ok := strings.Cut(h, ":")
		if ok {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpResponse, err := client.Do(req)
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

func DoRestPut[T any](host, endpoint string, request any, headers ...string) (response T, code int, err error) {
	fullURL, err := url.JoinPath(host, endpoint)
	if err != nil {
		return
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewReader(requestBytes))
	for _, h := range headers {
		k, v, ok := strings.Cut(h, ":")
		if ok {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpResponse, err := client.Do(req)
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
