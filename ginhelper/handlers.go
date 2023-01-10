// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Devuelve el handler de ruta inxistente (a.k.a. 404)
func HandleNoRoute(debug bool, proxy string) gin.HandlerFunc {
	if !debug {
		return func(c *gin.Context) {
			// Producción
			ghLog.Infof("No encontrado: %s %s", c.Request.Method, c.Request.URL.Path)
			c.PureJSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
		}
	}
	if proxy == "" {
		return func(c *gin.Context) {
			// Depuración sin proxy
			ghLog.Infof("No implementado: %s %s", c.Request.Method, c.Request.URL.Path)
			c.PureJSON(http.StatusNotImplemented, gin.H{"error": "No implementado"})
		}
	}
	return func(c *gin.Context) {
		// Depuración con proxy
		ghLog.Infof("Redirección: %s %s%s", c.Request.Method, proxy, c.Request.URL.Path)
		director := func(req *http.Request) {
			proxyUrl, _ := url.Parse(proxy)
			req.URL.Scheme = proxyUrl.Scheme
			req.URL.Host = proxyUrl.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Gestiona la documentación swagger
func HandleSwagger(router *gin.Engine, ruta string, mapaSVG map[string]string) {
	svgHandler := func(c *gin.Context) {
		file := c.Param("any")
		if !strings.HasSuffix(file, ".svg") {
			return
		}
		file = strings.TrimPrefix(file, "/")
		c.Writer.Header().Add("Content-Type", "image/svg+xml")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.WriteString(mapaSVG[file])
		c.Abort()
	}
	ruta = path.Join(ruta, "swagger")
	router.GET(path.Join(ruta, "*any"), svgHandler, ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true), ginSwagger.DefaultModelsExpandDepth(-1)))
}
