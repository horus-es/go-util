// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/postgres"
)

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
			ghLog.Infof(fmt, statusCode, statusText, method, path, latency)
		} else {
			// Respuesta lenta
			ghLog.Warnf(fmt, statusCode, statusText, method, path, latency)
		}
		return
	}
	if statusCode >= 300 && statusCode <= 499 {
		// Error de solicitud o redirección
		ghLog.Warnf(fmt, statusCode, statusText, method, path, latency)
	} else {
		// Error de servidor
		ghLog.Errorf(fmt, statusCode, statusText, method, path, latency)
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
	ghLog.Infof("%s %s", method, path)
	if c.Request.ContentLength != 0 {
		// Registramos solicitud duplicando reader
		buffer, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buffer))
		if len(buffer) < 3000 {
			ghLog.Infof(string(buffer))
		} else {
			ghLog.Infof(string(buffer[:3000]) + "···")
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
	ghLog.Infof("HTTP %d %s - %dms", statusCode, statusText, latency)
	// Registramos respuesta
	ct := c.Writer.Header().Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") && blw.body.Len() > 0 {
		ghLog.Infof(blw.body.String())
	} else {
		if ct != "" {
			ghLog.Infof("Content-Type: %s", ct)
		}
		ghLog.Infof("Content-Length: %d", blw.body.Len())
	}
	ghLog.Infof("==================================================")
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
		ghLog.Errorf("Error de red: %v", causa)
	} else {
		// TODO: ¿añadir 404 para claves no halladas?
		// ¿Es debido a una violación de restricción SQL o custom?
		errorSQL, msg := postgres.GetErrorSQL(causa.(error))
		switch errorSQL {
		case postgres.INTEGRITY_CONSTRAINT_VIOLATION:
			// Si es una violación de restriccion SQL, se supone que la culpa es del cliente
			// TODO: ¿Cambiar mensaje según tipo de violación?
			c.PureJSON(http.StatusBadRequest, BadRequestResponse("Valor duplicado", causa))
		case postgres.PL_PGSQL_RAISE_EXCEPTION:
			// Si es una excepcion levantada en un procedimiento, se supone que la culpa es del cliente
			c.PureJSON(http.StatusBadRequest, BadRequestResponse(msg, causa))
		default:
			// Otros errores seguramente sean de programación o de sistema
			ghLog.Errorf("Error interno: %v", causa)
			c.PureJSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		}
	}
	c.Abort()
}

// Determina si un error es debido a un corte de red (broken pipe / connection reset by peer)
func isNetworkError(err error) bool {
	return errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
}
