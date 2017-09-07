package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	crossRequestOriginKey  = "Origin"
	crossRequestMethodsKey = "Access-Control-Request-Method"
	crossRequestHeadersKey = "Access-Control-Request-Headers"
	crossAllowOriginKey    = "Access-Control-Allow-Origin"
	crossAllowMethodsKey   = "Access-Control-Allow-Methods"
	crossAllowHeadersKey   = "Access-Control-Allow-Headers"
)

func CrossMiddleware() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		allowOrigin := c.GetHeader(crossRequestOriginKey)
		if len(allowOrigin) != 0 {
			c.Header(crossAllowOriginKey, allowOrigin)
		}
		allowMethods := c.GetHeader(crossRequestMethodsKey)
		if len(allowMethods) != 0 {
			c.Header(crossAllowMethodsKey, allowMethods)
		}
		allowHeaders := c.GetHeader(crossRequestHeadersKey)
		if len(allowHeaders) != 0 {
			c.Header(crossAllowHeadersKey, allowHeaders)
		}
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.Next()
		}
	}

	return fn
}
