package models

import (
	"mango-admin/pkg/sdk"
	"gorm.io/gorm"
)

var (
	WealthPtr          *Wealth
	WealthChangeLogPtr *WealthChangeLog
)

func propertydb() *gorm.DB {
	return sdk.Runtime.GetDbByKey("property")
}

type Wealth struct {
	SysId  uint64 `json:"sys_id"`
	UserId uint64 `json:"user_id"`
	Leaf   uint64 `json:"leaf"`
	Oxygen uint64 `json:"oxygen"`
}

func (w *Wealth) DB() *gorm.DB {
	return propertydb().Model(w)
}

type WealthChangeLog struct {
	SysId       uint64 `json:"sys_id"`
	UserId      uint64 `json:"user_id"`
	ChangeId    uint8  `json:"change_id"` // 1=leaf、2=oxygen、3=water
	ChangeCount int64  `json:"change_count"`
}

func (wcl *WealthChangeLog) DB() *gorm.DB {
	return propertydb().Model(wcl)
}
