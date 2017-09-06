package middleware

import (
	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/common"
)

func CrossMiddleware() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		crossValue := c.GetHeader(common.CrossHeaderKey)
		if len(crossValue) == 0 {
			c.Header(common.CrossHeaderKey, "*")
		}
		c.Next()
	}

	return fn
}
