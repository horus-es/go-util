package ginhelper_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/v2/ginhelper"
)

func Example() {
	// Modo producción
	gin.SetMode(gin.ReleaseMode)
	// Crea el router
	router := gin.New()
	// Evitamos los redirect si falta la barra final
	router.RedirectTrailingSlash = false
	// Middleware recuperación errores
	router.Use(ginhelper.MiddlewarePanic())
	// Middleware CORS
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsCfg))
	// Middleware no implementado y debugLogger
	router.Use(ginhelper.MiddlewareNotImplemented(), ginhelper.MiddlewareLogger(true))

	// Rutas
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/json", func(c *gin.Context) {
		c.String(201, `{"response":"Sorry, the market was closed ..."}`)
	})
	router.POST("/multipart", func(c *gin.Context) {
		c.PureJSON(201, gin.H{})
	})

	// Solicitud Ping/Pong
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Solicitud JSON
	req, _ = http.NewRequest("POST", "/json", bytes.NewBufferString(`{"request":"Buy cheese and bread for breakfast."}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Solicitud multipart/form-data
	body := `-----------------------------9051914041544843365972754266
Content-Disposition: form-data; name="text"

text default
-----------------------------9051914041544843365972754266
Content-Disposition: form-data; name="file1"; filename="a.txt"
Content-Type: text/plain

Content of a.txt.

-----------------------------9051914041544843365972754266
Content-Disposition: form-data; name="file2"; filename="a.html"
Content-Type: text/html

<!DOCTYPE html><title>Content of a.html.</title>

-----------------------------9051914041544843365972754266--`

	req, _ = http.NewRequest("POST", "/multipart", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------9051914041544843365972754266")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Output:
	// INFO: GET /ping
	// INFO: HTTP 200 OK - 0ms
	// INFO: pong
	// INFO: ==================================================
	// INFO: POST /json
	// INFO: {"request":"Buy cheese and bread for breakfast."}
	// INFO: HTTP 201 Created - 0ms
	// INFO: {"response":"Sorry, the market was closed ..."}
	// INFO: ==================================================
	// INFO: POST /multipart
	// INFO: Content-Type: multipart/form-data
	// INFO: Content-Length: 538
	// INFO: HTTP 201 Created - 0ms
	// INFO: {}
	// INFO: ==================================================

}
