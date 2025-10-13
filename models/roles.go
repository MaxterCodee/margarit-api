package models

import (
	"time"

	"gorm.io/gorm"
)

type Rol struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Nombre         string         `gorm:"not null" json:"nombre"`
	Descripcion    string         `gorm:"type:text" json:"descripcion"`
	ParaEstudiante bool           `gorm:"not null" json:"para_estudiante"`
	ParaPersonal   bool           `gorm:"not null" json:"para_personal"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Rol) TableName() string {
	return "roles"
}

func (r *Rol) BeforeCreate(tx *gorm.DB) error {
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return nil
}

func (r *Rol) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}
