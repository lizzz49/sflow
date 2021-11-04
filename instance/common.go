package instance

import (
	"gorm.io/gorm"
	"time"
)

type ControlBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
