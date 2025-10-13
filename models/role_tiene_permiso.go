package models

type RoleTienePermiso struct {
	RoleID    uint    `gorm:"primaryKey" json:"role_id"`
	Rol       Rol     `gorm:"foreignKey:RoleID" json:"rol"`
	PermisoID uint    `gorm:"primaryKey" json:"permiso_id"`
	Permiso   Permiso `gorm:"foreignKey:PermisoID" json:"permiso"`
}
