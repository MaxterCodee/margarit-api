package models

import (
	"time"

	"gorm.io/gorm"
)

type Estudiante struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	UserID          uint         `gorm:"not null" json:"user_id"`
	User            User         `gorm:"foreignKey:UserID" json:"user"`
	Matricula       string       `gorm:"unique;not null" json:"matricula"`
	Nacionalidad    string       `gorm:"not null" json:"nacionalidad"`
	FechaNacimiento time.Time    `gorm:"not null" json:"fecha_nacimiento"`
	EdoOrigen       string       `gorm:"not null" json:"edo_origen"`
	MpioOrigen      string       `gorm:"not null" json:"mpio_origen"`
	EdoCivil        string       `gorm:"not null" json:"edo_civil"`
	Telefono        string       `gorm:"not null" json:"telefono"`
	PlantelID       uint         `gorm:"not null" json:"plantel_id"`
	Plantel         Plantel      `gorm:"foreignKey:PlantelID" json:"plantel"`
	NivelEscolarID  uint         `gorm:"not null" json:"nivel_escolar_id"`
	NivelEscolar    NivelEscolar `gorm:"foreignKey:NivelEscolarID" json:"nivel_escolar"`
	GrupoID         uint         `gorm:"not null" json:"grupo_id"`
	Grupo           Grupo        `gorm:"foreignKey:GrupoID" json:"grupo"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func (e *Estudiante) BeforeCreate(tx *gorm.DB) error {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Estudiante) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}
