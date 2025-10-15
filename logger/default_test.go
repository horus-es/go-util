package logger_test

import (
	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/v3/logger"
)

func Example() {
	logger.InitLogger("", true)
	logger.Infof(nil, "Mensaje de %s", "información")
	logger.Warnf(nil, "Mensaje de %s", "advertencia")
	logger.Errorf(nil, "Mensaje de %s", "error") // Este aparece por STDERR
	logger.CloseLogger()
	// Output:
	// INFO: Mensaje de información
	// WARN: Mensaje de advertencia
}

func ExampleInfof() {
	logger.InitLogger("", true)
	logger.Infof(nil, "sin parámetros")
	logger.Infof(nil, "con parámetro %q", "parámetro")
	logger.InitLogger("", false)
	logger.Infof(nil, "con debug=false, no se registra el mensaje")
	logger.CloseLogger()
	// Output:
	// INFO: sin parámetros
	// INFO: con parámetro "parámetro"
}

func ExampleWarnf() {
	logger.InitLogger("", true)
	logger.Warnf(nil, "sin parámetros")
	logger.Warnf(nil, "con parámetro %q", "parámetro")
	logger.InitLogger("", false)
	logger.Warnf(nil, "con debug=false, se registra el mensaje en fichero o en STDOUT")
	logger.CloseLogger()
	// Output:
	// WARN: sin parámetros
	// WARN: con parámetro "parámetro"
	// WARN: con debug=false, se registra el mensaje en fichero o en STDOUT
}

func ExampleErrorf() {
	logger.InitLogger("", true)
	logger.Errorf(nil, "sin parámetros")
	logger.Errorf(nil, "con parámetro %q", "parámetro")
	logger.InitLogger("", false)
	logger.Errorf(nil, "con debug=false, se registra el mensaje en fichero o en STDERR")
	logger.CloseLogger()
	// Output:
}

func ExampleFlush() {
	logger.InitLogger("testlog", false)
	c := &gin.Context{}
	logger.Errorf(c, "sin parámetros")
	logger.Warnf(c, "con parámetro %q", "parámetro")
	logger.Infof(c, "esta linea se incluye porque hay ERROR y WARN en el buffer")
	logger.Flush(c)
	logger.Infof(c, "esta linea no se incluye porque no hay ni ERROR ni WARN")
	logger.Flush(c)
	logger.CloseLogger()
	logger.InitLogger("testlog", true)
	logger.Infof(c, "esta linea se incluye porque debug=true")
	logger.Flush(c)
	logger.CloseLogger()
	// Output:
	// WARN: con parámetro "parámetro"
}
