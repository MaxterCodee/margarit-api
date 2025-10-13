package models

import (
	"time"

	"gorm.io/gorm"
)

type Direccion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Estado    string    `gorm:"not null" json:"estado"`
	Municipio string    `gorm:"not null" json:"municipio"`
	CPostal   string    `gorm:"not null" json:"c_postal"`
	Localidad string    `gorm:"not null" json:"localidad"`
	Direccion string    `gorm:"not null" json:"direccion"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Direccion) BeforeCreate(tx *gorm.DB) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

func (d *Direccion) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
