package models

import (
	"time"

	"gorm.io/gorm"
)

type Aula struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nombre      string    `gorm:"not null" json:"nombre"`
	Descripcion string    `gorm:"not null" json:"descripcion"`
	Grupos      []Grupo   `gorm:"foreignKey:AulaID" json:"grupos"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Aula) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Aula) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}
