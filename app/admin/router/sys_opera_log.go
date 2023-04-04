package router

import (
	"github.com/gin-gonic/gin"
	"mango-admin/app/admin/apis"
	"mango-admin/common/middleware"
	jwt "mango-admin/pkg/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysOperaLogRouter)
}

// 需认证的路由代码
func registerSysOperaLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysOperaLog{}
	r := v1.Group("/sys-opera-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.DELETE("", api.Delete)
	}
}
