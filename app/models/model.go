package models

import (
	"goblog/pkg/types"
	"time"
)

type BaseModel struct {
	ID uint64 `gorm:"cloumn:id;primaryKey;autoIncrement;not null"`

	CreatedAt time.Time `gorm:"cloumn:created_at;index"`
	UpdatedAt time.Time `gorm:"cloumn:updated_at;index"`
}

func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}
