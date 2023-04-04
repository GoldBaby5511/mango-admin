package apis

import (
	"mango-admin/app/admin/service"
	"mango-admin/app/admin/service/dto"

	"github.com/gin-gonic/gin"
	"mango-admin/pkg/sdk/pkg/response"
)

type User struct {
	Service service.User
}

// 参数验证，务必在false时return。 若非强校验，使用c.ShouldBind(&req)
func Bind(c *gin.Context, to interface{}) bool {
	if err := c.ShouldBind(to); err != nil {
		response.Error(c, 400, err, err.Error())
		return false
	}
	return true
}

func Echo(c *gin.Context, resp interface{}, err error) {
	if err != nil {
		response.Error(c, 500, err, err.Error())
		return
	}
	response.OK(c, resp, "ok")
}

func (u User) List(c *gin.Context) {
	var req dto.UserListReq
	if !Bind(c, &req) {
		return
	}
	var resp dto.UserListResp
	err := u.Service.List(req, &resp)
	Echo(c, &resp, err)
}

func (u User) Darshboard(c *gin.Context) {
	var resp dto.DashboardResp
	err := u.Service.Dashboard(&resp)
	Echo(c, &resp, err)
}

func (u User) DashboardTable(c *gin.Context) {
	var req dto.DashboardTableReq
	if !Bind(c, &req) {
		return
	}
	resp, err := u.Service.DashboardTable(req)
	Echo(c, &resp, err)
}

func (u User) DashboardOnline(c *gin.Context) {
	var req dto.DashboardOnlineReq
	if !Bind(c, &req) {
		return
	}
	resp, err := u.Service.DashboardOnline(req)
	Echo(c, &resp, err)
}

func (u User) WholeData(c *gin.Context) {
	var req dto.WholeDataReq
	if !Bind(c, &req) {
		return
	}
	resp, err := u.Service.WholeData(req)
	Echo(c, &resp, err)
}

// OnlineData
func (u User) OnlineData(c *gin.Context) {
	var req dto.OnlineDataReq
	if !Bind(c, &req) {
		return
	}
	var resp dto.OnlineDataResp
	err := u.Service.OnlineData(req, &resp)
	Echo(c, &resp, err)
}

// RechargeData
func (u User) RechargeData(c *gin.Context) {
	var req dto.RechargeDataReq
	if !Bind(c, &req) {
		return
	}
	var resp dto.RechargeDataResp
	err := u.Service.RechargeData(req, &resp)
	Echo(c, &resp, err)
}

// Statistic
func (u User) Statistic(c *gin.Context) {
	var req struct {
		Date string `form:"date" binding:"required"`
	}
	if !Bind(c, &req) {
		return
	}
	u.Service.Statistic(req.Date, true)
	Echo(c, "ok", nil)
}

// StatisticRemainCount
func (u User) StatisticRemainCount(c *gin.Context) {
	var req struct {
		Date string `form:"date" binding:"required"`
	}
	if !Bind(c, &req) {
		return
	}
	for i := 1; i <= 7; i++ {
		u.Service.StatisticRemainCount(req.Date, i)
	}
	Echo(c, "ok", nil)
}

func (u User) DailyStatistics(c *gin.Context) {
	var req dto.DailyStatisticsListReq
	if !Bind(c, &req) {
		return
	}
	var resp dto.DailyStatisticsListResq
	err := u.Service.DailyStatistics(req, &resp)
	Echo(c, &resp, err)
}
