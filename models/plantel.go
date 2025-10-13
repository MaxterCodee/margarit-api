package models

import (
	"time"

	"gorm.io/gorm"
)

type Plantel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nombre      string    `gorm:"not null" json:"nombre"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	Ubicacion   string    `gorm:"not null" json:"ubicacion"`
	Telefono    string    `gorm:"not null" json:"telefono"`
	Correo      string    `gorm:"not null" json:"correo"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Plantel) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Plantel) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
