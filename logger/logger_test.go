package logger_test

import (
	"os"
	"testing"
	"time"

	"github.com/horus-es/go-util/v2/logger"
)

func ExampleLogger() {
	logger := logger.NewLogger("", true)
	logger.Infof(nil, "Mensaje de %s", "información")
	logger.Warnf(nil, "Mensaje de %s", "advertencia")
	logger.Errorf(nil, "Mensaje de %s", "error") // Este aparece por STDERR
	// Output:
	// INFO: Mensaje de información
	// WARN: Mensaje de advertencia
}

func ExampleLogger_Infof() {
	log := logger.NewLogger("", true)
	log.Infof(nil, "sin parámetros")
	log.Infof(nil, "con parámetro %q", "parámetro")
	log = logger.NewLogger("", false)
	log.Infof(nil, "con debug=false, no se registra el mensaje")
	// Output:
	// INFO: sin parámetros
	// INFO: con parámetro "parámetro"
}

func ExampleLogger_Warnf() {
	log := logger.NewLogger("", true)
	log.Warnf(nil, "sin parámetros")
	log.Warnf(nil, "con parámetro %q", "parámetro")
	log = logger.NewLogger("", false)
	log.Warnf(nil, "con debug=false, se registra el mensaje en fichero o en STDOUT")
	// Output:
	// WARN: sin parámetros
	// WARN: con parámetro "parámetro"
	// WARN: con debug=false, se registra el mensaje en fichero o en STDOUT
}

func ExampleLogger_Errorf() {
	log := logger.NewLogger("", true)
	log.Errorf(nil, "sin parámetros")
	log.Errorf(nil, "con parámetro %q", "parámetro")
	log = logger.NewLogger("", false)
	log.Errorf(nil, "con debug=false, se registra el mensaje en fichero o en STDERR")
	// Output:
}

func TestRotacion(t *testing.T) {
	log := logger.NewLogger("testlog", true)
	log.Infof(nil, "Prueba 1 de logger: %s", "info")
	log.Warnf(nil, "Prueba 1 de logger: %s", "warn")
	log.Errorf(nil, "Prueba 1 de logger: %s", "error")
	log.CloseLogger()
	ayer := time.Now().AddDate(0, 0, -1)
	os.Chtimes("testlog.log", ayer, ayer)
	log = logger.NewLogger("testlog", true)
	log.Infof(nil, "Prueba 2 de logger: %s", "info")
	log.Warnf(nil, "Prueba 2 de logger: %s", "warn")
	log.Errorf(nil, "Prueba 2 de logger: %s", "error")
	log.CloseLogger()
}
