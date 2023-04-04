package dto

import (
	"mango-admin/app/admin/models"
	"mango-admin/common/dto"
)

type UserListReq struct {
	dto.Pagination `search:"-"`
}

type UserListResp struct {
	Count int64 `json:"count" `
	List  []struct {
		SysId      int    `json:"sys_id"`
		UserId     int    `json:"user_id"`
		Account    string `json:"account"`
		Nickname   string `json:"nickname"`
		Balance    uint64 `json:"balance"`
		IsOnline   bool   `json:"is_online" gorm:"-"`
		ChangeTime string `json:"change_time" gorm:"-"`
	} `json:"list"`
}

type DashboardResp struct {
	Cards  []DashboardRespCard `json:"cards"`
	Oxygen struct {
		Total   int64 `json:"total"`
		Consume int64 `json:"consume"`
	} `json:"o2"`
	Table DashboardTable `json:"table"`
	Leaf  struct {
		Total    int64 `json:"total"`
		Consume  int64 `json:"consume"`
		Recovery int64 `json:"recovery"`
	} `json:"leaf"`
	NFT struct {
		Total int `json:"total"`
	} `json:"nft"`
}

type DashboardTableReq struct {
	LastDays int `form:"lastDays" binding:"required"`
}

type DashboardTableResp = DashboardTable

type DashboardTable []struct {
	TheDate           string `json:"the_date"`
	SignupCount       int64  `json:"signup_count" `
	LoginCount        int64  `json:"login_count"`
	RechargeCount     int64  `json:"recharge_count" `
	RechargeUserCount int64  `json:"recharge_user_count"`
}
type DashboardOnlineReq struct {
	Start string `form:"start"`
	End   string `form:"end"`
}

type DashboardOnlineResp struct {
	Avg    int64           `json:"avg"`
	XAxios []string        `json:"xAxios"`
	YAxios [][]int         `json:"yAxios"`
	List   [][]interface{} `json:"list"`
}

type DashboardRespCard struct {
	Date     string `json:"date"`
	Count    int64  `json:"count"`
	Compare1 int64  `json:"compare1"`
	Compare2 int64  `json:"compare2"`
}

type WholeDataReq struct {
	Start string `form:"start"`
	End   string `form:"end"`
}

type WholeDataResp struct {
	// 7日留存
	Remain []struct {
		TheDate     string          `json:"the_date"`
		RemainCount models.IntArray `json:"remain_count"`
	} `json:"remain"`
	// 次日流失
	Loss struct {
		Avg  float32         `json:"avg"`
		List [][]interface{} `json:"list"`
	} `json:"loss"`
	// 来源分布
	Source []struct {
		Source string `json:"source"`
		Count  int    `json:"count"`
	} `json:"source"`
	// // 游戏时长
	// GameTime []struct {
	// 	Type1 string `json:"type1"`
	// 	Type2 string `json:"type2"`
	// 	Type3 string `json:"type3"`
	// }
	// // 在线城市
	// OnlineCity []struct {
	// 	City  string `json:"city"`
	// 	Count int    `json:"count"`
	// }
	// // 城市分布
	// City []struct {
	// 	City  string `json:"city"`
	// 	Count int    `json:"count"`
	// }
}

type ActiveDataReq struct {
	Date        string `form:"date" binding:"required"`
	StartMinute string `form:"startMinute" binding:"required"`
}

type ActiveQueryList []struct {
	T string `json:"t"`
	N int    `json:"n"`
}

type ActiveDataResp struct {
	Total int         `json:"total"`
	List  [][2]string `json:"list"`
}

type OnlineDataReq = ActiveDataReq

type OnlineDataResp = ActiveDataResp

type RechargeDataReq = ActiveDataReq

type RechargeDataResp = ActiveDataResp

type DailyStatisticsListReq struct {
	dto.Pagination `search:"-"`
}

type DailyStatistics struct {
	SysId           int    `json:"sys_id" gorm:"-"`
	Date            string `json:"date" gorm:"->"`
	RegisterCount   int    `json:"register_count" gorm:"->"`
	LoginCount      int    `json:"login_count" gorm:"-"`
	FirstLoginCount int    `json:"first_time_login_count" gorm:"-"`
	MaxOnlineCount  int    `json:"max_online_count" gorm:"-"`
	OnlineCount1    int    `json:"online_count_3" gorm:"-"`
	OnlineCount2    int    `json:"online_count_15" gorm:"-"`
	OnlineCount3    int    `json:"online_count_30" gorm:"-"`
}

type DailyStatisticsListResq struct {
	Count            int64             `json:"count" gorm:"-"`
	SumRegisterCount int64             `json:"sum_register_count" gorm:"-"`
	SumLoginCount    int64             `json:"sum_login_count" gorm:"-"`
	List             []DailyStatistics `json:"list"`
}
