package apis

import (
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type GameServerApi struct {
	api.Api
}

// 代理到游戏服务控制
func (e GameServerApi) Index(c *gin.Context) {
	if user.GetRoleName(c) != "admin" && user.GetRoleName(c) != "系统管理员" {
		e.Error(403, nil, "权限不足")
		return
	}

	url := "http://127.0.0.1:15052/"
	path := c.Param("path")

	req, _ := http.NewRequest(c.Request.Method, url+path, c.Request.Body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.Error(500, nil, "游戏服务接口出错: "+err.Error())
		return
	}
	defer resp.Body.Close()

	io.Copy(c.Writer, resp.Body)
}

// 生成游戏ID -
func (e GameServerApi) GetPage(c *gin.Context) {
	s := service.Game{}
	req := dto.GameIdListReq{}
	err := e.MakeContext(c).
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	s.Orm = sdk.Runtime.GetDbByKey("user")

	list := make([]models.GameIdNormal, 0)
	var count int64
	err = s.GetPage(&req, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	max := s.GetMax()
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), strconv.Itoa(int(max)))
}

func (e GameServerApi) GenerateIds(c *gin.Context) {
	var req struct {
		Start int `form:"start" json:"start"`
		End   int `form:"end" json:"end"`
	}
	s := service.Game{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	s.Orm = sdk.Runtime.GetDbByKey("user")

	err = s.GenerateId(req.Start, req.End)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	e.OK(nil, "生成成功")
}
