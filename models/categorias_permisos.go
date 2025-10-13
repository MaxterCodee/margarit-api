package models

import (
	"time"

	"gorm.io/gorm"
)

type CategoriaPermiso struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Titulo      string         `gorm:"not null" json:"titulo"`
	Descripcion string         `gorm:"type:text" json:"descripcion"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *CategoriaPermiso) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *CategoriaPermiso) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
