package models

import (
	"time"

	"gorm.io/gorm"
)

type Condicion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Titulo      string    `gorm:"not null" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	ContratoID  uint      `gorm:"not null" json:"contrato_id"`
	Contrato    Contrato  `gorm:"foreignKey:ContratoID" json:"contrato"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Condicion) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Condicion) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
