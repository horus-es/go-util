// Archivos de registro (Logger)
package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Logger struct {
	debugging bool
	file      string
	date      string
	writer    *os.File
	mutex     sync.Mutex
}

var defaultLogger *Logger // Logger por defecto

// Crea un logger. Si file=="", la salida se producirá por consola.
// A file se le añade el sufijo .log automáticamente.
func NewLogger(file string, debug bool) *Logger {
	logger := Logger{}
	logger.debugging = debug
	logger.file = file
	if logger.file != "" {
		file := logger.file + ".log"
		stat, err := os.Stat(file)
		if err == nil {
			logger.date = stat.ModTime().Format("20060102")
			logger.writer, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0666)
		} else if errors.Is(err, os.ErrNotExist) {
			logger.date = time.Now().Format("20060102")
			logger.writer, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
		}
		if err != nil {
			panic("Error abriendo " + file)
		}
	}
	return &logger
}

// Inicializa el logger por defecto
func InitLogger(file string, debug bool) {
	defaultLogger.CloseLogger()
	defaultLogger = NewLogger(file, debug)
}

func (logger *Logger) CloseLogger() {
	if logger == nil || logger.writer == nil {
		return
	}
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	if logger.writer == nil {
		return
	}
	logger.writer.Close()
	logger.writer = nil
}

// Escribe en el fichero de log, rotándolo si es preciso.
func writeFile(logger *Logger, prefix, format string, v ...any) bool {
	if logger == nil || logger.writer == nil {
		return false
	}
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	ahora := time.Now()
	hoy := ahora.Format("20060102")
	if logger.date != hoy {
		logger.writer.Close()
		logger.writer = nil
		f1 := logger.file + ".log"
		f2 := logger.file + "-" + logger.date[6:8] + ".log"
		err := os.Rename(f1, f2)
		if err == nil {
			logger.date = hoy
			logger.writer, err = os.OpenFile(f1, os.O_CREATE|os.O_WRONLY, 0666)
		}
		if err != nil {
			logger.writer = nil
			return false
		}
	}
	ok := writeLine(logger.writer, ahora.Format("2006-01-02 15:04:05 ")+prefix, format, v...)
	logger.writer.Sync()
	return ok
}

// Escribe una línea con prefijo
func writeLine(w io.Writer, prefix, format string, v ...any) bool {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if len(v) > 0 {
		_, err := fmt.Fprintf(w, prefix+format, v...)
		return err == nil
	} else {
		_, err := fmt.Fprint(w, prefix+format)
		return err == nil
	}
}

// Registra un INFO (interno)
func infof(logger *Logger, format string, v ...any) {
	if logger != nil && !logger.debugging {
		return
	}
	prefix := "INFO: "
	writeLine(os.Stdout, prefix, format, v...)
	writeFile(logger, prefix, format, v...)
}

// Registra un INFO (solo en modo DEBUG)
func (logger *Logger) Infof(format string, v ...any) {
	if logger == nil {
		infof(defaultLogger, format, v...)
	} else {
		infof(logger, format, v...)
	}
}

// Registra un INFO usando el logger por defecto (solo en modo DEBUG)
func Infof(format string, v ...any) {
	infof(defaultLogger, format, v...)
}

// Registra un WARN (interno)
func warnf(logger *Logger, format string, v ...any) {
	prefix := "WARN: "
	if !writeFile(logger, prefix, format, v...) || logger.debugging {
		writeLine(os.Stdout, prefix, format, v...)
	}
}

// Registra un WARN
func (logger *Logger) Warnf(format string, v ...any) {
	if logger == nil {
		warnf(defaultLogger, format, v...)
	} else {
		warnf(logger, format, v...)
	}
}

// Registra un WARN usando el logger por defecto
func Warnf(format string, v ...any) {
	warnf(defaultLogger, format, v...)
}

// Registra un ERROR (interno)
func errorf(logger *Logger, format string, v ...any) {
	prefix := "ERROR: "
	if !writeFile(logger, prefix, format, v...) || logger.debugging {
		writeLine(os.Stderr, prefix, format, v...)
	}
}

// Registra un ERROR
func (logger *Logger) Errorf(format string, v ...any) {
	if logger == nil {
		errorf(defaultLogger, format, v...)
	} else {
		errorf(logger, format, v...)
	}
}

// Registra un ERROR usando el logger por defecto
func Errorf(format string, v ...any) {
	errorf(defaultLogger, format, v...)
}
