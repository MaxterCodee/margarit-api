package gestionusuarios

import (
	"api-margaritai/database"
	"api-margaritai/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ObtenerEstudiantes obtiene todos los estudiantes con su usuario relacionado
func ObtenerEstudiantes(c *gin.Context) {
	var estudiantes []models.Estudiante

	if err := database.DB.Preload("User").Find(&estudiantes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los estudiantes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Estudiantes obtenidos correctamente",
		"estudiantes": estudiantes,
	})
}

// InsertarEstudiante crea un usuario y un estudiante asociado con control avanzado de errores
func InsertarEstudiante(c *gin.Context) {
	type UserInput struct {
		Nombre    string `json:"nombre" binding:"required"`
		ApellidoP string `json:"apellido_p" binding:"required"`
		ApellidoM string `json:"apellido_m" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		CURP      string `json:"curp" binding:"required"`
		Password  string `json:"password" binding:"required"`
		FechaNac  string `json:"fecha_nac" binding:"required"` // Formato YYYY-MM-DD
		GeneroID  uint   `json:"genero_id" binding:"required"`
		RolID     uint   `json:"rol_id" binding:"required"`
	}

	type EstudianteInput struct {
		Matricula         string `json:"matricula" binding:"required"`
		Nacionalidad      string `json:"nacionalidad" binding:"required"`
		FechaNacimiento   string `json:"fecha_nacimiento" binding:"required"` // YYYY-MM-DD
		EdoOrigen         string `json:"edo_origen" binding:"required"`
		MpioOrigen        string `json:"mpio_origen" binding:"required"`
		EdoCivil          string `json:"edo_civil" binding:"required"`
		Telefono          string `json:"telefono" binding:"required"`
		PlantelID         uint   `json:"plantel_id" binding:"required"`
		NivelEscolarID    uint   `json:"nivel_escolar_id" binding:"required"`
		GrupoID           uint   `json:"grupo_id" binding:"required"`
		EnProcesoAdmision *bool  `json:"en_proceso_admision"`
	}

	var input struct {
		UserInput
		EstudianteInput
	}

	// Manejo detallado de errores de bind
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Estructura de datos inválida",
			"details": err.Error(),
		})
		return
	}

	// Check unicidad de email y curp con error específico
	var count int64
	if tx := database.DB.Model(&models.User{}).
		Where("email = ?", input.Email).
		Or("curp = ?", input.CURP).
		Count(&count); tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar unicidad de usuario",
			"details": tx.Error.Error(),
		})
		return
	}
	if count > 0 {
		// Determinar duplicado exacto
		var existingUser models.User
		e := database.DB.Where("email = ?", input.Email).First(&existingUser)
		if e.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya existe."})
			return
		}
		e = database.DB.Where("curp = ?", input.CURP).First(&existingUser)
		if e.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La CURP ya existe."})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "El email o CURP ya existe."})
		return
	}

	// Check unicidad matrícula
	if tx := database.DB.Model(&models.Estudiante{}).
		Where("matricula = ?", input.Matricula).
		Count(&count); tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar unicidad de matrícula",
			"details": tx.Error.Error(),
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La matrícula ya existe.",
		})
		return
	}

	// Parse fechas user y estudiante
	fechaNacUser, err := time.Parse("2006-01-02", input.UserInput.FechaNac)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Formato de fecha_nac inválido (user): utilice YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}
	fechaNacEstudiante, err := time.Parse("2006-01-02", input.EstudianteInput.FechaNacimiento)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Formato de fecha_nacimiento inválido (estudiante): utilice YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	// Crear usuario con password hasheado
	user := models.User{
		Nombre:    input.Nombre,
		ApellidoP: input.ApellidoP,
		ApellidoM: input.ApellidoM,
		Email:     input.Email,
		CURP:      input.CURP,
		FechaNac:  fechaNacUser,
		GeneroID:  input.GeneroID,
		RolID:     input.RolID,
		EsActivo:  true,
	}
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al hashear el password.",
			"details": err.Error(),
		})
		return
	}

	// Transaccion para atomicidad entre usuario y estudiante
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error inesperado al crear estudiante.",
				"details": r,
			})
		}
	}()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al guardar el usuario en base de datos.",
			"details": err.Error(),
		})
		return
	}

	enProceso := true
	if input.EnProcesoAdmision != nil {
		enProceso = *input.EnProcesoAdmision
	}
	est := models.Estudiante{
		UserID:            user.ID,
		Matricula:         input.Matricula,
		Nacionalidad:      input.Nacionalidad,
		FechaNacimiento:   fechaNacEstudiante,
		EdoOrigen:         input.EdoOrigen,
		MpioOrigen:        input.MpioOrigen,
		EdoCivil:          input.EdoCivil,
		Telefono:          input.Telefono,
		PlantelID:         input.PlantelID,
		NivelEscolarID:    input.NivelEscolarID,
		GrupoID:           input.GrupoID,
		EnProcesoAdmision: enProceso,
	}

	if err := tx.Create(&est).Error; err != nil {
		tx.Rollback()
		// Intentar limpiar el usuario insertado
		database.DB.Unscoped().Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al guardar al estudiante en base de datos.",
			"details": err.Error(),
		})
		return
	}

	// Pre-cargar datos y responder
	if err := tx.Preload("User").First(&est, est.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al recuperar datos del estudiante insertado.",
			"details": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al finalizar la transacción.",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Estudiante creado correctamente",
		"estudiante": est,
	})
}

// EditarEstudiante edita los datos del estudiante y su usuario
func EditarEstudiante(c *gin.Context) {
	id := c.Param("id")
	var estudiante models.Estudiante

	if err := database.DB.Preload("User").First(&estudiante, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado"})
		return
	}

	// ASUMIMOS exactamente el mismo input que al crear
	type UserInput struct {
		Nombre    string `json:"nombre"`
		ApellidoP string `json:"apellido_p"`
		ApellidoM string `json:"apellido_m"`
		Email     string `json:"email"`
		CURP      string `json:"curp"`
		Password  string `json:"password"`
		FechaNac  string `json:"fecha_nac"`
		GeneroID  *uint  `json:"genero_id"`
		RolID     *uint  `json:"rol_id"`
	}

	type EstudianteInput struct {
		Matricula         string `json:"matricula"`
		Nacionalidad      string `json:"nacionalidad"`
		FechaNacimiento   string `json:"fecha_nacimiento"` // YYYY-MM-DD
		EdoOrigen         string `json:"edo_origen"`
		MpioOrigen        string `json:"mpio_origen"`
		EdoCivil          string `json:"edo_civil"`
		Telefono          string `json:"telefono"`
		PlantelID         *uint  `json:"plantel_id"`
		NivelEscolarID    *uint  `json:"nivel_escolar_id"`
		GrupoID           *uint  `json:"grupo_id"`
		EnProcesoAdmision *bool  `json:"en_proceso_admision"`
	}

	var input struct {
		UserInput
		EstudianteInput
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estructura de datos inválida", "details": err.Error()})
		return
	}

	// Actualización selectiva de campos del usuario
	user := &estudiante.User
	if input.Nombre != "" {
		user.Nombre = input.Nombre
	}
	if input.ApellidoP != "" {
		user.ApellidoP = input.ApellidoP
	}
	if input.ApellidoM != "" {
		user.ApellidoM = input.ApellidoM
	}
	if input.Email != "" && input.Email != user.Email {
		var count int64
		database.DB.Model(&models.User{}).Where("email = ? AND id <> ?", input.Email, user.ID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya existe para otro usuario"})
			return
		}
		user.Email = input.Email
	}
	if input.CURP != "" && input.CURP != user.CURP {
		var count int64
		database.DB.Model(&models.User{}).Where("curp = ? AND id <> ?", input.CURP, user.ID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La CURP ya existe para otro usuario"})
			return
		}
		user.CURP = input.CURP
	}
	if input.FechaNac != "" {
		fecha, err := time.Parse("2006-01-02", input.FechaNac)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha_nac inválido (user)"})
			return
		}
		user.FechaNac = fecha
	}
	if input.Password != "" {
		if err := user.HashPassword(input.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hashear password"})
			return
		}
	}
	if input.GeneroID != nil {
		user.GeneroID = *input.GeneroID
	}
	if input.RolID != nil {
		user.RolID = *input.RolID
	}

	// Actualización selectiva del estudiante
	if input.Matricula != "" && input.Matricula != estudiante.Matricula {
		var count int64
		database.DB.Model(&models.Estudiante{}).Where("matricula = ? AND id <> ?", input.Matricula, estudiante.ID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La matrícula ya existe para otro estudiante"})
			return
		}
		estudiante.Matricula = input.Matricula
	}
	if input.Nacionalidad != "" {
		estudiante.Nacionalidad = input.Nacionalidad
	}
	if input.FechaNacimiento != "" {
		fecha, err := time.Parse("2006-01-02", input.FechaNacimiento)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha_nacimiento inválido (estudiante)"})
			return
		}
		estudiante.FechaNacimiento = fecha
	}
	if input.EdoOrigen != "" {
		estudiante.EdoOrigen = input.EdoOrigen
	}
	if input.MpioOrigen != "" {
		estudiante.MpioOrigen = input.MpioOrigen
	}
	if input.EdoCivil != "" {
		estudiante.EdoCivil = input.EdoCivil
	}
	if input.Telefono != "" {
		estudiante.Telefono = input.Telefono
	}
	if input.PlantelID != nil {
		estudiante.PlantelID = *input.PlantelID
	}
	if input.NivelEscolarID != nil {
		estudiante.NivelEscolarID = *input.NivelEscolarID
	}
	if input.GrupoID != nil {
		estudiante.GrupoID = *input.GrupoID
	}
	if input.EnProcesoAdmision != nil {
		estudiante.EnProcesoAdmision = *input.EnProcesoAdmision
	}

	// Guardar usuario y luego estudiante
	if err := database.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}
	if err := database.DB.Save(&estudiante).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estudiante"})
		return
	}

	// Responder con el estudiante actualizado
	database.DB.Preload("User").First(&estudiante, estudiante.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Estudiante actualizado correctamente",
		"estudiante": estudiante,
	})
}

// EliminarEstudiante elimina el estudiante y el usuario asociado
func EliminarEstudiante(c *gin.Context) {
	id := c.Param("id")
	var estudiante models.Estudiante

	// Primero, obtener el estudiante con su UserID
	if err := database.DB.First(&estudiante, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado"})
		return
	}

	userID := estudiante.UserID

	// Eliminar el estudiante
	if err := database.DB.Delete(&estudiante).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar estudiante"})
		return
	}

	// Eliminar el usuario asociado
	if err := database.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Estudiante eliminado, pero error al eliminar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estudiante y usuario asociados eliminados correctamente"})
}
