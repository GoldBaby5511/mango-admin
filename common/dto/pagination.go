package dto

import (
	"gorm.io/gorm"
)

type Pagination struct {
	PageIndex int `form:"pageIndex" json:"pageIndex"`
	PageSize  int `form:"pageSize" json:"pageSize"`

	Count int64       `json:"count" form:"-"`
	List  interface{} `json:"list" form:"-"`
}

func (m *Pagination) GetPageIndex() int {
	if m.PageIndex <= 0 {
		m.PageIndex = 1
	}
	return m.PageIndex
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

func (p *Pagination) AutoFind(tx *gorm.DB, to interface{}) {
	tx.Count(&p.Count)
	if p.Count > 0 {
		offset := (p.GetPageIndex() - 1) * p.GetPageSize()
		tx.Offset(offset).Limit(p.GetPageSize()).Find(to)
	}
}
