package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Nombre      string      `gorm:"not null" json:"nombre"`
	ApellidoP   string      `gorm:"not null" json:"apellido_p"`
	ApellidoM   string      `gorm:"not null" json:"apellido_m"`
	Email       string      `gorm:"unique;not null" json:"email"`
	CURP        string      `gorm:"unique;not null" json:"curp"`
	Password    string      `gorm:"not null" json:"-"`
	FechaNac    time.Time   `gorm:"not null" json:"fecha_nac"`
	GeneroID    uint        `gorm:"not null" json:"genero_id"`
	Genero      Genero      `gorm:"foreignKey:GeneroID" json:"genero"`
	RolID       uint        `gorm:"not null" json:"rol_id"`                                                      // Debe hacer referencia a un rol existente para evitar error de FK
	Rol         Rol         `gorm:"foreignKey:RolID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"rol"` // FK explícito
	Direcciones []Direccion `gorm:"foreignKey:UserID" json:"direcciones"`
	Grupos      []Grupo     `gorm:"foreignKey:UserID" json:"grupos"`
	Planteles   []Plantel   `gorm:"foreignKey:UserID" json:"planteles"`
	Tutores     []Tutor     `gorm:"foreignKey:UserID" json:"tutores"`
	EsActivo    bool        `gorm:"not null;default:true" json:"es_activo"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// IMPORTANTE: Si se va a asignar un Rol al crear/actualizar un usuario, RolID debe corresponder a un registro existente en la tabla "roles".
// Si se inserta un valor inválido, la DB rechazará la operación por la restricción de clave foránea.

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
