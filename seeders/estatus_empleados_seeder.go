package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarEstatusEmpleadosIniciales inserta los registros iniciales de estatus de empleados
func InsertarEstatusEmpleadosIniciales() {
	estatus := []models.EstatusEmpleado{
		{Titulo: "Activo"},
		{Titulo: "Inactivo"},
		{Titulo: "Suspendido"},
		{Titulo: "Terminado"},
	}

	for _, est := range estatus {
		var existing models.EstatusEmpleado
		result := database.DB.Where("titulo = ?", est.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&est).Error; err != nil {
				log.Printf("Error insertando estatus empleado %s: %v", est.Titulo, err)
			} else {
				log.Printf("Estatus empleado '%s' insertado exitosamente", est.Titulo)
			}
		} else {
			log.Printf("Estatus empleado '%s' ya existe, omitiendo", est.Titulo)
		}
	}
}