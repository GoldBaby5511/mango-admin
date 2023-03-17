package models

import (
	"database/sql/driver"
	"log"
	"time"

	"gorm.io/gorm"
)

type ControlBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
}

type ModelTime struct {
	CreatedAt Datetime       `json:"createdAt" gorm:"comment:创建时间" `
	UpdatedAt Datetime       `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type ID struct {
	Id uint32 `json:"id" gorm:"primaryKey; <-:false; type:int unsigned auto_increment; comment:主键;"`
}

func (id ID) IsValid() bool {
	return id.Id > 0
}

type Times struct {
	CreatedAt Datetime `json:"created_at" gorm:"<-:false; type:datetime; not null; default:now(); comment:创建时间;"`
	UpdatedAt Datetime `json:"updated_at" gorm:"<-:false; type:datetime; not null; default:now() ON UPDATE now(); comment:更新时间;"`
}

type Datetime string

// Value
func (d Datetime) Value() (driver.Value, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}

func (d Datetime) ToTime() time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", string(d), time.Local)
	if err != nil {
		log.Println(err)
	}
	return t
}
