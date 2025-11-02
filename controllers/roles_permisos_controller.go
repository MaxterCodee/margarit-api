package controllers

import (
	"net/http"
	"strconv"

	"api-margaritai/database"
	"api-margaritai/models"

	"github.com/gin-gonic/gin"
)

// Obtener todas las relaciones rol-permiso
func GetRolesTienenPermisos(c *gin.Context) {
	var relaciones []models.RoleTienePermiso
	if err := database.DB.Preload("Rol").Preload("Permiso").Find(&relaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo las relaciones rol-permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Relaciones rol-permiso obtenidas exitosamente",
		"relaciones": relaciones,
	})
}

// Obtener una relación específica rol-permiso
func GetRoleTienePermiso(c *gin.Context) {
	roleID := c.Param("role_id")
	permisoID := c.Param("permiso_id")

	var relacion models.RoleTienePermiso
	if err := database.DB.Preload("Rol").Preload("Permiso").
		Where("role_id = ? AND permiso_id = ?", roleID, permisoID).
		First(&relacion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relación rol-permiso no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Relación rol-permiso obtenida exitosamente",
		"relacion": relacion,
	})
}

// Crear una nueva relación rol-permiso
type CreateRoleTienePermisoInput struct {
	RoleID    uint `json:"role_id" binding:"required"`
	PermisoID uint `json:"permiso_id" binding:"required"`
}

func CreateRoleTienePermiso(c *gin.Context) {
	var input CreateRoleTienePermisoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el RoleID existe
	var role models.Rol
	if err := database.DB.First(&role, input.RoleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role no encontrado"})
		return
	}

	// Verificar si el PermisoID existe
	var permiso models.Permiso
	if err := database.DB.First(&permiso, input.PermisoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permiso no encontrado"})
		return
	}

	relacion := models.RoleTienePermiso{
		RoleID:    input.RoleID,
		PermisoID: input.PermisoID,
	}

	if err := database.DB.Create(&relacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando la relación rol-permiso"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Relación rol-permiso creada exitosamente",
		"relacion": relacion,
	})
}

// Eliminar una relación rol-permiso
func DeleteRoleTienePermiso(c *gin.Context) {
	roleID := c.Param("role_id")
	permisoID := c.Param("permiso_id")

	var relacion models.RoleTienePermiso
	if err := database.DB.Where("role_id = ? AND permiso_id = ?", roleID, permisoID).First(&relacion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relación rol-permiso no encontrada"})
		return
	}

	if err := database.DB.Delete(&relacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando la relación rol-permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relación rol-permiso eliminada exitosamente",
	})
}

// Obtener todos los permisos de un rol
func GetPermisosDeRol(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	var relaciones []models.RoleTienePermiso
	if err := database.DB.Preload("Permiso.CategoriaPermiso").Where("role_id = ?", roleID).Find(&relaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo los permisos del rol"})
		return
	}

	var permisos []models.Permiso
	for _, rel := range relaciones {
		permisos = append(permisos, rel.Permiso)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Permisos del rol obtenidos exitosamente",
		"permisos": permisos,
	})
}

// Obtener todos los roles de un permiso
func GetRolesDePermiso(c *gin.Context) {
	permisoIDStr := c.Param("id")
	permisoID, err := strconv.ParseUint(permisoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de permiso inválido"})
		return
	}

	var relaciones []models.RoleTienePermiso
	if err := database.DB.Preload("Rol").Where("permiso_id = ?", permisoID).Find(&relaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo los roles del permiso"})
		return
	}

	var roles []models.Rol
	for _, rel := range relaciones {
		roles = append(roles, rel.Rol)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roles del permiso obtenidos exitosamente",
		"roles":   roles,
	})
}
