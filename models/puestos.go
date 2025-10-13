package models

import (
	"time"

	"gorm.io/gorm"
)

type Puesto struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	PagoXHr   float64   `gorm:"not null" json:"pago_x_hr"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Puesto) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Puesto) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
