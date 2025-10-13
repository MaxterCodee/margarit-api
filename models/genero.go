package models

import (
	"time"

	"gorm.io/gorm"
)

type Genero struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"not null;unique" json:"nombre"`
	Users     []User    `gorm:"foreignKey:GeneroID" json:"users"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (g *Genero) BeforeCreate(tx *gorm.DB) error {
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Genero) BeforeUpdate(tx *gorm.DB) error {
	g.UpdatedAt = time.Now()
	return nil
}
