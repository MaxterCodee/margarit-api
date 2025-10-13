// grados_academicos.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type GradoAcademico struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (g *GradoAcademico) BeforeCreate(tx *gorm.DB) error {
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return nil
}

func (g *GradoAcademico) BeforeUpdate(tx *gorm.DB) error {
	g.UpdatedAt = time.Now()
	return nil
}
