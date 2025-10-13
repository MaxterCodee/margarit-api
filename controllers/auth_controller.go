// controllers/auth_controller.go
package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-margaritai/database"
	"api-margaritai/middleware"
	"api-margaritai/models"
)

type RegisterInput struct {
	Nombre    string `json:"nombre" binding:"required"`
	ApellidoP string `json:"apellido_p" binding:"required"`
	ApellidoM string `json:"apellido_m" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	CURP      string `json:"curp" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	FechaNac  string `json:"fecha_nac" binding:"required"`
	GeneroID  uint   `json:"genero_id" binding:"required"`
	RolID     uint   `json:"rol_id" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el email ya existe
	var existingUserByEmail models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUserByEmail).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ya ha sido registrado", "status": http.StatusBadRequest})
		return
	}

	// Verificar si el CURP ya existe
	var existingUserByCURP models.User
	if err := database.DB.Where("curp = ?", input.CURP).First(&existingUserByCURP).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CURP ya ha sido registrado", "status": http.StatusBadRequest})
		return
	}

	// Verificar si el género existe
	var genero models.Genero
	if err := database.DB.Where("id = ?", input.GeneroID).First(&genero).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Género no encontrado", "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando género", "status": http.StatusInternalServerError})
		}
		return
	}

	// Verificar si el rol existe
	var rol models.Rol
	if err := database.DB.Where("id = ?", input.RolID).First(&rol).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rol no encontrado", "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando rol", "status": http.StatusInternalServerError})
		}
		return
	}

	// Parsear la fecha de nacimiento
	fechaNac, err := time.Parse("2006-01-02", input.FechaNac)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
		return
	}

	user := models.User{
		Nombre:    input.Nombre,
		ApellidoP: input.ApellidoP,
		ApellidoM: input.ApellidoM,
		Email:     input.Email,
		CURP:      input.CURP,
		FechaNac:  fechaNac,
		GeneroID:  input.GeneroID,
		RolID:     input.RolID,
	}

	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hashear contraseña", "status": http.StatusInternalServerError})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando usuario", "status": http.StatusInternalServerError})
		return
	}

	// Cargar la información del género y rol
	if err := database.DB.Preload("Genero").Preload("Rol").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cargando información del usuario", "status": http.StatusInternalServerError})
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando token", "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"nombre":     user.Nombre,
			"apellido_p": user.ApellidoP,
			"apellido_m": user.ApellidoM,
			"email":      user.Email,
			"curp":       user.CURP,
			"fecha_nac":  user.FechaNac.Format("2006-01-02"),
			"genero_id":  user.GeneroID,
			"genero": gin.H{
				"id":     user.Genero.ID,
				"nombre": user.Genero.Nombre,
			},
			"rol_id": user.RolID,
			"rol": gin.H{
				"id":              user.Rol.ID,
				"nombre":          user.Rol.Nombre,
				"descripcion":     user.Rol.Descripcion,
				"para_estudiante": user.Rol.ParaEstudiante,
				"para_personal":   user.Rol.ParaPersonal,
				"created_at":      user.Rol.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at":      user.Rol.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		},
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Verificar si el correo existe antes de intentar cargar el usuario
	var count int64
	if err := database.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos, revisar conexión", "status": 500})
		return
	}
	if count == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "El correo electrónico no existe", "status": 401})
		return
	}

	// Preload Rol y Genero
	if err := database.DB.Preload("Genero").Preload("Rol").Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos, revisar conexión", "status": 500})
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "La contraseña es incorrecta"})
		return
	}

	// Obtener los permisos del rol del usuario
	var permisos []models.Permiso
	err := database.DB.
		Joins("JOIN role_tiene_permisoS ON role_tiene_permisoS.permiso_id = permisos.id").
		Where("role_tiene_permisoS.role_id = ?", user.RolID).
		Preload("CategoriaPermiso").
		Find(&permisos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo permisos del rol"})
		return
	}

	// Construir la lista de permisos para la respuesta
	permisosResponse := make([]gin.H, 0, len(permisos))
	for _, permiso := range permisos {
		permisosResponse = append(permisosResponse, gin.H{
			"id":                   permiso.ID,
			"titulo":               permiso.Titulo,
			"descripcion":          permiso.Descripcion,
			"categoria_permiso_id": permiso.CategoriaPermisoID,
			"categoria_permiso": gin.H{
				"id": permiso.CategoriaPermiso.ID,
				// No hay campo "nombre" en CategoriaPermiso, así que solo devolvemos el ID
			},
			"created_at": permiso.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": permiso.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando token"})
		return
	}

	// Crear una sesión que expira en 2 horas
	expiraEn := time.Now().Add(2 * time.Hour)
	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiraEn,
	}
	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error creando sesión",
			"status": 500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Inicio de sesión exitoso",
		"token":      token,
		"expires_at": expiraEn.Format("2006-01-02 15:04:05"),
		"user": gin.H{
			"id":         user.ID,
			"nombre":     user.Nombre,
			"apellido_p": user.ApellidoP,
			"apellido_m": user.ApellidoM,
			"email":      user.Email,
			"curp":       user.CURP,
			"fecha_nac":  user.FechaNac.Format("2006-01-02"),
			"genero_id":  user.GeneroID,
			"genero": gin.H{
				"id":     user.Genero.ID,
				"nombre": user.Genero.Nombre,
			},
			"rol_id": user.RolID,
			"rol": gin.H{
				"id":              user.Rol.ID,
				"nombre":          user.Rol.Nombre,
				"descripcion":     user.Rol.Descripcion,
				"para_estudiante": user.Rol.ParaEstudiante,
				"para_personal":   user.Rol.ParaPersonal,
				"created_at":      user.Rol.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at":      user.Rol.UpdatedAt.Format("2006-01-02 15:04:05"),
				"permisos":        permisosResponse,
			},
		},
	})
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token no proporcionado"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Invalidar el token
	middleware.InvalidateToken(tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout exitoso",
		"status":  http.StatusOK,
	})
}

// Función para la ruta /validate-token
// Valida si el token es válido consultando la sesión en base de datos
func ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token no proporcionado"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	var session models.Session
	if err := database.DB.Where("token = ?", tokenString).First(&session).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Por favor inicia sesión"})
		return
	}

	// Validar expiración
	if session.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "La sesión ha expirado, por favor inicia sesión nuevamente.",
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Token válido",
		"status":     http.StatusOK,
		"expires_at": session.ExpiresAt.Format("2006-01-02 15:04:05"),
	})
}
