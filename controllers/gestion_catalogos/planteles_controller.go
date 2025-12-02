package gestioncatalogos

import (
	"api-margaritai/database"
	"api-margaritai/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ObtenerPlanteles obtiene todos los planteles existentes
func ObtenerPlanteles(c *gin.Context) {
	var planteles []models.Plantel

	if err := database.DB.Preload("User").Find(&planteles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los planteles", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"planteles": planteles,
	})
}

// CrearPlantel crea un nuevo plantel
func CrearPlantel(c *gin.Context) {
	type PlantelInput struct {
		Nombre      string `json:"nombre" binding:"required"`
		Descripcion string `json:"descripcion"`
		Ubicacion   string `json:"ubicacion" binding:"required"`
		Telefono    string `json:"telefono" binding:"required"`
		Correo      string `json:"correo" binding:"required,email"`
		UserID      uint   `json:"user_id" binding:"required"`
	}

	var input PlantelInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	plantel := models.Plantel{
		Nombre:      input.Nombre,
		Descripcion: input.Descripcion,
		Ubicacion:   input.Ubicacion,
		Telefono:    input.Telefono,
		Correo:      input.Correo,
		UserID:      input.UserID,
	}

	if err := database.DB.Create(&plantel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el plantel", "details": err.Error()})
		return
	}

	// Preload para devolver info del usuario relacionado, si es necesario
	if err := database.DB.Preload("User").First(&plantel, plantel.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Plantel creado exitosamente, pero hubo un problema obteniendo la información ampliada",
			"plantel": plantel,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plantel creado exitosamente",
		"plantel": plantel,
	})
}

// EditarPlantel actualiza la información de un plantel existente
func EditarPlantel(c *gin.Context) {
	// El ID puede venir como parámetro de la ruta: /planteles/:id
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del plantel requerido"})
		return
	}

	type PlantelUpdateInput struct {
		Nombre      *string `json:"nombre"`
		Descripcion *string `json:"descripcion"`
		Ubicacion   *string `json:"ubicacion"`
		Telefono    *string `json:"telefono"`
		Correo      *string `json:"correo"`
		UserID      *uint   `json:"user_id"`
	}

	var input PlantelUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	var plantel models.Plantel
	if err := database.DB.First(&plantel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plantel no encontrado"})
		return
	}

	// Solo actualizar campos que vienen en el JSON (no nulos)
	if input.Nombre != nil {
		plantel.Nombre = *input.Nombre
	}
	if input.Descripcion != nil {
		plantel.Descripcion = *input.Descripcion
	}
	if input.Ubicacion != nil {
		plantel.Ubicacion = *input.Ubicacion
	}
	if input.Telefono != nil {
		plantel.Telefono = *input.Telefono
	}
	if input.Correo != nil {
		plantel.Correo = *input.Correo
	}
	if input.UserID != nil {
		plantel.UserID = *input.UserID
	}

	if err := database.DB.Save(&plantel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el plantel", "details": err.Error()})
		return
	}

	// Preload de usuario asociado actualizado
	if err := database.DB.Preload("User").First(&plantel, plantel.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Plantel editado, pero hubo un problema obteniendo la información ampliada",
			"plantel": plantel,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plantel actualizado correctamente",
		"plantel": plantel,
	})
}

// EliminarPlantel elimina un plantel si no tiene estudiantes ni niveles escolares asociados
func EliminarPlantel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del plantel requerido"})
		return
	}

	plantelID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verifica que exista el plantel
	var plantel models.Plantel
	if err := database.DB.First(&plantel, plantelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plantel no encontrado"})
		return
	}

	// Verifica si existen estudiantes asociados
	var countEstudiantes int64
	if err := database.DB.Model(&models.Estudiante{}).Where("plantel_id = ?", plantelID).Count(&countEstudiantes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo verificar estudiantes asociados", "details": err.Error()})
		return
	}
	if countEstudiantes > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar el plantel porque existen estudiantes asociados"})
		return
	}

	// Verifica si existen niveles escolares asociados
	var countNiveles int64
	if err := database.DB.Model(&models.NivelEscolar{}).Where("plantel_id = ?", plantelID).Count(&countNiveles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo verificar niveles escolares asociados", "details": err.Error()})
		return
	}
	if countNiveles > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar el plantel porque existen niveles escolares asociados"})
		return
	}

	// Ahora sí, eliminar el plantel
	if err := database.DB.Delete(&models.Plantel{}, plantelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el plantel", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plantel eliminado exitosamente",
		"plantel": plantel,
	})
}
