package gestioncatalogos

import (
	"net/http"
	"strconv"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// ObtenerGrupos maneja la consulta de todos los grupos
func ObtenerGrupos(c *gin.Context) {
	var grupos []models.Grupo
	if err := database.DB.Preload("User").Preload("NivelEscolar").Find(&grupos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los grupos", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grupos)
}

// InsertarGrupo maneja la creación de un nuevo grupo
func InsertarGrupo(c *gin.Context) {
	var input struct {
		Titulo         string `json:"titulo" binding:"required"`
		UserID         uint   `json:"user_id" binding:"required"`
		NivelEscolarID uint   `json:"nivel_escolar_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	grupo := models.Grupo{
		Titulo:         input.Titulo,
		UserID:         input.UserID,
		NivelEscolarID: input.NivelEscolarID,
	}

	if err := database.DB.Create(&grupo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el grupo", "details": err.Error()})
		return
	}

	// Preload relaciones para la respuesta
	database.DB.Preload("User").Preload("NivelEscolar").First(&grupo, grupo.ID)
	c.JSON(http.StatusCreated, grupo)
}

// EditarGrupo maneja la edición de un grupo existente
func EditarGrupo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var grupo models.Grupo
	if err := database.DB.First(&grupo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grupo no encontrado"})
		return
	}

	var input struct {
		Titulo         string `json:"titulo"`
		UserID         uint   `json:"user_id"`
		NivelEscolarID uint   `json:"nivel_escolar_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Solo actualizar campos provistos
	if input.Titulo != "" {
		grupo.Titulo = input.Titulo
	}
	if input.UserID != 0 {
		grupo.UserID = input.UserID
	}
	if input.NivelEscolarID != 0 {
		grupo.NivelEscolarID = input.NivelEscolarID
	}

	if err := database.DB.Save(&grupo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar grupo", "details": err.Error()})
		return
	}

	database.DB.Preload("User").Preload("NivelEscolar").First(&grupo, id)
	c.JSON(http.StatusOK, grupo)
}

// EliminarGrupo maneja la eliminación de un grupo existente
func EliminarGrupo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var grupo models.Grupo
	if err := database.DB.First(&grupo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grupo no encontrado"})
		return
	}

	if err := database.DB.Delete(&grupo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar grupo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Grupo eliminado exitosamente"})
}
