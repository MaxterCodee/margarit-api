package gestioncatalogos

import (
	"net/http"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// obtenerGradoAcademico: obtiene todos los grados académicos
func ObtenerGradoAcademico(c *gin.Context) {
	var grados []models.GradoAcademico
	if err := database.DB.Find(&grados).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los grados académicos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Grados académicos obtenidos correctamente",
		"data":    grados,
	})
}

// insertarGradoAcademico: inserta un nuevo grado académico
func InsertarGradoAcademico(c *gin.Context) {
	var input struct {
		Titulo string `json:"titulo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grado := models.GradoAcademico{
		Titulo: input.Titulo,
	}

	if err := database.DB.Create(&grado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el grado académico"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Grado académico creado correctamente",
		"data":    grado,
	})
}

// editarGradoAcademico: edita un grado académico por ID
func EditarGradoAcademico(c *gin.Context) {
	id := c.Param("id")
	var grado models.GradoAcademico

	if err := database.DB.First(&grado, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grado académico no encontrado"})
		return
	}

	var input struct {
		Titulo string `json:"titulo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grado.Titulo = input.Titulo

	if err := database.DB.Save(&grado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el grado académico"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Grado académico actualizado correctamente",
		"data":    grado,
	})
}

// eliminarGradoAcademico: elimina un grado académico por ID
func EliminarGradoAcademico(c *gin.Context) {
	id := c.Param("id")
	var grado models.GradoAcademico

	if err := database.DB.First(&grado, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grado académico no encontrado"})
		return
	}

	if err := database.DB.Delete(&grado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el grado académico"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Grado académico eliminado correctamente",
	})
}
