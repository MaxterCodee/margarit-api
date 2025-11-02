package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarGenerosIniciales inserta los registros iniciales de género
func InsertarGenerosIniciales() {
	generos := []models.Genero{
		{Nombre: "Masculino"},
		{Nombre: "Femenino"},
		{Nombre: "No especificado"},
	}

	for _, genero := range generos {
		// Verificar si el género ya existe
		var existingGenero models.Genero
		result := database.DB.Where("nombre = ?", genero.Nombre).First(&existingGenero)

		if result.Error != nil {
			// Si no existe, crearlo
			if err := database.DB.Create(&genero).Error; err != nil {
				log.Printf("Error insertando género %s: %v", genero.Nombre, err)
			} else {
				log.Printf("Género '%s' insertado exitosamente", genero.Nombre)
			}
		} else {
			log.Printf("Género '%s' ya existe, omitiendo", genero.Nombre)
		}
	}
}