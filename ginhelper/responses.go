// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/horus-es/go-util/v3/logger"
)

var ghLog *logger.Logger

// Establece el logger. Si el logger es nil, se usa el logger por defecto.
func InitGinHelper(logger *logger.Logger) {
	ghLog = logger
}

// Genera una respuesta json/REST a una solicitud incorrecta.
// Incluye un mensaje de error para el usuario y opcionalmente la causa del error para depuración.
func BadRequestResponse(c *gin.Context, msg string, causa any) map[string]any {
	if causa == nil {
		ghLog.Warnf(c, "%s", msg)
		return map[string]any{"error": msg}
	}
	ghLog.Warnf(c, "%s: %v", msg, causa)
	return map[string]any{"error": msg, "causa": fmt.Sprint(causa)}
}
