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

	// Estructura para agrupar permisos por categoría
	type PermisosPorCategoria struct {
		Categoria models.CategoriaPermiso `json:"categoria"`
		Permisos  []models.Permiso        `json:"permisos"`
	}

	permisosAgrupadosMap := make(map[uint]PermisosPorCategoria)
	for _, rel := range relaciones {
		categoriaID := rel.Permiso.CategoriaPermiso.ID
		if _, ok := permisosAgrupadosMap[categoriaID]; !ok {
			permisosAgrupadosMap[categoriaID] = PermisosPorCategoria{
				Categoria: rel.Permiso.CategoriaPermiso,
				Permisos:  []models.Permiso{},
			}
		}
		p := permisosAgrupadosMap[categoriaID]
		p.Permisos = append(p.Permisos, rel.Permiso)
		permisosAgrupadosMap[categoriaID] = p
	}

	var permisosAgrupados []PermisosPorCategoria
	for _, v := range permisosAgrupadosMap {
		permisosAgrupados = append(permisosAgrupados, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Permisos del rol obtenidos exitosamente",
		"permisos_agrupados": permisosAgrupados,
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

// Estructura para recibir múltiples permisos a asignar a un rol
type AsignarPermisosARolInput struct {
	RoleID                uint   `json:"role_id" binding:"required"`
	PermisosPorAsignar    []uint `json:"permisos_por_asignar"`
	PermisosPorDesasignar []uint `json:"permisos_por_desasignar"`
}

// Asignar múltiples permisos a un rol
func AsignarPermisosARol(c *gin.Context) {
	var input AsignarPermisosARolInput
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

	// Iniciar una transacción para asegurar que todas las operaciones se realicen correctamente
	tx := database.DB.Begin()

	// --- Procesar permisos a desasignar ---
	for _, permisoID := range input.PermisosPorDesasignar {
		var relacion models.RoleTienePermiso
		if err := tx.Where("role_id = ? AND permiso_id = ?", input.RoleID, permisoID).First(&relacion).Error; err != nil {
			// Si no se encuentra la relación, simplemente la ignoramos o registramos un aviso
			// Por ahora, la ignoraremos para permitir la desasignación de permisos que quizás ya no existan
			continue
		}

		if err := tx.Delete(&relacion).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error desasignando permiso con ID " + strconv.FormatUint(uint64(permisoID), 10) + " del rol"})
			return
		}
	}

	// --- Procesar permisos a asignar ---
	permisosAsignados := []models.RoleTienePermiso{}
	permisosYaAsignados := []uint{}

	for _, permisoID := range input.PermisosPorAsignar {
		// Verificar que el permiso exista
		var permiso models.Permiso
		if err := database.DB.First(&permiso, permisoID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Permiso con ID " + strconv.FormatUint(uint64(permisoID), 10) + " no encontrado"})
			return
		}

		var relacion models.RoleTienePermiso
		if err := tx.Where("role_id = ? AND permiso_id = ?", input.RoleID, permisoID).First(&relacion).Error; err == nil {
			// El permiso ya está asignado al rol
			permisosYaAsignados = append(permisosYaAsignados, permisoID)
			continue
		}

		// Crear la nueva relación
		relacion = models.RoleTienePermiso{
			RoleID:    input.RoleID,
			PermisoID: permisoID,
		}

		if err := tx.Create(&relacion).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error asignando permiso con ID " + strconv.FormatUint(uint64(permisoID), 10) + " al rol"})
			return
		}

		permisosAsignados = append(permisosAsignados, relacion)
	}

	// Confirmar la transacción
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":               "Operación de permisos completada exitosamente",
		"permisos_asignados":    permisosAsignados,
		"permisos_ya_asignados": permisosYaAsignados,
		"permisos_desasignados": input.PermisosPorDesasignar, // Se puede mejorar para mostrar solo los que realmente se desasignaron
	})
}

// Obtener todos los permisos con estado de asignación para un rol específico
func GetRolePermisosConEstadoAsignacion(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	// Verificar si el rol existe
	var role models.Rol
	if err := database.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado"})
		return
	}

	// Obtener todos los permisos
	var permisos []models.Permiso
	if err := database.DB.Preload("CategoriaPermiso").Find(&permisos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo los permisos"})
		return
	}

	// Obtener los permisos asignados al rol
	var relaciones []models.RoleTienePermiso
	if err := database.DB.Where("role_id = ?", roleID).Find(&relaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo las relaciones rol-permiso"})
		return
	}

	// Crear un mapa para verificar rápidamente si un permiso está asignado
	permisosAsignados := make(map[uint]bool)
	for _, rel := range relaciones {
		permisosAsignados[rel.PermisoID] = true
	}

	// Estructura para representar un permiso con su estado de asignación
	type PermisoConEstado struct {
		ID               uint                   `json:"id"`
		Nombre           string                 `json:"nombre"`
		Descripcion      string                 `json:"descripcion"`
		CategoriaID      uint                   `json:"categoria_id"`
		CategoriaPermiso models.CategoriaPermiso `json:"categoria_permiso"`
		Asignado         bool                   `json:"asignado"`
	}

	// Estructura para agrupar permisos por categoría
	type PermisosPorCategoria struct {
		Categoria models.CategoriaPermiso `json:"categoria"`
		Permisos  []PermisoConEstado      `json:"permisos"`
	}

	// Agrupar permisos por categoría y añadir estado de asignación
	permisosAgrupadosMap := make(map[uint]PermisosPorCategoria)
	for _, permiso := range permisos {
		categoriaID := permiso.CategoriaPermisoID
		if _, ok := permisosAgrupadosMap[categoriaID]; !ok {
			permisosAgrupadosMap[categoriaID] = PermisosPorCategoria{
				Categoria: permiso.CategoriaPermiso,
				Permisos:  []PermisoConEstado{},
			}
		}

		permisoConEstado := PermisoConEstado{
			ID:               permiso.ID,
			Nombre:           permiso.Titulo,
			Descripcion:      permiso.Descripcion,
			CategoriaID:      permiso.CategoriaPermisoID,
			CategoriaPermiso: permiso.CategoriaPermiso,
			Asignado:         permisosAsignados[permiso.ID],
		}

		p := permisosAgrupadosMap[categoriaID]
		p.Permisos = append(p.Permisos, permisoConEstado)
		permisosAgrupadosMap[categoriaID] = p
	}

	var permisosAgrupados []PermisosPorCategoria
	for _, v := range permisosAgrupadosMap {
		permisosAgrupados = append(permisosAgrupados, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Permisos con estado de asignación obtenidos exitosamente",
		"permisos_agrupados": permisosAgrupados,
	})
}
