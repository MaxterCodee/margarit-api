package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarGradosAcademicosIniciales inserta los registros iniciales de grados académicos
func InsertarGradosAcademicosIniciales() {
	grados := []models.GradoAcademico{
		{Titulo: "Licenciatura"},
		{Titulo: "Maestría"},
		{Titulo: "Doctorado"},
		{Titulo: "Técnico"},
		{Titulo: "Bachillerato"},
	}

	for _, grado := range grados {
		var existing models.GradoAcademico
		result := database.DB.Where("titulo = ?", grado.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&grado).Error; err != nil {
				log.Printf("Error insertando grado académico %s: %v", grado.Titulo, err)
			} else {
				log.Printf("Grado académico '%s' insertado exitosamente", grado.Titulo)
			}
		} else {
			log.Printf("Grado académico '%s' ya existe, omitiendo", grado.Titulo)
		}
	}
}