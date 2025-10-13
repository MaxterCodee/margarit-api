package controllers

import (
	"net/http"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

type CreateCategoriaPermisoInput struct {
	Titulo      string `json:"titulo" binding:"required"`
	Descripcion string `json:"descripcion"`
}

type UpdateCategoriaPermisoInput struct {
	Titulo      *string `json:"titulo"`
	Descripcion *string `json:"descripcion"`
}

// Obtener todas las categorías de permisos
func GetCategoriasPermisos(c *gin.Context) {
	var categorias []models.CategoriaPermiso
	if err := database.DB.Find(&categorias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo categorías de permisos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Categorías de permisos obtenidas exitosamente",
		"categorias": categorias,
	})
}

// Obtener una categoría de permiso específica por ID
func GetCategoriaPermiso(c *gin.Context) {
	var categoria models.CategoriaPermiso
	id := c.Param("id")

	if err := database.DB.First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría de permiso no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Categoría de permiso obtenida exitosamente",
		"categoria": categoria,
	})
}

// Crear una nueva categoría de permiso
func CreateCategoriaPermiso(c *gin.Context) {
	var input CreateCategoriaPermisoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoria := models.CategoriaPermiso{
		Titulo:      input.Titulo,
		Descripcion: input.Descripcion,
	}

	if err := database.DB.Create(&categoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando la categoría de permiso"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Categoría de permiso creada exitosamente",
		"categoria": categoria,
	})
}

// Actualizar una categoría de permiso existente
func UpdateCategoriaPermiso(c *gin.Context) {
	var categoria models.CategoriaPermiso
	id := c.Param("id")

	if err := database.DB.First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría de permiso no encontrada"})
		return
	}

	var input UpdateCategoriaPermisoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Titulo != nil {
		categoria.Titulo = *input.Titulo
	}
	if input.Descripcion != nil {
		categoria.Descripcion = *input.Descripcion
	}

	if err := database.DB.Save(&categoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando la categoría de permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Categoría de permiso actualizada exitosamente",
		"categoria": categoria,
	})
}

// Eliminar una categoría de permiso (soft delete)
func DeleteCategoriaPermiso(c *gin.Context) {
	var categoria models.CategoriaPermiso
	id := c.Param("id")

	if err := database.DB.First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría de permiso no encontrada"})
		return
	}

	if err := database.DB.Delete(&categoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando la categoría de permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categoría de permiso eliminada exitosamente",
	})
}
