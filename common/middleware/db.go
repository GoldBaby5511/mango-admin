package middleware

import (
	"github.com/gin-gonic/gin"
	"mango-admin/pkg/sdk"
)

func WithContextDb(c *gin.Context) {
	c.Set("db", sdk.Runtime.GetDbByKey("default").WithContext(c))
	c.Next()
}
