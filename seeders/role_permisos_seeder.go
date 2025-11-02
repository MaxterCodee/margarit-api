package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// AsignarPermisosAdministrador asigna los permisos de gestión de roles al rol "Administrador"
func AsignarPermisosAdministrador() {
	var adminRole models.Rol
	result := database.DB.Where("nombre = ?", "Administrador").First(&adminRole)

	if result.Error != nil {
		log.Fatalf("Error: Rol 'Administrador' no encontrado: %v", result.Error)
	}

	permisosTitulos := []string{
		"Ver roles",
		"Crear roles",
		"Editar roles",
		"Eliminar roles",
	}

	var permisos []models.Permiso
	result = database.DB.Where("titulo IN ?", permisosTitulos).Find(&permisos)

	if result.Error != nil {
		log.Fatalf("Error: No se pudieron encontrar los permisos: %v", result.Error)
	}

	// Asignar permisos al rol de Administrador
	for _, permiso := range permisos {
		var roleTienePermiso models.RoleTienePermiso
		result := database.DB.Where("role_id = ? AND permiso_id = ?", adminRole.ID, permiso.ID).First(&roleTienePermiso)

		if result.Error != nil {
			// Si no existe la asignación, crearla
			if err := database.DB.Create(&models.RoleTienePermiso{RoleID: adminRole.ID, PermisoID: permiso.ID}).Error; err != nil {
				log.Printf("Error asignando permiso '%s' al rol 'Administrador': %v", permiso.Titulo, err)
			} else {
				log.Printf("Permiso '%s' asignado exitosamente al rol 'Administrador'", permiso.Titulo)
			}
		} else {
			log.Printf("Permiso '%s' ya asignado al rol 'Administrador', omitiendo", permiso.Titulo)
		}
	}
}