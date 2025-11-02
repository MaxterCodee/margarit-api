package seeders

import (
	"log"

	"api-margaritai/database"
	"api-margaritai/models"
)

// InsertarPermisosIniciales inserta los registros iniciales de permisos
func InsertarPermisosIniciales() {
	var categoriaPermiso models.CategoriaPermiso
	result := database.DB.Where("titulo = ?", "Gestión de roles y permisos").First(&categoriaPermiso)

	if result.Error != nil {
		log.Fatalf("Error: Categoría de permiso 'Gestión de roles y permisos' no encontrada: %v", result.Error)
	}

	permisos := []models.Permiso{
		{Titulo: "Ver roles", Descripcion: "Permite ver los roles del sistema", CategoriaPermisoID: categoriaPermiso.ID},
		{Titulo: "Crear roles", Descripcion: "Permite crear nuevos roles en el sistema", CategoriaPermisoID: categoriaPermiso.ID},
		{Titulo: "Editar roles", Descripcion: "Permite editar roles existentes en el sistema", CategoriaPermisoID: categoriaPermiso.ID},
		{Titulo: "Eliminar roles", Descripcion: "Permite eliminar roles del sistema", CategoriaPermisoID: categoriaPermiso.ID},
	}

	for _, permiso := range permisos {
		var existing models.Permiso
		result := database.DB.Where("titulo = ?", permiso.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&permiso).Error; err != nil {
				log.Printf("Error insertando permiso %s: %v", permiso.Titulo, err)
			} else {
				log.Printf("Permiso '%s' insertado exitosamente", permiso.Titulo)
			}
		} else {
			log.Printf("Permiso '%s' ya existe, omitiendo", permiso.Titulo)
		}
	}
}