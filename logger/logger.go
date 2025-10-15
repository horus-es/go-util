// Archivos de registro (Logger)
package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/v2/errores"
	"github.com/pkg/errors"
)

type Logger struct {
	debugging bool       // Modo depuración: si true graba los INFO
	filename  string     // Nombre del fichero sin el sufijo .log
	date      string     // Fecha de escritura del último log, usado para rotar
	file      *os.File   // Handler del fichero de salida, o nil para STDOUT/STDERR
	mutex     sync.Mutex // Para evitar que se mezclen bloques
}

type logBuffer struct {
	buf     bytes.Buffer // Buffer del log
	errores int          // Número de WARN/ERROR en el buffer, usado para decidir si se graba el bufffer
}

var defaultLogger *Logger // Logger por defecto

const (
	ctxCtxKey = "7d0952d6-1a52-4e1c-ae54-25efdf2e669a" // Clave única para el buffer del contexto
	kINFO     = "INFO: "
	kWARN     = "WARN: "
	kERROR    = "ERROR: "
	kTIME     = "2006-01-02 15:04:05.000 "
	kSEP      = "=================================================="
)

// Crea un logger. Si filename esta vacío, la salida se producirá por STDOUT y STDERR.
// A filename se le añade el sufijo .log automáticamente.
// Los ficheros se rotan diariamente con formato filename-XX.log, donde XX es el día del mes
func NewLogger(filename string, debug bool) *Logger {
	logger := Logger{}
	logger.debugging = debug
	logger.filename = filename
	// Abrimos fichero
	if logger.filename != "" {
		fn := logger.filename + ".log"
		stat, err := os.Stat(fn)
		if err == nil {
			logger.date = stat.ModTime().Format("20060102")
			logger.file, err = os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, 0666)
		} else if errors.Is(err, os.ErrNotExist) {
			logger.date = time.Now().Format("20060102")
			logger.file, err = os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0666)
		}
		errores.PanicIfError(err)
	}
	return &logger
}

// Inicializa el logger por defecto
func InitLogger(filename string, debug bool) {
	closeLogger(defaultLogger)
	defaultLogger = NewLogger(filename, debug)
}

// Cierra fichero de log
func closeLogger(logger *Logger) {
	if logger == nil || logger.file == nil {
		return
	}
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.file.Close()
	logger.file = nil
}

// Cierra fichero de log
func (logger *Logger) CloseLogger() {
	if logger == nil {
		closeLogger(defaultLogger)
	} else {
		closeLogger(logger)
	}
}

// Cierra fichero de log por defecto
func CloseLogger() {
	closeLogger(defaultLogger)
}

func getLogBuf(c *gin.Context) *logBuffer {
	p, _ := c.Get(ctxCtxKey)
	if p == nil {
		lb := &logBuffer{}
		c.Set(ctxCtxKey, lb)
		return lb
	} else {
		return p.(*logBuffer)
	}
}

// Escribe en el fichero de log o en buffer
func writeFileOrBuffer(c *gin.Context, logger *Logger, prefix, format string, v ...any) {
	ahora := time.Now()
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	if c == nil {
		rotateLog(logger, ahora)
		writeLine(logger.file, ahora.Format(kTIME)+prefix, format, v...)
		logger.file.Sync()
	} else {
		logBuf := getLogBuf(c)
		writeLine(&logBuf.buf, ahora.Format(kTIME)+prefix, format, v...)
		if prefix != kINFO {
			logBuf.errores++
		}
	}
}

// Escribe una línea con prefijo
func writeLine(w io.Writer, prefix, format string, v ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if len(v) > 0 {
		_, err := fmt.Fprintf(w, prefix+format, v...)
		errores.PanicIfError(err)
	} else {
		_, err := fmt.Fprint(w, prefix+format)
		errores.PanicIfError(err)
	}
}

