package gestionusuarios

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"api-margaritai/database"
	"api-margaritai/models"
)

// ObtenerPersonal: devuelve la lista de personal con su usuario asociado.
func ObtenerPersonal(c *gin.Context) {
	var personal []models.Personal
	result := database.DB.Preload("User").Preload("GradoAcademico").
		Preload("EstatusLaboral").Preload("Puesto").Preload("EstatusEmpleado").
		Find(&personal)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar personal", "details": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, personal)
}

// InsertarPersonal: crea personal y usuario asociado.
func InsertarPersonal(c *gin.Context) {
	// Utilizar map[string]interface{} para bindear, debido a la ambigüedad con los campos no-exportados (Password) y el binding de json anidados
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// Validar que el usuario esté presente y sea un map[string]interface{}
	userMap, ok := payload["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo user es obligatorio"})
		return
	}

	// Parsear password manualmente para evitar problemas de binding
	password, passOk := userMap["password"].(string)
	if !passOk || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere contraseña"})
		return
	}

	// Checar email y curp únicos
	email, _ := userMap["email"].(string)
	curp, _ := userMap["curp"].(string)
	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo electrónico ya está registrado"})
		return
	}
	database.DB.Model(&models.User{}).Where("curp = ?", curp).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La CURP ya está registrada"})
		return
	}

	// Construir el usuario manualmente para asegurar que Password está correctamente presente
	usr := models.User{}
	// Los tipos deben ser correctos para el modelo. Parsear y asignar cada campo.
	if nombre, ok := userMap["nombre"].(string); ok {
		usr.Nombre = nombre
	}
	if apellidoP, ok := userMap["apellido_p"].(string); ok {
		usr.ApellidoP = apellidoP
	}
	if apellidoM, ok := userMap["apellido_m"].(string); ok {
		usr.ApellidoM = apellidoM
	}
	if email != "" {
		usr.Email = email
	}
	if curp != "" {
		usr.CURP = curp
	}
	if fechaNacStr, ok := userMap["fecha_nac"].(string); ok && fechaNacStr != "" {
		// Parse fecha_nac, asume formato RFC3339
		t, err := time.Parse(time.RFC3339, fechaNacStr)
		if err == nil {
			usr.FechaNac = t
		}
	}
	if generoID, ok := userMap["genero_id"].(float64); ok {
		usr.GeneroID = uint(generoID)
	}
	if rolID, ok := userMap["rol_id"].(float64); ok {
		usr.RolID = uint(rolID)
	}
	if esActivo, ok := userMap["es_activo"].(bool); ok {
		usr.EsActivo = esActivo
	}
	// Hash de password
	if err := usr.HashPassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo procesar la contraseña", "details": err.Error()})
		return
	}
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = usr.CreatedAt

	if err := database.DB.Create(&usr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario", "details": err.Error()})
		return
	}

	personal := models.Personal{
		UserID: usr.ID,
		// Los siguientes campos se obtienen del payload root, convertir cada uno adecuado
	}

	if rfc, ok := payload["rfc"].(string); ok {
		personal.RFC = rfc
	}
	if numEmp, ok := payload["numero_empleado"].(string); ok {
		personal.NumeroEmpleado = numEmp
	}
	if tel1, ok := payload["telefono_1"].(string); ok {
		personal.Telefono1 = tel1
	}
	if tel2, ok := payload["telefono_2"].(string); ok {
		personal.Telefono2 = tel2
	}
	if carrera, ok := payload["carrera"].(string); ok {
		personal.Carrera = carrera
	}
	if esProfesor, ok := payload["es_profesor"].(bool); ok {
		personal.EsProfesor = esProfesor
	}
	if gradoAcademicoID, ok := payload["grado_academico_id"].(float64); ok {
		personal.GradoAcademicoID = uint(gradoAcademicoID)
	}
	if estatusLaboralID, ok := payload["estatus_laboral_id"].(float64); ok {
		personal.EstatusLaboralID = uint(estatusLaboralID)
	}
	if puestoID, ok := payload["puesto_id"].(float64); ok {
		personal.PuestoID = uint(puestoID)
	}
	if estatusEmpleadoID, ok := payload["estatus_empleado_id"].(float64); ok {
		personal.EstatusEmpleadoID = uint(estatusEmpleadoID)
	}
	personal.CreatedAt = time.Now()
	personal.UpdatedAt = time.Now()

	if err := database.DB.Create(&personal).Error; err != nil {
		// Rollback usuario si no se creó el personal
		database.DB.Delete(&models.User{}, usr.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear personal", "details": err.Error()})
		return
	}

	var personalCreado models.Personal
	database.DB.Preload("User").First(&personalCreado, personal.ID)
	c.JSON(http.StatusCreated, personalCreado)
}

