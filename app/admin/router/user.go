package router

import (
	"go-admin/app/admin/apis"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerUserRouter)
}

func registerUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.User{}
	r := v1.Group("/user")
	{
		r.GET("list", api.List)
		r.GET("dashboard", api.Darshboard)
		r.GET("dashboard-table", api.DashboardTable)
		r.GET("dashboard-online", api.DashboardOnline)
		r.GET("whole-data", api.WholeData)
		r.GET("online-data", api.OnlineData)
		r.GET("recharge-data", api.RechargeData)
		r.GET("statistic", api.Statistic)
		r.GET("statistic-remain", api.StatisticRemainCount)
		r.GET("daily-statistics", api.DailyStatistics)
	}
}
