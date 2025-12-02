package gestioncatalogos

import (
	"net/http"
	"time"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// ObtenerGrados obtiene todos los grados registrados
func ObtenerGrados(c *gin.Context) {
	var grados []models.Grado
	if err := database.DB.Find(&grados).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los grados"})
		return
	}
	c.JSON(http.StatusOK, grados)
}

// InsertarGrado inserta un nuevo grado
func InsertarGrado(c *gin.Context) {
	var input struct {
		Titulo         string `json:"titulo" binding:"required"`
		Descripcion    string `json:"descripcion"`
		NivelEscolarID uint   `json:"nivel_escolar_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	grado := models.Grado{
		Titulo:         input.Titulo,
		Descripcion:    input.Descripcion,
		NivelEscolarID: input.NivelEscolarID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := database.DB.Create(&grado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el grado"})
		return
	}
	c.JSON(http.StatusCreated, grado)
}

// EditarGrado edita los datos de un grado existente
func EditarGrado(c *gin.Context) {
	id := c.Param("id")
	var grado models.Grado
	if err := database.DB.First(&grado, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grado no encontrado"})
		return
	}

	var input struct {
		Titulo         string `json:"titulo"`
		Descripcion    string `json:"descripcion"`
		NivelEscolarID uint   `json:"nivel_escolar_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	if input.Titulo != "" {
		grado.Titulo = input.Titulo
	}
	grado.Descripcion = input.Descripcion
	if input.NivelEscolarID != 0 {
		grado.NivelEscolarID = input.NivelEscolarID
	}
	grado.UpdatedAt = time.Now()

	if err := database.DB.Save(&grado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el grado"})
		return
	}
	c.JSON(http.StatusOK, grado)
}

// EliminarGrado elimina un grado solo si no tiene materias relacionadas
func EliminarGrado(c *gin.Context) {
	id := c.Param("id")
	var grado models.Grado
	if err := database.DB.First(&grado, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grado no encontrado"})
		return
	}

	// Verificar que no existan materias asociadas a este grado
	var materiasCount int64
	if err := database.DB.Model(&models.Materia{}).Where("grado_id = ?", grado.ID).Count(&materiasCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo validar materias relacionadas"})
		return
	}
	if materiasCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar el grado porque existen materias relacionadas"})
		return
	}

	if err := database.DB.Delete(&grado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el grado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Grado eliminado correctamente"})
}
