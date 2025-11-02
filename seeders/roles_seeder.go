package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarRolesIniciales inserta los registros iniciales de roles
func InsertarRolesIniciales() {
	roles := []models.Rol{
		{Nombre: "Administrador", Descripcion: "Acceso completo al sistema", ParaEstudiante: false, ParaPersonal: true},
		{Nombre: "Director", Descripcion: "Gestión de plantel", ParaEstudiante: false, ParaPersonal: true},
		{Nombre: "Profesor", Descripcion: "Gestión de grupos y estudiantes", ParaEstudiante: false, ParaPersonal: true},
		{Nombre: "Estudiante", Descripcion: "Acceso a información académica", ParaEstudiante: true, ParaPersonal: false},
		{Nombre: "Tutor", Descripcion: "Acceso a información del estudiante", ParaEstudiante: false, ParaPersonal: false},
	}

	for _, rol := range roles {
		var existing models.Rol
		result := database.DB.Where("nombre = ?", rol.Nombre).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&rol).Error; err != nil {
				log.Printf("Error insertando rol %s: %v", rol.Nombre, err)
			} else {
				log.Printf("Rol '%s' insertado exitosamente", rol.Nombre)
			}
		} else {
			log.Printf("Rol '%s' ya existe, omitiendo", rol.Nombre)
		}
	}
}