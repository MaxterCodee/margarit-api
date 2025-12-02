package gestioncatalogos

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-margaritai/database"
	"api-margaritai/models"
)

// ObtenerNivelesEscolares retorna todos los niveles escolares (posiblemente filtrados por plantel_id si es pasado como query param)
func ObtenerNivelesEscolares(c *gin.Context) {
	db := database.DB

	var niveles []models.NivelEscolar

	plantelIDParam := c.Query("plantel_id")
	query := db.Preload("Plantel").Order("id")
	if plantelIDParam != "" {
		plantelID, err := strconv.ParseUint(plantelIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "plantel_id inválido"})
			return
		}
		query = query.Where("plantel_id = ?", plantelID)
	}
	if err := query.Find(&niveles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los niveles escolares"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"niveles_escolares": niveles})
}

// CrearNivelEscolar crea un nuevo nivel escolar
func CrearNivelEscolar(c *gin.Context) {
	db := database.DB

	type NivelEscolarInput struct {
		Titulo      string  `json:"titulo" binding:"required"`
		Descripcion string  `json:"descripcion"`
		Mensualidad float64 `json:"mensualidad" binding:"required"`
		PlantelID   uint    `json:"plantel_id" binding:"required"`
	}

	var input NivelEscolarInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// Verifica que exista el plantel
	var plantel models.Plantel
	if err := db.First(&plantel, input.PlantelID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plantel no encontrado"})
		return
	}

	nivel := models.NivelEscolar{
		Titulo:      input.Titulo,
		Descripcion: input.Descripcion,
		Mensualidad: input.Mensualidad,
		PlantelID:   input.PlantelID,
	}

	if err := db.Create(&nivel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el nivel escolar", "details": err.Error()})
		return
	}

	// Preload del plantel
	if err := db.Preload("Plantel").First(&nivel, nivel.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":       "Nivel escolar creado. Hubo un error cargando la información ampliada.",
			"nivel_escolar": nivel,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Nivel escolar creado exitosamente",
		"nivel_escolar": nivel,
	})
}

// EditarNivelEscolar actualiza un nivel escolar existente
func EditarNivelEscolar(c *gin.Context) {
	db := database.DB

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del nivel escolar requerido"})
		return
	}

	type NivelEscolarUpdateInput struct {
		Titulo      *string  `json:"titulo"`
		Descripcion *string  `json:"descripcion"`
		Mensualidad *float64 `json:"mensualidad"`
		PlantelID   *uint    `json:"plantel_id"`
	}

	var input NivelEscolarUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	var nivel models.NivelEscolar
	if err := db.First(&nivel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nivel escolar no encontrado"})
		return
	}

	if input.Titulo != nil {
		nivel.Titulo = *input.Titulo
	}
	if input.Descripcion != nil {
		nivel.Descripcion = *input.Descripcion
	}
	if input.Mensualidad != nil {
		nivel.Mensualidad = *input.Mensualidad
	}
	if input.PlantelID != nil {
		// Verifica que exista el plantel nuevo
		var plantel models.Plantel
		if err := db.First(&plantel, *input.PlantelID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nuevo plantel no encontrado"})
			return
		}
		nivel.PlantelID = *input.PlantelID
	}

	if err := db.Save(&nivel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el nivel escolar", "details": err.Error()})
		return
	}

	// Preload del plantel
	if err := db.Preload("Plantel").First(&nivel, nivel.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":       "Nivel escolar editado, pero hubo un problema obteniendo la información ampliada",
			"nivel_escolar": nivel,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Nivel escolar actualizado correctamente",
		"nivel_escolar": nivel,
	})
}

// EliminarNivelEscolar elimina un nivel escolar (solo si no existen estudiantes asociados, si aplica)
func EliminarNivelEscolar(c *gin.Context) {
	db := database.DB

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del nivel escolar requerido"})
		return
	}

	nivelID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verifica que exista el nivel escolar
	var nivel models.NivelEscolar
	if err := db.First(&nivel, nivelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nivel escolar no encontrado"})
		return
	}

	// (Opcional): Verifica si existen estudiantes asociados a este nivel escolar
	// Si tienes un modelo Estudiante que tiene NivelEscolarID
	var countEstudiantes int64
	if err := db.Table("estudiantes").Where("nivel_escolar_id = ?", nivelID).Count(&countEstudiantes).Error; err == nil && countEstudiantes > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar el nivel escolar porque existen estudiantes asociados"})
		return
	}

	if err := db.Delete(&models.NivelEscolar{}, nivelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el nivel escolar", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nivel escolar eliminado exitosamente"})
}
