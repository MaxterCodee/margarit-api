package models

import (
	"time"

	"gorm.io/gorm"
)

type Contrato struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	PersonalID     uint         `gorm:"not null" json:"personal_id"`
	Personal       Personal     `gorm:"foreignKey:PersonalID" json:"personal"`
	TipoContratoID uint         `gorm:"not null" json:"tipo_contrato_id"`
	TipoContrato   TipoContrato `gorm:"foreignKey:TipoContratoID" json:"tipo_contrato"`
	FechaInicio    time.Time    `gorm:"not null" json:"fecha_inicio"`
	FechaFin       time.Time    `json:"fecha_fin"`
	SalarioInicial float64      `gorm:"not null" json:"salario_inicial"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func (c *Contrato) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Contrato) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
