// Middlewares para GIN-GONIC
package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/errores"
	"github.com/horus-es/go-util/postgres"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var mwLog *errores.Logger

// Establece el logger. Si logger es nil, todos los mensajes se muestran en la consola.
func InitMidleware(logger *errores.Logger) {
	mwLog = logger
}

// Middleware para inicializar estado a no implementado
func MiddlewareNotImplemented() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNotImplemented)
		c.Next()
	}
}

// Middleware de registro de actividad en modo produccion
func productionLogger(c *gin.Context) {
	t := time.Now()
	c.Next()
	latency := time.Since(t).Milliseconds()
	statusCode := c.Writer.Status()
	statusText := http.StatusText(statusCode)
	method := c.Request.Method
	path := c.Request.URL.String()
	path, _ = url.PathUnescape(path)
	path = strings.ReplaceAll(path, " ", "+")
	const fmt = "%d %s - %s %s - %dms"
	const max = 1000
	if statusCode >= 200 && statusCode <= 299 {
		if latency < max {
			// Todo OK
			mwLog.Infof(fmt, statusCode, statusText, method, path, latency)
		} else {
			// Respuesta lenta
			mwLog.Warnf(fmt, statusCode, statusText, method, path, latency)
		}
		return
	}
	if statusCode >= 300 && statusCode <= 499 {
		// Error de solicitud o redirección
		mwLog.Warnf(fmt, statusCode, statusText, method, path, latency)
	} else {
		// Error de servidor
		mwLog.Errorf(fmt, statusCode, statusText, method, path, latency)
	}
}

// Middleware de registro de actividad en modo depuración

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func debugLogger(c *gin.Context) {
	// Registramos metodo y ruta
	method := c.Request.Method
	path := c.Request.URL.String()
	path, _ = url.PathUnescape(path)
	path = strings.ReplaceAll(path, " ", "+")
	mwLog.Infof("%s %s", method, path)
	if c.Request.ContentLength != 0 {
		// Registramos solicitud duplicando reader
		buffer, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buffer))
		if len(buffer) < 3000 {
			mwLog.Infof(string(buffer))
		} else {
			mwLog.Infof(string(buffer[:3000]) + "···")
		}
	}
	// Duplicamos writer
	blw := &bodyLogWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
	c.Writer = blw
	// Siguiente en cadena
	t := time.Now()
	c.Next()
	latency := time.Since(t).Milliseconds()
	// Registramos status
	statusCode := c.Writer.Status()
	statusText := http.StatusText(statusCode)
	mwLog.Infof("HTTP %d %s - %dms", statusCode, statusText, latency)
	// Registramos respuesta
	ct := c.Writer.Header().Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") && blw.body.Len() > 0 {
		mwLog.Infof(blw.body.String())
	} else {
		if ct != "" {
			mwLog.Infof("Content-Type: %s", ct)
		}
		mwLog.Infof("Content-Length: %d", blw.body.Len())
	}
	mwLog.Infof("==================================================")
}

// Devuelve el logger de depuración o de producción
func MiddlewareLogger(debug bool) gin.HandlerFunc {
	if debug {
		return debugLogger
	} else {
		return productionLogger
	}
}

// Middleware de gestión de transacciones
func MiddlewareTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := postgres.StartTX()
		defer postgres.RollbackTX(tx)
		c.Next()
		statusCode := c.Writer.Status()
		if statusCode >= 200 && statusCode <= 299 {
			postgres.CommitTX(tx)
		}
	}
}

// Devuelve el handler de ruta inxistente (a.k.a. 404)
func HandleNoRoute(debug bool, proxy string) gin.HandlerFunc {
	if !debug {
		return func(c *gin.Context) {
			// Producción
			mwLog.Infof("No encontrado: %s %s", c.Request.Method, c.Request.URL.Path)
			c.PureJSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
		}
	}
	if proxy == "" {
		return func(c *gin.Context) {
			// Depuración sin proxy
			mwLog.Infof("No implementado: %s %s", c.Request.Method, c.Request.URL.Path)
			c.PureJSON(http.StatusNotImplemented, gin.H{"error": "No implementado"})
		}
	}
	return func(c *gin.Context) {
		// Depuración con proxy
		mwLog.Infof("Redirección: %s %s%s", c.Request.Method, proxy, c.Request.URL.Path)
		director := func(req *http.Request) {
			proxyUrl, _ := url.Parse(proxy)
			req.URL.Scheme = proxyUrl.Scheme
			req.URL.Host = proxyUrl.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Middleware de recuperación de errores
func MiddlewarePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recuperaDiferido(c)
		c.Next()
	}
}

// Auxiliar de MiddlewarePanic
func recuperaDiferido(c *gin.Context) {
	causa := recover()
	if causa == nil {
		// No panic
		return
	}
	// ¿Es debido a un error de red?
	networkError := isNetworkError(causa.(error))
	if networkError {
		// Si hay error de red, no podemos responder nada ...
		mwLog.Errorf("Error de red: %v", causa)
	} else {
		// TODO: ¿añadir 404 para claves no halladas?
		// ¿Es debido a una violación de restricción SQL o custom?
		errorSQL, msg := postgres.GetErrorSQL(causa.(error))
		switch errorSQL {
		case postgres.INTEGRITY_CONSTRAINT_VIOLATION:
			// Si es una violación de restriccion SQL, se supone que la culpa es del cliente
			// TODO: ¿Cambiar mensaje según tipo de violación?
			c.PureJSON(http.StatusBadRequest, mwLog.BadHttpRequest("Valor duplicado", causa))
		case postgres.PL_PGSQL_RAISE_EXCEPTION:
			// Si es una excepcion levantada en un procedimiento, se supone que la culpa es del cliente
			c.PureJSON(http.StatusBadRequest, mwLog.BadHttpRequest(msg, causa))
		default:
			// Otros errores seguramente sean de programación o de sistema
			mwLog.Errorf("Error interno: %v", causa)
			c.PureJSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		}
	}
	c.Abort()
}

// Determina si un error es debido a un corte de red (broken pipe / connection reset by peer)
func isNetworkError(err error) bool {
	return errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
}

// Gestiona la documentación swagger
func HandleSwagger(router *gin.Engine, ruta string, mapaSVG map[string]string) {
	svgHandler := func(c *gin.Context) {
		file := c.Param("any")
		if !strings.HasSuffix(file, ".svg") {
			return
		}
		file = strings.TrimPrefix(file, "/")
		c.Writer.Header().Add("Content-Type", "image/svg+xml")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.WriteString(mapaSVG[file])
		c.Abort()
	}
	ruta = path.Join(ruta, "swagger")
	router.GET(path.Join(ruta, "*any"), svgHandler, ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true), ginSwagger.DefaultModelsExpandDepth(-1)))
}
