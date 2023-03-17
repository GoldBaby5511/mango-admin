package router

import (
	"go-admin/app/admin/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerGameServerRouter)
}

func registerGameServerRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.GameServerApi{}
	r := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.Any("game_server/:path", api.Index)
	}

	r = v1.Group("game")
	{
		r.GET("id_list", api.GetPage)
		r.POST("generate_ids", api.GenerateIds)
	}
}
