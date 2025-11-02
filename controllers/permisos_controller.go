package controllers

import (
	"net/http"
	"strconv"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

type CreatePermisoInput struct {
	Titulo             string `json:"titulo" binding:"required"`
	Descripcion        string `json:"descripcion"`
	CategoriaPermisoID uint   `json:"categoria_permiso_id" binding:"required"`
}

type UpdatePermisoInput struct {
	Titulo             *string `json:"titulo"`
	Descripcion        *string `json:"descripcion"`
	CategoriaPermisoID *uint   `json:"categoria_permiso_id"`
}

// Crear un nuevo permiso
func CreatePermiso(c *gin.Context) {
	var input CreatePermisoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permiso := models.Permiso{
		Titulo:             input.Titulo,
		Descripcion:        input.Descripcion,
		CategoriaPermisoID: input.CategoriaPermisoID,
	}

	if err := database.DB.Create(&permiso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando permiso"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permiso creado exitosamente",
		"permiso": permiso,
	})
}

// Obtener todos los permisos
func GetPermisos(c *gin.Context) {
	var permisos []models.Permiso
	if err := database.DB.Preload("CategoriaPermiso").Find(&permisos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo permisos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Permisos obtenidos exitosamente",
		"permisos": permisos,
	})
}

// Obtener un permiso específico por ID
func GetPermiso(c *gin.Context) {
	var permiso models.Permiso
	id := c.Param("id")

	if err := database.DB.Preload("CategoriaPermiso").First(&permiso, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permiso no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permiso obtenido exitosamente",
		"permiso": permiso,
	})
}

// Actualizar un permiso existente
func UpdatePermiso(c *gin.Context) {
	var permiso models.Permiso
	id := c.Param("id")

	if err := database.DB.First(&permiso, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permiso no encontrado"})
		return
	}

	var input UpdatePermisoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Titulo != nil {
		permiso.Titulo = *input.Titulo
	}
	if input.Descripcion != nil {
		permiso.Descripcion = *input.Descripcion
	}
	if input.CategoriaPermisoID != nil {
		permiso.CategoriaPermisoID = *input.CategoriaPermisoID
	}

	if err := database.DB.Save(&permiso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permiso actualizado exitosamente",
		"permiso": permiso,
	})
}

// Eliminar un permiso (soft delete)
func DeletePermiso(c *gin.Context) {
	var permiso models.Permiso
	id := c.Param("id")

	if err := database.DB.First(&permiso, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permiso no encontrado"})
		return
	}

	if err := database.DB.Delete(&permiso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permiso eliminado exitosamente",
	})
}

// Estructura para representar un permiso con su estado de asignación
type PermisoConAsignacion struct {
	ID          uint   `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Asignado    bool   `json:"asignado"`
}

// Estructura para representar una categoría con sus permisos
type CategoriaConPermisos struct {
	ID          uint                 `json:"id"`
	Titulo      string               `json:"titulo"`
	Descripcion string               `json:"descripcion"`
	Icono       string               `json:"icono"`
	Permisos    []PermisoConAsignacion `json:"permisos"`
}

// GetPermisosConEstadoAsignacion obtiene todos los permisos del sistema agrupados por categoría y verifica cuáles están asignados a un rol específico
func GetPermisosConEstadoAsignacion(c *gin.Context) {
	roleID := c.Param("role_id")
	
	// Convertir roleID a uint para comparaciones
	roleIDUint, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}
	
	// Obtener todas las categorías de permisos
	var categorias []models.CategoriaPermiso
	if err := database.DB.Find(&categorias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo categorías de permisos"})
		return
	}
	
	// Obtener todos los permisos del sistema con sus categorías
	var permisos []models.Permiso
	if err := database.DB.Preload("CategoriaPermiso").Find(&permisos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo permisos"})
		return
	}
	
	// Obtener los permisos asignados al rol
	var permisosDelRol []models.RoleTienePermiso
	if err := database.DB.Where("role_id = ?", roleID).Find(&permisosDelRol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo permisos del rol"})
		return
	}
	
	// Crear un mapa para verificar rápidamente si un permiso está asignado al rol
	permisosAsignados := make(map[uint]bool)
	for _, rolPermiso := range permisosDelRol {
		permisosAsignados[rolPermiso.PermisoID] = true
	}
	
	// Crear un mapa para agrupar permisos por categoría
	permisosPorCategoria := make(map[uint][]PermisoConAsignacion)
	
	// Agrupar permisos por categoría
	for _, permiso := range permisos {
		permisoConAsignacion := PermisoConAsignacion{
			ID:          permiso.ID,
			Titulo:      permiso.Titulo,
			Descripcion: permiso.Descripcion,
			Asignado:    permisosAsignados[permiso.ID],
		}
		
		permisosPorCategoria[permiso.CategoriaPermisoID] = append(permisosPorCategoria[permiso.CategoriaPermisoID], permisoConAsignacion)
	}
	
	// Crear la lista de categorías con sus permisos
	var categoriasConPermisos []CategoriaConPermisos
	
	for _, categoria := range categorias {
		categoriasConPermisos = append(categoriasConPermisos, CategoriaConPermisos{
			ID:          categoria.ID,
			Titulo:      categoria.Titulo,
			Descripcion: categoria.Descripcion,
			Icono:       categoria.Icono,
			Permisos:    permisosPorCategoria[categoria.ID],
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Permisos agrupados por categoría con estado de asignación obtenidos exitosamente",
		"categorias": categoriasConPermisos,
	})
}
