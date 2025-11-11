// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
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
		defer recuperaLogger(c)
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

// Auxiliar de MiddlewareLogger
func recuperaLogger(c *gin.Context) {
	causa := recover()
	if causa == nil {
		// No panic
		return
	}
	// Registramos el error y el stack
	ghLog.Errorf(c, "panic: %v\n%s", causa, debug.Stack())
	ghLog.Flush(c)
	panic(causa)
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
		defer recuperaPanic(c)
		c.Next()
	}
}

// Auxiliar de MiddlewarePanic
func recuperaPanic(c *gin.Context) {
	causa := recover()
	if causa == nil {
		// No panic
		return
	}
	e, ok := causa.(error)
	if !ok {
		ghLog.Errorf(c, "panic: %v\n%s", causa, debug.Stack())
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
	} else {
		// ¿Es un error de red?
		if isNetworkError(e) {
			// Como hay error de red no podemos responder nada ...
			ghLog.Errorf(c, "Error de red: %v", causa)
		} else {
			errorSQL, msg := postgres.GetErrorSQL(e)
			switch errorSQL {
			case postgres.INTEGRITY_CONSTRAINT_VIOLATION:
				c.PureJSON(http.StatusBadRequest, BadRequestResponse(c, "Valor duplicado", causa))
			case postgres.PL_PGSQL_RAISE_EXCEPTION:
				c.PureJSON(http.StatusBadRequest, BadRequestResponse(c, msg, causa))
			default:
				ghLog.Errorf(c, "panic: %v\n%s", causa, debug.Stack())
				c.PureJSON(http.StatusInternalServerError, gin.H{"error": "Error interno", "causa": fmt.Sprint(causa)})
			}
		}
	}
	c.Abort()
}

// Determina si un error es debido a un corte de red (broken pipe / connection reset by peer)
func isNetworkError(err error) bool {
	return errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
}
