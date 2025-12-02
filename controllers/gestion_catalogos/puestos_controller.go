package gestioncatalogos

import (
	"net/http"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// obtenerPuestos: obtiene todos los puestos
func ObtenerPuestos(c *gin.Context) {
	var puestos []models.Puesto
	if err := database.DB.Find(&puestos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los puestos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Puestos obtenidos correctamente",
		"data":    puestos,
	})
}

// insertarPuesto: inserta un nuevo puesto
func InsertarPuesto(c *gin.Context) {
	var input struct {
		Titulo  string  `json:"titulo" binding:"required"`
		PagoXHr float64 `json:"pago_x_hr" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	puesto := models.Puesto{
		Titulo:  input.Titulo,
		PagoXHr: input.PagoXHr,
	}

	if err := database.DB.Create(&puesto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el puesto"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Puesto creado correctamente",
		"data":    puesto,
	})
}

// editarPuesto: edita un puesto por ID
func EditarPuesto(c *gin.Context) {
	id := c.Param("id")
	var puesto models.Puesto

	if err := database.DB.First(&puesto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Puesto no encontrado"})
		return
	}

	var input struct {
		Titulo  string  `json:"titulo" binding:"required"`
		PagoXHr float64 `json:"pago_x_hr" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	puesto.Titulo = input.Titulo
	puesto.PagoXHr = input.PagoXHr

	if err := database.DB.Save(&puesto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el puesto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Puesto actualizado correctamente",
		"data":    puesto,
	})
}

// eliminarPuesto: elimina un puesto por ID
func EliminarPuesto(c *gin.Context) {
	id := c.Param("id")
	var puesto models.Puesto

	if err := database.DB.First(&puesto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Puesto no encontrado"})
		return
	}

	if err := database.DB.Delete(&puesto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el puesto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Puesto eliminado correctamente",
	})
}
