package models

import (
	"time"

	"gorm.io/gorm"
)

type Tutor struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	UserID            uint              `gorm:"not null" json:"user_id"`
	User              User              `gorm:"foreignKey:UserID" json:"user"`
	Nombre            string            `gorm:"not null" json:"nombre"`
	Telefono          string            `gorm:"not null" json:"telefono"`
	Telefono2         string            `gorm:"not null" json:"telefono2"`
	EstudianteTutores []EstudianteTutor `gorm:"foreignKey:TutorID" json:"estudiante_tutores"` // Relación muchos-a-muchos explícita para posibles usos avanzados
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
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
