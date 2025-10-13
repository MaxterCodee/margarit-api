package models

import (
	"time"

	"gorm.io/gorm"
)

type Permiso struct {
	ID                 uint             `gorm:"primaryKey" json:"id"`
	Titulo             string           `gorm:"not null" json:"titulo"`
	Descripcion        string           `gorm:"type:text" json:"descripcion"`
	CategoriaPermisoID uint             `gorm:"not null" json:"categoria_permiso_id"`
	CategoriaPermiso   CategoriaPermiso `gorm:"foreignKey:CategoriaPermisoID" json:"categoria_permiso"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (p *Permiso) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Permiso) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
