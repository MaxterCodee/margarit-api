package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarTiposContratosIniciales inserta los registros iniciales de tipos de contratos
func InsertarTiposContratosIniciales() {
	tipos := []models.TipoContrato{
		{Titulo: "Tiempo completo"},
		{Titulo: "Medio tiempo"},
		{Titulo: "Por horas"},
		{Titulo: "Temporal"},
		{Titulo: "Pasantía"},
	}

	for _, tipo := range tipos {
		var existing models.TipoContrato
		result := database.DB.Where("titulo = ?", tipo.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&tipo).Error; err != nil {
				log.Printf("Error insertando tipo contrato %s: %v", tipo.Titulo, err)
			} else {
				log.Printf("Tipo contrato '%s' insertado exitosamente", tipo.Titulo)
			}
		} else {
			log.Printf("Tipo contrato '%s' ya existe, omitiendo", tipo.Titulo)
		}
	}
}