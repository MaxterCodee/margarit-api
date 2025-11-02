package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarEstatusLaboralesIniciales inserta los registros iniciales de estatus laborales
func InsertarEstatusLaboralesIniciales() {
	estatus := []models.EstatusLaboral{
		{Titulo: "Contratado"},
		{Titulo: "Por horas"},
		{Titulo: "Temporal"},
		{Titulo: "Pasantía"},
	}

	for _, est := range estatus {
		var existing models.EstatusLaboral
		result := database.DB.Where("titulo = ?", est.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&est).Error; err != nil {
				log.Printf("Error insertando estatus laboral %s: %v", est.Titulo, err)
			} else {
				log.Printf("Estatus laboral '%s' insertado exitosamente", est.Titulo)
			}
		} else {
			log.Printf("Estatus laboral '%s' ya existe, omitiendo", est.Titulo)
		}
	}
}