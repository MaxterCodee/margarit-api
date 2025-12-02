package gestioncatalogos

import (
	"net/http"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// obtenerEstatusLaborales: obtiene todos los estatus laborales
func ObtenerEstatusLaborales(c *gin.Context) {
	var estatus []models.EstatusLaboral
	if err := database.DB.Find(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los estatus laborales"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Estatus laborales obtenidos correctamente",
		"data":    estatus,
	})
}

// insertarEstatusLaborales: inserta un nuevo estatus laboral
func InsertarEstatusLaborales(c *gin.Context) {
	var input struct {
		Titulo string `json:"titulo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estatus := models.EstatusLaboral{
		Titulo: input.Titulo,
	}

	if err := database.DB.Create(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el estatus laboral"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Estatus laboral creado correctamente",
		"data":    estatus,
	})
}

// editarEstatusLaborales: edita un estatus laboral por ID
func EditarEstatusLaborales(c *gin.Context) {
	id := c.Param("id")
	var estatus models.EstatusLaboral

	if err := database.DB.First(&estatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estatus laboral no encontrado"})
		return
	}

	var input struct {
		Titulo string `json:"titulo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estatus.Titulo = input.Titulo

	if err := database.DB.Save(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estatus laboral"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estatus laboral actualizado correctamente",
		"data":    estatus,
	})
}

// eliminarEstatusLaborales: elimina un estatus laboral por ID
func EliminarEstatusLaborales(c *gin.Context) {
	id := c.Param("id")
	var estatus models.EstatusLaboral

	if err := database.DB.First(&estatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estatus laboral no encontrado"})
		return
	}

	if err := database.DB.Delete(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el estatus laboral"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estatus laboral eliminado correctamente",
	})
}
