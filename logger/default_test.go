package logger

func Example() {
	InitLogger("", true)
	Infof("Mensaje de %s", "información")
	Warnf("Mensaje de %s", "advertencia")
	Errorf("Mensaje de %s", "error") // Este aparece por STDERR
	// Output:
	// INFO: Mensaje de información
	// WARN: Mensaje de advertencia
}

func ExampleInfof() {
	InitLogger("", true)
	Infof("sin parámetros")
	Infof("con parámetro %q", "parámetro")
	InitLogger("", false)
	Infof("con debug=false, no se registra el mensaje")
	// Output:
	// INFO: sin parámetros
	// INFO: con parámetro "parámetro"
}

func ExampleWarnf() {
	InitLogger("", true)
	Warnf("sin parámetros")
	Warnf("con parámetro %q", "parámetro")
	InitLogger("", false)
	Warnf("con debug=false, se registra el mensaje en fichero o en STDOUT")
	// Output:
	// WARN: sin parámetros
	// WARN: con parámetro "parámetro"
	// WARN: con debug=false, se registra el mensaje en fichero o en STDOUT
}

func ExampleErrorf() {
	InitLogger("", true)
	Errorf("sin parámetros")
	Errorf("con parámetro %q", "parámetro")
	InitLogger("", false)
	Errorf("con debug=false, se registra el mensaje en fichero o en STDERR")
	// Output:
}
