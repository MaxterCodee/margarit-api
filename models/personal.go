package models

import (
	"time"

	"gorm.io/gorm"
)

type Personal struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	RFC               string          `gorm:"not null" json:"rfc"`
	NumeroEmpleado    string          `gorm:"not null" json:"numero_empleado"`
	Telefono1         string          `gorm:"not null" json:"telefono_1"`
	Telefono2         string          `json:"telefono_2"`
	Carrera           string          `gorm:"not null" json:"carrera"`
	EsProfesor        bool            `gorm:"not null" json:"es_profesor"`
	GradoAcademicoID  uint            `gorm:"not null" json:"grado_academico_id"`
	GradoAcademico    GradoAcademico  `gorm:"foreignKey:GradoAcademicoID" json:"grado_academico"`
	EstatusLaboralID  uint            `gorm:"not null" json:"estatus_laboral_id"`
	EstatusLaboral    EstatusLaboral  `gorm:"foreignKey:EstatusLaboralID" json:"estatus_laboral"`
	PuestoID          uint            `gorm:"not null" json:"puesto_id"`
	Puesto            Puesto          `gorm:"foreignKey:PuestoID" json:"puesto"`
	EstatusEmpleadoID uint            `gorm:"not null" json:"estatus_empleado_id"`
	EstatusEmpleado   EstatusEmpleado `gorm:"foreignKey:EstatusEmpleadoID" json:"estatus_empleado"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (p *Personal) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Personal) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
