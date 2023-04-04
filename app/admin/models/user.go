package models

import (
	"mango-admin/pkg/sdk"
	"gorm.io/gorm"
)

var (
	UserAccountPtr *UserAccount
	UserLoginPtr   *UserLogin
	UserOnlinePtr  *UserOnline
)

func userdb() *gorm.DB {
	return sdk.Runtime.GetDbByKey("user")
}

type UserAccount struct {
	UserId     int    `json:"user_id"`
	Account    string `json:"account"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Balance    uint64 `json:"balance"`
	CreateTime string `json:"create_time"`
}

func (ua *UserAccount) DB() *gorm.DB {
	return userdb().Model(ua)
}

type UserLogin struct {
	UserId      int    `json:"user_id"`
	LogoutTime  string `json:"logout_time"`
	RefreshTime string `json:"refresh_time"`
	CreateTime  string `json:"created_time"`
}

func (ul *UserLogin) DB() *gorm.DB {
	return userdb().Model(ul)
}

type UserOnline struct {
	SysId        uint64 `json:"sys_id"`
	UserId       uint64 `json:"user_id"`
	OfflineTime  uint32 `json:"offline_time" `
	OnlineSecond uint32 `json:"online_second" `
	CreateTime   string `json:"created_time"`
	ChangeTime   string `json:"change_time"`
}

func (uo *UserOnline) DB() *gorm.DB {
	return userdb().Model(uo)
}

type GameIdNormal struct {
	SysId  uint64 `json:"sys_id"`
	GameId uint64 `json:"game_id" `
	UserId uint64 `json:"user_id" `
}

type GameIdExcellent struct {
	SysId  uint64 `json:"sys_id"`
	GameId uint64 `json:"game_id" `
	UserId uint64 `json:"user_id" `
}

func (gi *GameIdNormal) DB() *gorm.DB {
	return userdb().Model(gi)
}

func (ri *GameIdExcellent) DB() *gorm.DB {
	return userdb().Model(ri)
}

func logicdb() *gorm.DB {
	return sdk.Runtime.GetDbByKey("logic")
}

type UserLoginLog struct {
	UserId int `json:"user_id"`
	//Account  string `json:"account"`
	Reason         int    `json:"reason"`
	LoginTime      string `json:"login_time"`
	LogoutTime     string `json:"logout_time"`
	Date           string `json:"date"`
	ActiveDuration int64  `json:"active_duration" gorm:"-"` //活跃时长
}

func (ull *UserLoginLog) DB() *gorm.DB {
	return logicdb().Model(ull).Table("login_log")
}

type UserLogicAccount struct {
	UserId int `json:"user_id"`
	//Account  string `json:"account"`
	RegTime string `json:"reg_time"`
}

func (ula *UserLogicAccount) DB() *gorm.DB {
	return logicdb().Model(ula).Table("account")
}
