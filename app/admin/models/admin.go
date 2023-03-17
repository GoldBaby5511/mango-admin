package models

import (
	"go-admin/common/models"

	"github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"
)

var (
	StatisticPtr *Statistic
)

func admindb() *gorm.DB {
	return sdk.Runtime.GetDbByKey("default")
}

type Statistic struct {
	models.ID
	TheDate             string   `json:"the_date" gorm:"type:date; not null; comment:日期;"`
	SignupCount         int64    `json:"signup_count" gorm:""`
	RemainCount         IntArray `json:"remain_count" gorm:"size:100; comment:后7天留存;"`
	LoginCount          int64    `json:"login_count" gorm:"size:4; comment:登录人数;"`
	RechargeCount       int64    `json:"recharge_count"`
	RechargeUserCount   int64    `json:"recharge_user_count"`
	OxygenCount         int64    `json:"oxygen_count"`          // 用户现有氧气总量
	OxygenGenerateCount int64    `json:"oxygen_generate_count"` // 用户氧气产生
	OxygenConsumeCount  int64    `json:"oxygen_consume_count"`  // 用户氧气消耗
	LeafCount           int64    `json:"leaf_count"`            // 我们账户内叶总量
	LeafGenerateCount   int64    `json:"leaf_generate_count"`   // 用户叶子获得
	LeafConsumeCount    int64    `json:"leaf_consume_count"`    // 用户叶子消耗

	models.Times
}

func (us *Statistic) DB() *gorm.DB {
	return admindb().Model(us)
}
