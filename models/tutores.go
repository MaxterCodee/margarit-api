package models

import (
	"time"

	"gorm.io/gorm"
)

type Tutor struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Nombre       string     `gorm:"not null" json:"nombre"`
	Telefono     string     `gorm:"not null" json:"telefono"`
	Correo       string     `gorm:"not null" json:"correo"`
	EstudianteID uint       `gorm:"not null" json:"estudiante_id"`
	Estudiante   Estudiante `gorm:"foreignKey:EstudianteID" json:"estudiante"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (t *Tutor) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Tutor) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
