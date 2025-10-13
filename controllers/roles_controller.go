// controllers/roles_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-margaritai/database"
	"api-margaritai/models"
)

type CreateRoleInput struct {
	Nombre         string `json:"nombre" binding:"required"`
	Descripcion    string `json:"descripcion"`
	ParaEstudiante bool   `json:"para_estudiante"`
	ParaPersonal   bool   `json:"para_personal"`
}

func CreateRole(c *gin.Context) {
	var input CreateRoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el rol ya existe
	var existingRole models.Rol
	if err := database.DB.Where("nombre = ?", input.Nombre).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "El rol ya existe. Por favor, elija un nombre diferente o verifique los roles existentes antes de intentar crear uno nuevo.",
			"status": http.StatusBadRequest,
		})
		return
	}

	// Crear el nuevo rol
	rol := models.Rol{
		Nombre:         input.Nombre,
		Descripcion:    input.Descripcion,
		ParaEstudiante: input.ParaEstudiante,
		ParaPersonal:   input.ParaPersonal,
	}

	if err := database.DB.Create(&rol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando rol", "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rol creado exitosamente",
		"rol": gin.H{
			"id":              rol.ID,
			"nombre":          rol.Nombre,
			"descripcion":     rol.Descripcion,
			"para_estudiante": rol.ParaEstudiante,
			"para_personal":   rol.ParaPersonal,
			"created_at":      rol.CreatedAt,
			"updated_at":      rol.UpdatedAt,
		},
	})
}

// GetRoles obtiene todos los roles
func GetRoles(c *gin.Context) {
	var roles []models.Rol

	if err := database.DB.Order("id desc").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo roles", "status": http.StatusInternalServerError})
		return
	}

	var rolesResponse []gin.H
	for _, rol := range roles {
		tipo := ""
		if rol.ParaEstudiante {
			tipo = "Para estudiante"
		} else if rol.ParaPersonal {
			tipo = "Para personal"
		} else {
			tipo = ""
		}

		rolesResponse = append(rolesResponse, gin.H{
			"id":              rol.ID,
			"nombre":          rol.Nombre,
			"descripcion":     rol.Descripcion,
			"para_estudiante": rol.ParaEstudiante,
			"para_personal":   rol.ParaPersonal,
			"created_at":      rol.CreatedAt,
			"updated_at":      rol.UpdatedAt,
			"tipo":            tipo,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roles obtenidos exitosamente",
		"roles":   rolesResponse,
	})
}

// GetRole obtiene un rol específico por ID
func GetRole(c *gin.Context) {
	var rol models.Rol
	rolID := c.Param("id")

	if err := database.DB.First(&rol, rolID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado", "status": http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rol obtenido exitosamente",
		"rol":     rol,
	})
}

type UpdateRoleInput struct {
	Nombre         *string `json:"nombre"`
	Descripcion    *string `json:"descripcion"`
	ParaEstudiante *bool   `json:"para_estudiante"`
	ParaPersonal   *bool   `json:"para_personal"`
}

// UpdateRole actualiza un rol existente
func UpdateRole(c *gin.Context) {
	var rol models.Rol
	rolID := c.Param("id")

	// Verificar si el rol existe
	if err := database.DB.First(&rol, rolID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado", "status": http.StatusNotFound})
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el nuevo nombre ya existe (si se está actualizando)
	if input.Nombre != nil && *input.Nombre != rol.Nombre {
		var existingRole models.Rol
		if err := database.DB.Where("nombre = ? AND id != ?", *input.Nombre, rolID).First(&existingRole).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Ya existe un rol con ese nombre. Por favor, elija un nombre diferente para el rol. No se permiten duplicados para mantener la integridad de los datos.",
				"status":  http.StatusBadRequest,
				"detalle": "El nombre proporcionado ya está siendo utilizado por otro rol en el sistema. Si necesita actualizar el nombre, asegúrese de que sea único.",
			})
			return
		}
	}

	// Actualizar campos si se proporcionan
	if input.Nombre != nil {
		rol.Nombre = *input.Nombre
	}
	if input.Descripcion != nil {
		rol.Descripcion = *input.Descripcion
	}
	if input.ParaEstudiante != nil {
		rol.ParaEstudiante = *input.ParaEstudiante
	}
	if input.ParaPersonal != nil {
		rol.ParaPersonal = *input.ParaPersonal
	}

	if err := database.DB.Save(&rol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando rol", "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rol actualizado exitosamente",
		"rol":     rol,
	})
}

// DeleteRole elimina un rol (soft delete)
func DeleteRole(c *gin.Context) {
	var rol models.Rol
	rolID := c.Param("id")

	// Verificar si el rol existe
	if err := database.DB.First(&rol, rolID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado", "status": http.StatusNotFound})
		return
	}

	// Realizar soft delete
	if err := database.DB.Delete(&rol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando rol", "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rol eliminado exitosamente",
	})
}
