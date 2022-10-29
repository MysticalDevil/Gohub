// Package models Model common properties and methods
package models

import "time"

// BaseModel Model base class
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField TimeStamp
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"created_at;index;" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"updated_at;index;" json:"updated_at,omitempty"`
}
