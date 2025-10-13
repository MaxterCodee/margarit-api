package models

import (
	"time"

	"gorm.io/gorm"
)

type EstatusEmpleado struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *EstatusEmpleado) BeforeCreate(tx *gorm.DB) error {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return nil
}

func (e *EstatusEmpleado) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}
