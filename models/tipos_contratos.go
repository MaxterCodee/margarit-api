// tipos_contratos.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type TipoContrato struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *TipoContrato) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *TipoContrato) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
