package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarCategoriasPermisosIniciales inserta los registros iniciales de categorías de permisos
func InsertarCategoriasPermisosIniciales() {
	categorias := []models.CategoriaPermiso{
		{Titulo: "Gestión de roles y permisos", Descripcion: "Permisos relacionados con la administración de roles y sus permisos asociados.", Icono: "security"},
	}

	for _, categoria := range categorias {
		var existing models.CategoriaPermiso
		result := database.DB.Where("titulo = ?", categoria.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&categoria).Error; err != nil {
				log.Printf("Error insertando categoría de permiso %s: %v", categoria.Titulo, err)
			} else {
				log.Printf("Categoría de permiso '%s' insertada exitosamente", categoria.Titulo)
			}
		} else {
			log.Printf("Categoría de permiso '%s' ya existe, omitiendo", categoria.Titulo)
		}
	}
}