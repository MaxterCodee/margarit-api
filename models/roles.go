package models

import (
	"time"

	"gorm.io/gorm"
)

type Rol struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Nombre         string         `gorm:"not null" json:"nombre"`
	Descripcion    string         `gorm:"type:text" json:"descripcion"`
	Icono          string         `gorm:"type:varchar(255)" json:"icono"` // nuevo campo para el nombre del icono de Material Icons
	ParaEstudiante bool           `gorm:"not null" json:"para_estudiante"`
	ParaPersonal   bool           `gorm:"not null" json:"para_personal"`
	ParaTutor      bool           `gorm:"not null" json:"para_tutor"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Permisos       []Permiso      `gorm:"many2many:role_tiene_permiso;" json:"permisos"`
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
