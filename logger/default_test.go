package logger_test

import "github.com/horus-es/go-util/v2/logger"

func Example() {
	logger.InitLogger("", true)
	logger.Infof("Mensaje de %s", "información")
	logger.Warnf("Mensaje de %s", "advertencia")
	logger.Errorf("Mensaje de %s", "error") // Este aparece por STDERR
	// Output:
	// INFO: Mensaje de información
	// WARN: Mensaje de advertencia
}

func ExampleInfof() {
	logger.InitLogger("", true)
	logger.Infof("sin parámetros")
	logger.Infof("con parámetro %q", "parámetro")
	logger.InitLogger("", false)
	logger.Infof("con debug=false, no se registra el mensaje")
	// Output:
	// INFO: sin parámetros
	// INFO: con parámetro "parámetro"
}

func ExampleWarnf() {
	logger.InitLogger("", true)
	logger.Warnf("sin parámetros")
	logger.Warnf("con parámetro %q", "parámetro")
	logger.InitLogger("", false)
	logger.Warnf("con debug=false, se registra el mensaje en fichero o en STDOUT")
	// Output:
	// WARN: sin parámetros
	// WARN: con parámetro "parámetro"
	// WARN: con debug=false, se registra el mensaje en fichero o en STDOUT
}

func ExampleErrorf() {
	logger.InitLogger("", true)
	logger.Errorf("sin parámetros")
	logger.Errorf("con parámetro %q", "parámetro")
	logger.InitLogger("", false)
	logger.Errorf("con debug=false, se registra el mensaje en fichero o en STDERR")
	// Output:
}
