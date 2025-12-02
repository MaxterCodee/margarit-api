package gestioncatalogos

import (
	"net/http"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// obtenerEstatusEmpleados: obtiene todos los estatus de empleados
func ObtenerEstatusEmpleados(c *gin.Context) {
	var estatus []models.EstatusEmpleado
	if err := database.DB.Find(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los estatus de empleados"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Estatus de empleados obtenidos correctamente",
		"data":    estatus,
	})
}

// insertarEstatusEmpleado: inserta un nuevo estatus de empleado
func InsertarEstatusEmpleado(c *gin.Context) {
	var input struct {
		Titulo string `json:"titulo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estatus := models.EstatusEmpleado{
		Titulo: input.Titulo,
	}

	if err := database.DB.Create(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el estatus de empleado"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Estatus de empleado creado correctamente",
		"data":    estatus,
	})
}

// editarEstatusEmpleado: edita un estatus de empleado por ID
func EditarEstatusEmpleado(c *gin.Context) {
	id := c.Param("id")
	var estatus models.EstatusEmpleado

	if err := database.DB.First(&estatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estatus de empleado no encontrado"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estatus de empleado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estatus de empleado actualizado correctamente",
		"data":    estatus,
	})
}

// eliminarEstatusEmpleado: elimina un estatus de empleado por ID
func EliminarEstatusEmpleado(c *gin.Context) {
	id := c.Param("id")
	var estatus models.EstatusEmpleado

	if err := database.DB.First(&estatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estatus de empleado no encontrado"})
		return
	}

	if err := database.DB.Delete(&estatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el estatus de empleado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estatus de empleado eliminado correctamente",
	})
}
