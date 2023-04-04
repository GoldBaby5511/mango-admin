package router

import (
	"github.com/gin-gonic/gin"
	jwt "mango-admin/pkg/sdk/pkg/jwtauth"

	"mango-admin/app/admin/apis"
	"mango-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysRoleRouter)
}

// 需认证的路由代码
func registerSysRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysRole{}
	r := v1.Group("/role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
	r1 := v1.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		r1.PUT("/role-status", api.Update2Status)
		r1.PUT("/roledatascope", api.Update2DataScope)
	}
}
