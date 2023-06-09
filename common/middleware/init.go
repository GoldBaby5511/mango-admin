package middleware

import (
	"mango-admin/common/actions"

	"github.com/gin-gonic/gin"
	"mango-admin/pkg/sdk"
	jwt "mango-admin/pkg/jwtauth"
)

const (
	JwtTokenCheck   string = "JwtToken"
	RoleCheck       string = "AuthCheckRole"
	PermissionCheck string = "PermissionAction"
)

func InitMiddleware(r *gin.Engine) {
	// 自定义错误处理
	r.Use(CustomRecovery())

	// 数据库链接
	r.Use(WithContextDb)

	// 日志处理
	r.Use(ApiLog())

	// NoCache is a middleware function that appends headers
	r.Use(NoCache)
	// 跨域处理
	r.Use(Options)
	// Secure is a middleware function that appends security
	r.Use(Secure)
	//r.Use(DemoEvn())
	// 链路追踪
	//r.Use(middleware.Trace())
	sdk.Runtime.SetMiddleware(JwtTokenCheck, (*jwt.GinJWTMiddleware).MiddlewareFunc)
	sdk.Runtime.SetMiddleware(RoleCheck, AuthCheckRole())
	sdk.Runtime.SetMiddleware(PermissionCheck, actions.PermissionAction())
}
