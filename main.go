package main

import (
	"mango-admin/cmd"
)

//go:generate swag init --parseDependency --parseDepth=6

// @title mango-admin API
// @version 2.0.0
// @description 基于Gin + Vue + Element UI的前后端分离权限管理系统的接口文档
// @description 添加qq群: 521386980 进入技术交流群 请先star，谢谢！
// @license.name MIT
// @license.url https://github.com/mango-admin-team/mango-admin/blob/master/LICENSE.md

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
