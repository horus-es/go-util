// Funciones auxiliares para GIN-GONIC
package ginhelper

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
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