// EditarPersonal: edita personal y su usuario
func EditarPersonal(c *gin.Context) {
	var personal models.Personal
	id := c.Param("id")
	if err := database.DB.Preload("User").First(&personal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal no encontrado"})
		return
	}

	var input struct {
		User              *models.User `json:"user"`
		RFC               string       `json:"rfc"`
		NumeroEmpleado    string       `json:"numero_empleado"`
		Telefono1         string       `json:"telefono_1"`
		Telefono2         string       `json:"telefono_2"`
		Carrera           string       `json:"carrera"`
		EsProfesor        bool         `json:"es_profesor"`
		GradoAcademicoID  uint         `json:"grado_academico_id"`
		EstatusLaboralID  uint         `json:"estatus_laboral_id"`
		PuestoID          uint         `json:"puesto_id"`
		EstatusEmpleadoID uint         `json:"estatus_empleado_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// Edita datos del usuario
	if input.User != nil {
		user := personal.User

		if input.User.Nombre != "" {
			user.Nombre = input.User.Nombre
		}
		if input.User.ApellidoP != "" {
			user.ApellidoP = input.User.ApellidoP
		}
		if input.User.ApellidoM != "" {
			user.ApellidoM = input.User.ApellidoM
		}
		if input.User.Email != "" && input.User.Email != user.Email {
			var count int64
			database.DB.Model(&models.User{}).Where("email = ? AND id <> ?", input.User.Email, user.ID).Count(&count)
			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "El correo electrónico ya está registrado por otro usuario"})
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
		if !input.User.FechaNac.IsZero() {
			user.FechaNac = input.User.FechaNac
		}
		if input.User.GeneroID != 0 {
			user.GeneroID = input.User.GeneroID
		}
		if input.User.RolID != 0 {
			user.RolID = input.User.RolID
		}
		if input.User.EsActivo != user.EsActivo {
			user.EsActivo = input.User.EsActivo
		}
		// Si hay nueva contraseña
		if input.User.Password != "" {
			if err := user.HashPassword(input.User.Password); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la contraseña"})
				return
			}
		}
		user.UpdatedAt = time.Now()
		if err := database.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando User", "details": err.Error()})
			return
		}
	}

	// Edita datos de Personal
	toUpdate := map[string]interface{}{}
	if input.RFC != "" {
		toUpdate["rfc"] = input.RFC
	}
	if input.NumeroEmpleado != "" {
		toUpdate["numero_empleado"] = input.NumeroEmpleado
	}
	if input.Telefono1 != "" {
		toUpdate["telefono1"] = input.Telefono1
	}
	if input.Telefono2 != "" {
		toUpdate["telefono2"] = input.Telefono2
	}
	if input.Carrera != "" {
		toUpdate["carrera"] = input.Carrera
	}
	toUpdate["es_profesor"] = input.EsProfesor
	if input.GradoAcademicoID != 0 {
		toUpdate["grado_academico_id"] = input.GradoAcademicoID
	}
	if input.EstatusLaboralID != 0 {
		toUpdate["estatus_laboral_id"] = input.EstatusLaboralID
	}
	if input.PuestoID != 0 {
		toUpdate["puesto_id"] = input.PuestoID
	}
	if input.EstatusEmpleadoID != 0 {
		toUpdate["estatus_empleado_id"] = input.EstatusEmpleadoID
	}
	toUpdate["updated_at"] = time.Now()
	if err := database.DB.Model(&personal).Updates(toUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando Personal", "details": err.Error()})
		return
	}

	var actualizado models.Personal
	database.DB.Preload("User").First(&actualizado, personal.ID)
	c.JSON(http.StatusOK, actualizado)
}

// EliminarPersonal: elimina personal y su usuario
func EliminarPersonal(c *gin.Context) {
	id := c.Param("id")
	var personal models.Personal
	if err := database.DB.First(&personal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal no encontrado"})
		return
	}

	// Eliminar el Personal primero para evitar errores de restricción de clave foránea
	if err := database.DB.Delete(&models.Personal{}, personal.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el registro personal"})
		return
	}
	// Ahora eliminar el usuario asociado
	if err := database.DB.Delete(&models.User{}, personal.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo eliminar el usuario asociado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "Personal y usuario asociado eliminados exitosamente"})
}
