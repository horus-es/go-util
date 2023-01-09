package errores

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegistro(t *testing.T) {
	logger := NewLogger("", true)
	logger.Infof("Mensaje de %s", "información")
	logger.Warnf("Mensaje de %s", "advertencia")
	logger.Errorf("Mensaje de %s", "error")
}

func TestBadRequest(t *testing.T) {
	logger := NewLogger("", true)
	mensaje := "Mensaje para usuario"
	causa := "Causa del error para depuración"
	r := logger.BadHttpRequest(mensaje, errors.New(causa))
	assert.Equal(t, r["error"], mensaje)
	assert.Equal(t, r["causa"], causa)
	r = logger.BadHttpRequest(mensaje, nil)
	assert.Equal(t, r["error"], mensaje)
	assert.Equal(t, r["causa"], nil)
}

func ExampleLogger_Infof() {
	logger := NewLogger("", true)
	logger.Infof("sin parámetros")
	logger.Infof("con parámetro %q", "parámetro")
	logger = NewLogger("", false)
	logger.Infof("con debug=false, no se registra el mensaje")
	// Output:
	// INFO: sin parámetros
	// INFO: con parámetro "parámetro"
}

func ExampleLogger_Warnf() {
	logger := NewLogger("", true)
	logger.Warnf("sin parámetros")
	logger.Warnf("con parámetro %q", "parámetro")
	logger = NewLogger("", false)
	logger.Warnf("con debug=false, se registra el mensaje en fichero o en STDOUT")
	// Output:
	// WARN: sin parámetros
	// WARN: con parámetro "parámetro"
	// WARN: con debug=false, se registra el mensaje en fichero o en STDOUT
}

func ExampleLogger_Errorf() {
	logger := NewLogger("", true)
	logger.Errorf("sin parámetros")
	logger.Errorf("con parámetro %q", "parámetro")
	logger = NewLogger("", false)
	logger.Errorf("con debug=false, se registra el mensaje en fichero o en STDERR")
	// Output:
}

func ExampleLogger_BadHttpRequest() {
	logger := NewLogger("", true)
	mensaje := "Mensaje para usuario"
	causa := "Causa del error para depuración"
	r := logger.BadHttpRequest(mensaje, errors.New(causa))
	fmt.Println(r)
	// Output:
	// WARN: Mensaje para usuario: Causa del error para depuración
	// map[causa:Causa del error para depuración error:Mensaje para usuario]
}

func TestLogger(t *testing.T) {
	logger := NewLogger("testlog", true)
	logger.Infof("Prueba 1 de logger: %s", "info")
	logger.Warnf("Prueba 1 de logger: %s", "warn")
	logger.Errorf("Prueba 1 de logger: %s", "error")
	ayer := time.Now().AddDate(0, 0, -1)
	os.Chtimes("testlog.log", ayer, ayer)
	logger = NewLogger("testlog", true)
	logger.Infof("Prueba 2 de logger: %s", "info")
	logger.Warnf("Prueba 2 de logger: %s", "warn")
	logger.Errorf("Prueba 2 de logger: %s", "error")
}
