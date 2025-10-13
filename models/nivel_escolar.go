package models

import (
	"time"

	"gorm.io/gorm"
)

type NivelEscolar struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Titulo      string    `gorm:"not null" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	Mensualidad float64   `gorm:"not null" json:"mensualidad"`
	PlantelID   uint      `gorm:"not null" json:"plantel_id"`
	Plantel     Plantel   `gorm:"foreignKey:PlantelID" json:"plantel"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (n *NivelEscolar) BeforeCreate(tx *gorm.DB) error {
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	return nil
}

func (n *NivelEscolar) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = time.Now()
	return nil
}
