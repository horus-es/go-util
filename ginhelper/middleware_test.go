package ginhelper_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sync"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/v3/formato"
	"github.com/horus-es/go-util/v3/ginhelper"
	"github.com/horus-es/go-util/v3/logger"
	"github.com/horus-es/go-util/v3/postgres"
	"github.com/jackc/pgx/v5/pgtype"
)

func Example() {
	// Borramos el fichero de log para el ejemplo
	const logfile = "testlog"
	os.Remove(logfile + ".log")
	// Creamos el logger en modo debug
	logger.InitLogger(logfile, true)
	// Modo producción
	gin.SetMode(gin.ReleaseMode)
	// Creamos el router
	router := gin.New()
	// Evitamos los redirect si falta la barra final
	router.RedirectTrailingSlash = false
	// Middleware para recuperación errores
	router.Use(ginhelper.MiddlewarePanic())
	// Middleware CORS
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsCfg))
	// Middlewares para no implementado y logger
	router.Use(ginhelper.MiddlewareNotImplemented(), ginhelper.MiddlewareLogger(1000, ``))
	// Rutas
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/json", func(c *gin.Context) {
		logger.Warnf(c, "Las advertencias siempre salen por STDOUT")
		logger.Errorf(c, "Los errores siempre salen por STDERR")
		c.String(201, `{"response":"Sorry, the market was closed ..."}`)
	})
	router.POST("/multipart", func(c *gin.Context) {
		c.PureJSON(201, gin.H{})
	})
	router.GET("/panic", func(c *gin.Context) {
		panic("panico inesperado")
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
	body, _ := os.ReadFile("multipart.mime")
	req, _ = http.NewRequest("POST", "/multipart", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------9051914041544843365972754266")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Solicitud Panic
	req, _ = http.NewRequest("GET", "/panic", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	logger.CloseLogger()

	// Imprime el fichero de log quitando las fechas y el stack trace
	printLogFile(logfile)

	// Nota: los warnings siempre salen por STDOUT y los errores seimpre salen por STDERR, ademas de quedar registrados en el fichero de log.

	// Output:
	// WARN: Las advertencias siempre salen por STDOUT
	// INFO: GET /ping
	// INFO: HTTP 200 OK - 0ms
	// INFO: pong
	// ==================================================
	// INFO: POST /json
	// INFO: {"request":"Buy cheese and bread for breakfast."}
	// WARN: Las advertencias siempre salen por STDOUT
	// ERROR: Los errores siempre salen por STDERR
	// INFO: HTTP 201 Created - 0ms
	// INFO: {"response":"Sorry, the market was closed ..."}
	// ==================================================
	// INFO: POST /multipart
	// INFO: Content-Type: multipart/form-data
	// INFO: Content-Length: 538
	// INFO: HTTP 201 Created - 0ms
	// INFO: {}
	// ==================================================
	// INFO: GET /panic
	// ERROR: panic: panico inesperado
	// goroutine 1 [running]:
	// runtime/debug.Stack()
	// ... stack trace ...
}

// Imprime el fichero de log quitando las fechas y el stack trace
func printLogFile(logfile string) {
	data, _ := os.ReadFile(logfile + ".log")
	stackStart := []byte("runtime/debug.Stack()")
	if idx := bytes.Index(data, stackStart); idx != -1 {
		data = data[:idx+len(stackStart)]
	}
	data = regexp.MustCompile(`(?m)^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\s*`).ReplaceAll(data, nil)
	fmt.Println(string(data))
	fmt.Println("... stack trace ...")
}

func TestGin(t *testing.T) {

	log := logger.NewLogger("testlog", true)
	ginhelper.InitGinHelper(log)
	postgres.InitPool(`host=devel.horus.es port=43210 user=SPARK2 password=lahh4jaequ2I dbname=SPARK2 sslmode=disable application_name=_TEST_`, log)

	// Modo producción
	gin.SetMode(gin.ReleaseMode)

	// Crea el router
	router := gin.New()
	router.Use(ginhelper.MiddlewareLogger(100, `GET /gin_`), ginhelper.MiddlewareTransaction())

	// Rutas
	router.GET("/gin_test", func(c *gin.Context) {

		type item struct {
			Id         pgtype.UUID
			Nombre     string
			Fechas     formato.Fecha
			FhAvisos   pgtype.Timestamp
			NumAvisos  int
			FhErrores  pgtype.Timestamp
			NumErrores int
			Tarifas    []byte
		}
		var lista []item
		postgres.GetOrderedRows(c, &lista, `SELECT p.id, p.nombre::text, o.fechas::text, p.tarifas, (SELECT MAX(p2.desde) FROM problemas p2 LEFT JOIN mensajes m ON m.codigo = p2.codigo WHERE p2.hasta IS NULL AND m.nivel = 'WARN' AND p2.parking = p.id) AS fh_avisos, (SELECT COUNT(*) FROM problemas p2 LEFT JOIN mensajes m ON m.codigo = p2.codigo WHERE p2.hasta IS NULL AND m.nivel = 'WARN' AND p2.parking = p.id) AS num_avisos, (SELECT MAX(p2.desde) FROM problemas p2 LEFT JOIN mensajes m ON m.codigo = p2.codigo WHERE p2.hasta IS NULL AND m.nivel = 'ERROR' AND p2.parking = p.id) AS fh_errores, (SELECT COUNT(*) FROM problemas p2 LEFT JOIN mensajes m ON m.codigo = p2.codigo WHERE p2.hasta IS NULL AND m.nivel = 'ERROR' AND p2.parking = p.id) AS num_errores FROM operadores o JOIN parkings p ON o.id = p.operador JOIN personal pe ON o.id = pe.operador JOIN sesiones s ON pe.id = s.empleado WHERE s.id = 'd6d8770d-5619-4a9d-9d10-95e508a35b71' ORDER BY 2`)
	})

	// Solicitudes
	var wg sync.WaitGroup
	for range 200 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/gin_test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}()
	}
	wg.Wait()
	log.CloseLogger()
}
