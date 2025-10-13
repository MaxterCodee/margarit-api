package models

import (
	"time"

	"gorm.io/gorm"
)

type Grupo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	AulaID    uint      `gorm:"not null" json:"aula_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Aula      Aula      `gorm:"foreignKey:AulaID" json:"aula"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (g *Grupo) BeforeCreate(tx *gorm.DB) error {
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Grupo) BeforeUpdate(tx *gorm.DB) error {
	g.UpdatedAt = time.Now()
	return nil
}
