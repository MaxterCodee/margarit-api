package models

import (
	"time"

	"gorm.io/gorm"
)

type Grado struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	Titulo         string       `gorm:"not null" json:"titulo"`
	Descripcion    string       `gorm:"type:text" json:"descripcion"`
	NivelEscolarID uint         `gorm:"not null" json:"nivel_escolar_id"`
	NivelEscolar   NivelEscolar `gorm:"foreignKey:NivelEscolarID" json:"nivel_escolar"`
	Materias       []Materia    `gorm:"foreignKey:GradoID" json:"materias"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func (g *Grado) BeforeCreate(tx *gorm.DB) error {
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Grado) BeforeUpdate(tx *gorm.DB) error {
	g.UpdatedAt = time.Now()
	return nil
}
