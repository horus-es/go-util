// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/v3/postgres"
)

// Middleware para inicializar estado a no implementado
func MiddlewareNotImplemented() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNotImplemented)
		c.Next()
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

// Middleware logger.
// warnTime: milisegundos de respuesta por encima de los cuales la respuesta sale con WARN en vez de INFO
// reSlow: expresion regular de las URLs lentas, para las cuales no se tiene en cuenta el timeout
func MiddlewareLogger(warnTime int64, reSlow string) gin.HandlerFunc {
	var reTime *regexp.Regexp
	if reSlow != "" {
		reTime = regexp.MustCompile(reSlow)
	}
	return func(c *gin.Context) {
		path := c.Request.URL.String()
		path, _ = url.PathUnescape(path)
		path = strings.ReplaceAll(path, " ", "+")
		path = c.Request.Method + " " + path
		ghLog.Infof(c, path)
		if c.Request.ContentLength > 0 {
			// Duplicamos reader
			buffer, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(buffer))
			// Registramos solicitud
			ctype, _, _ := strings.Cut(c.Request.Header.Get("Content-Type"), ";")
			isTxt := ctype == "application/json" || ctype == "application/xml" || strings.HasPrefix(ctype, "text/")
			if isTxt {
				ghLog.Infof(c, string(buffer))
			} else {
				if len(ctype) > 0 {
					ghLog.Infof(c, "Content-Type: %s", ctype)
				}
				ghLog.Infof(c, "Content-Length: %d", c.Request.ContentLength)
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
		if (statusCode >= 200 && statusCode <= 299) || (statusCode >= 400 && statusCode <= 499) {
			if latency > warnTime && (reTime == nil || !reTime.MatchString(path)) {
				// Respuesta lenta
				ghLog.Warnf(c, "HTTP %d %s - %dms", statusCode, statusText, latency)
			} else {
				// Todo OK
				ghLog.Infof(c, "HTTP %d %s - %dms", statusCode, statusText, latency)
			}
		} else {
			// Error
			ghLog.Errorf(c, "HTTP %d %s - %dms", statusCode, statusText, latency)
		}
		// Registramos respuesta
		ctype, _, _ := strings.Cut(c.Writer.Header().Get("Content-Type"), ";")
		isTxt := ctype == "application/json" || ctype == "application/xml" || strings.HasPrefix(ctype, "text/")
		if isTxt {
			ghLog.Infof(c, blw.body.String())
		} else {
			if len(ctype) > 0 {
				ghLog.Infof(c, "Content-Type: %s", ctype)
			}
			ghLog.Infof(c, "Content-Length: %d", blw.body.Len())
		}
		ghLog.Flush(c)
	}
}

// Middleware de gestión de transacciones
func MiddlewareTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		postgres.StartTX(c)
		defer postgres.RollbackTX(c)
		c.Next()
		statusCode := c.Writer.Status()
		if statusCode >= 200 && statusCode <= 299 {
			postgres.CommitTX(c)
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
	e, ok := causa.(error)
	if ok && isNetworkError(e) {
		// Si hay error de red, no podemos responder nada ...
		ghLog.Errorf(c, "Error de red: %v", causa)
	} else {
		// TODO: ¿añadir 404 para claves no halladas?
		// ¿Es debido a una violación de restricción SQL o custom?
		errorSQL := postgres.NON_SQL
		var msg string
		if ok {
			errorSQL, msg = postgres.GetErrorSQL(e)
		}
		switch errorSQL {
		case postgres.INTEGRITY_CONSTRAINT_VIOLATION:
			// Si es una violación de restriccion SQL, se supone que la culpa es del cliente
			// TODO: ¿Cambiar mensaje según tipo de violación?
			c.PureJSON(http.StatusBadRequest, BadRequestResponse(c, "Valor duplicado", causa))
		case postgres.PL_PGSQL_RAISE_EXCEPTION:
			// Si es una excepcion levantada en un procedimiento, se supone que la culpa es del cliente
			c.PureJSON(http.StatusBadRequest, BadRequestResponse(c, msg, causa))
		default:
			// Otros errores seguramente sean de programación o de sistema
			ghLog.Errorf(c, "Error interno: %v", causa)
			c.PureJSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		}
	}
	c.Abort()
}

// Determina si un error es debido a un corte de red (broken pipe / connection reset by peer)
func isNetworkError(err error) bool {
	return errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
}
