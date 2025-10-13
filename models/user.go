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
	RolID       uint        `gorm:"not null" json:"rol_id"`
	Rol         Rol         `gorm:"foreignKey:RolID" json:"rol"`
	Direcciones []Direccion `gorm:"foreignKey:UserID" json:"direcciones"`
	Grupos      []Grupo     `gorm:"foreignKey:UserID" json:"grupos"`
	Planteles   []Plantel   `gorm:"foreignKey:UserID" json:"planteles"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// HashPassword hashes the provided password and stores it in the User struct.
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword compares the stored hashed password with a plaintext password.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields before inserting a new User.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// BeforeUpdate sets the UpdatedAt field before updating an existing User.
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
