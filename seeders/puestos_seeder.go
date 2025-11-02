package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarPuestosIniciales inserta los registros iniciales de puestos
func InsertarPuestosIniciales() {
	puestos := []models.Puesto{
		{Titulo: "Director", PagoXHr: 500.0},
		{Titulo: "Subdirector", PagoXHr: 400.0},
		{Titulo: "Coordinador", PagoXHr: 350.0},
		{Titulo: "Profesor", PagoXHr: 300.0},
		{Titulo: "Secretario", PagoXHr: 200.0},
		{Titulo: "Conserje", PagoXHr: 150.0},
	}

	for _, puesto := range puestos {
		var existing models.Puesto
		result := database.DB.Where("titulo = ?", puesto.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&puesto).Error; err != nil {
				log.Printf("Error insertando puesto %s: %v", puesto.Titulo, err)
			} else {
				log.Printf("Puesto '%s' insertado exitosamente", puesto.Titulo)
			}
		} else {
			log.Printf("Puesto '%s' ya existe, omitiendo", puesto.Titulo)
		}
	}
}