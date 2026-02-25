// Package models Model common properties and methods
package models

import (
	"time"

	"github.com/spf13/cast"
)

// BaseModel Model base class
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField TimeStamp
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"created_at;index;" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at;index;" json:"updated_at"`
}

// GetStringID Get ID in string format
func (a BaseModel) GetStringID() string {
	return cast.ToString(a.ID)
}