// Rota el fichero de log
func rotateLog(logger *Logger, ahora time.Time) {
	hoy := ahora.Format("20060102")
	if logger.date == hoy {
		return
	}
	logger.file.Close()
	logger.file = nil
	f1 := logger.filename + ".log"
	f2 := logger.filename + "-" + logger.date[6:8] + ".log"
	err := os.Rename(f1, f2)
	errores.PanicIfError(err)
	logger.date = hoy
	logger.file, err = os.OpenFile(f1, os.O_CREATE|os.O_WRONLY, 0666)
	errores.PanicIfError(err)
}

// Graba el buffer a fichero (interno).
func flush(c *gin.Context, logger *Logger) {
	errores.PanicIfTrue(c == nil, "Flush requiere contexto")
	if logger == nil || logger.file == nil {
		fmt.Fprintln(os.Stdout, kSEP)
		return // Flush es solo para fichero
	}
	ahora := time.Now()
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logBuf := getLogBuf(c)
	if !logger.debugging && logBuf.errores == 0 {
		// No estamos en modo depuración y no hay ningún WARN o ERROR
		logBuf.buf.Reset()
		return
	}
	fmt.Fprintln(&logBuf.buf, kSEP)
	rotateLog(logger, ahora)
	_, err := logBuf.buf.WriteTo(logger.file)
	errores.PanicIfError(err)
	logger.file.Sync()
	logBuf.buf.Reset()
	logBuf.errores = 0
}

// Registra un INFO  (interno): por fichero si procede, en otro caso por STDOUT en modo depuración
func infof(c *gin.Context, logger *Logger, format string, v ...any) {
	if logger == nil {
		writeLine(os.Stdout, kINFO, format, v...)
	} else if logger.file != nil {
		writeFileOrBuffer(c, logger, kINFO, format, v...)
	} else if logger.debugging {
		writeLine(os.Stdout, kINFO, format, v...)
	}
}

// Registra un WARN (interno): siempre por STDOUT y por fichero si procede
func warnf(c *gin.Context, logger *Logger, format string, v ...any) {
	writeLine(os.Stdout, kWARN, format, v...)
	if logger != nil && logger.file != nil {
		writeFileOrBuffer(c, logger, kWARN, format, v...)
	}
}

// Registra un ERROR (interno): siempre por STDERR y por fichero si procede
func errorf(c *gin.Context, logger *Logger, format string, v ...any) {
	writeLine(os.Stderr, kERROR, format, v...)
	if logger != nil && logger.file != nil {
		writeFileOrBuffer(c, logger, kERROR, format, v...)
	}
}

// Registra un INFO
func (logger *Logger) Infof(c *gin.Context, format string, v ...any) {
	if logger == nil {
		infof(c, defaultLogger, format, v...)
	} else {
		infof(c, logger, format, v...)
	}
}

// Registra un INFO usando el logger por defecto
func Infof(c *gin.Context, format string, v ...any) {
	infof(c, defaultLogger, format, v...)
}

// Registra un WARN
func (logger *Logger) Warnf(c *gin.Context, format string, v ...any) {
	if logger == nil {
		warnf(c, defaultLogger, format, v...)
	} else {
		warnf(c, logger, format, v...)
	}
}

// Registra un WARN usando el logger por defecto
func Warnf(c *gin.Context, format string, v ...any) {
	warnf(c, defaultLogger, format, v...)
}

// Registra un ERROR
func (logger *Logger) Errorf(c *gin.Context, format string, v ...any) {
	if logger == nil {
		errorf(c, defaultLogger, format, v...)
	} else {
		errorf(c, logger, format, v...)
	}
}

// Registra un ERROR usando el logger por defecto
func Errorf(c *gin.Context, format string, v ...any) {
	errorf(c, defaultLogger, format, v...)
}

// Graba el buffer en fichero
func (logger *Logger) Flush(c *gin.Context) {
	if logger == nil {
		flush(c, defaultLogger)
	} else {
		flush(c, logger)
	}
}

// Graba el buffer en fichero usando el logger por defecto
func Flush(c *gin.Context) {
	flush(c, defaultLogger)
}
