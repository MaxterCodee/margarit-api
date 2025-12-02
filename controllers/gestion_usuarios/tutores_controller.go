package gestionusuarios

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"api-margaritai/database"
	"api-margaritai/models"
)

// obtenerTutores: devuelve la lista de tutores con su usuario asociado.
func ObtenerTutores(c *gin.Context) {
	var tutores []models.Tutor
	result := database.DB.Preload("User").Find(&tutores)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar tutores", "details": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, tutores)
}

// insertarTutor: crea un tutor con su usuario asociado.
func InsertarTutor(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// Validar usuario presente y correcto
	userMap, ok := payload["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo user es obligatorio"})
		return
	}

	password, passOk := userMap["password"].(string)
	if !passOk || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere contraseña para el usuario del tutor"})
		return
	}

	email, _ := userMap["email"].(string)
	curp, _ := userMap["curp"].(string)
	if email != "" {
		var count int64
		database.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está registrado"})
			return
		}
	}
	if curp != "" {
		var count int64
		database.DB.Model(&models.User{}).Where("curp = ?", curp).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La CURP ya está registrada"})
			return
		}
	}

	var user models.User
	user.Nombre, _ = userMap["nombre"].(string)
	user.ApellidoP, _ = userMap["apellido_p"].(string)
	user.ApellidoM, _ = userMap["apellido_m"].(string)
	user.Email, _ = userMap["email"].(string)
	user.CURP, _ = userMap["curp"].(string)
	user.EsActivo = true
	// Parse FechaNac
	if fn, exists := userMap["fecha_nac"].(string); exists && fn != "" {
		if fecha, err := time.Parse("2006-01-02", fn); err == nil {
			user.FechaNac = fecha
		}
	}
	// GeneroID y RolID
	if generoID, ok := userMap["genero_id"].(float64); ok {
		user.GeneroID = uint(generoID)
	}
	if rolID, ok := userMap["rol_id"].(float64); ok {
		user.RolID = uint(rolID)
	}
	// Password - hash
	if err := user.HashPassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	// Guardar usuario primero
	if err := database.DB.Create(&user).Error; err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "violates foreign key constraint") && strings.Contains(errMsg, "fk_users_rol") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "No se pudo crear el usuario",
				"details": "El rol especificado no es válido o no existe en la base de datos.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario", "details": errMsg})
		return
	}

	// Crear Tutor con referencias
	var tutor models.Tutor
	tutor.UserID = user.ID

	if nombre, ok := payload["nombre"].(string); ok {
		tutor.Nombre = nombre
	}
	if telefono, ok := payload["telefono"].(string); ok {
		tutor.Telefono = telefono
	}
	if telefono2, ok := payload["telefono2"].(string); ok {
		tutor.Telefono2 = telefono2
	}

	if err := database.DB.Create(&tutor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el tutor", "details": err.Error()})
		return
	}

	database.DB.Preload("User").First(&tutor, tutor.ID)
	c.JSON(http.StatusCreated, tutor)
}

// editarTutor: edita un tutor (y su usuario correspondiente)
func EditarTutor(c *gin.Context) {
	id := c.Param("id")

	// Estructura para recibir datos anidados
	type UserInput struct {
		Nombre    string `json:"nombre"`
		ApellidoP string `json:"apellido_p"`
		ApellidoM string `json:"apellido_m"`
		Email     string `json:"email"`
		CURP      string `json:"curp"`
		Password  string `json:"password"`
		FechaNac  string `json:"fecha_nac"`
		GeneroID  uint   `json:"genero_id"`
		RolID     uint   `json:"rol_id"`
		EsActivo  *bool  `json:"es_activo"`
	}
	type TutorInput struct {
		Nombre    string    `json:"nombre"`
		Telefono  string    `json:"telefono"`
		Telefono2 string    `json:"telefono2"`
		User      UserInput `json:"user"`
	}

	var input TutorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	var tutor models.Tutor
	if err := database.DB.First(&tutor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tutor no encontrado"})
		return
	}

	// Actualizar datos de tutor
	tutorMap := map[string]interface{}{}
	if input.Nombre != "" {
		tutorMap["nombre"] = input.Nombre
	}
	if input.Telefono != "" {
		tutorMap["telefono"] = input.Telefono
	}
	if input.Telefono2 != "" {
		tutorMap["telefono2"] = input.Telefono2
	}
	tutorMap["updated_at"] = time.Now()
	if err := database.DB.Model(&tutor).Updates(tutorMap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando tutor", "details": err.Error()})
		return
	}

	// Actualizar usuario asociado
	var user models.User
	if err := database.DB.First(&user, tutor.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario asociado no encontrado"})
		return
	}

	// Checar email/curp únicos SOLO si cambian
	if input.User.Email != "" && input.User.Email != user.Email {
		var count int64
		database.DB.Model(&models.User{}).Where("email = ? AND id <> ?", input.User.Email, user.ID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está registrado por otro usuario"})
			return
		}
		user.Email = input.User.Email
	}
	if input.User.CURP != "" && input.User.CURP != user.CURP {
		var count int64
		database.DB.Model(&models.User{}).Where("curp = ? AND id <> ?", input.User.CURP, user.ID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La CURP ya está registrada por otro usuario"})
			return
		}
		user.CURP = input.User.CURP
	}
	if input.User.Nombre != "" {
		user.Nombre = input.User.Nombre
	}
	if input.User.ApellidoP != "" {
		user.ApellidoP = input.User.ApellidoP
	}
	if input.User.ApellidoM != "" {
		user.ApellidoM = input.User.ApellidoM
	}
	if input.User.FechaNac != "" {
		if fecha, err := time.Parse("2006-01-02", input.User.FechaNac); err == nil {
			user.FechaNac = fecha
		}
	}
	if input.User.GeneroID != 0 {
		user.GeneroID = input.User.GeneroID
	}
	if input.User.RolID != 0 {
		user.RolID = input.User.RolID
	}
	if input.User.EsActivo != nil {
		user.EsActivo = *input.User.EsActivo
	}
	// Password (si manda uno nuevo)
	if input.User.Password != "" {
		if err := user.HashPassword(input.User.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la contraseña"})
			return
		}
	}
	user.UpdatedAt = time.Now()
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando usuario asociado", "details": err.Error()})
		return
	}

	// Responder tutor actualizado con User
	var actualizado models.Tutor
	database.DB.Preload("User").First(&actualizado, tutor.ID)
	c.JSON(http.StatusOK, actualizado)
}

// eliminarTutor: elimina un tutor y su usuario asociado
func EliminarTutor(c *gin.Context) {
	id := c.Param("id")
	var tutor models.Tutor
	if err := database.DB.First(&tutor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tutor no encontrado"})
		return
	}

	// Eliminar tutor primero (para no violar FK)
	if err := database.DB.Delete(&models.Tutor{}, tutor.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el registro de tutor"})
		return
	}
	// Eliminar usuario asociado
	if err := database.DB.Delete(&models.User{}, tutor.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el usuario asociado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "Tutor y usuario asociado eliminados exitosamente"})
}
