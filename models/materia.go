package models

import (
	"time"

	"gorm.io/gorm"
)

type Materia struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Titulo      string    `gorm:"not null" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	GradoID     uint      `gorm:"not null" json:"grado_id"`
	Grado       Grado     `gorm:"foreignKey:GradoID" json:"grado"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Materia) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Materia) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
