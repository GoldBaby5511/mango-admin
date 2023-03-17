package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
)

func WithContextDb(c *gin.Context) {
	c.Set("db", sdk.Runtime.GetDbByKey("default").WithContext(c))
	c.Next()
}
