package controllers

import (
	"net/http"

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
