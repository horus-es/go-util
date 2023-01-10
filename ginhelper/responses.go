// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"fmt"

	"github.com/horus-es/go-util/v2/logger"
)

var ghLog *logger.Logger

// Establece el logger. Si logger es nil, todos los mensajes se muestran en la consola.
func InitGinHelper(logger *logger.Logger) {
	ghLog = logger
}

// Genera una respuesta json/REST a una solicitud incorrecta.
// Incluye un mensaje de error para el usuario y opcionalmente la causa del error para depuración.
func BadRequestResponse(msg string, causa any) map[string]any {
	if causa == nil {
		ghLog.Warnf("%s", msg)
		return map[string]any{"error": msg}
	}
	ghLog.Warnf("%s: %v", msg, causa)
	return map[string]any{"error": msg, "causa": fmt.Sprint(causa)}
}
