// routes/api.go
package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"api-margaritai/controllers"
	gestioncatalogos "api-margaritai/controllers/gestion_catalogos"
	gestionusuarios "api-margaritai/controllers/gestion_usuarios"
	"api-margaritai/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint de health check
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/validate-token", controllers.ValidateToken)
	}

	protected := api.Group("/protected")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{"message": "Acceso concedido", "user_id": userID})
		})
		// Agrega el endpoint de logout
		protected.POST("/logout", controllers.Logout)

		// Endpoints especiales de roles (para obtener por tipo)
		protected.GET("/roles/para_estudiante", controllers.ObtenerRolesEstudiante)
		protected.GET("/roles/para_personal", controllers.ObtenerRolesPersonal)
		protected.GET("/roles/para_tutor", controllers.ObtenerRolesTutor)

		// Endpoints para roles
		protected.GET("/roles", controllers.GetRoles)
		protected.POST("/roles", controllers.CreateRole)

		// Rutas específicas de roles (deben ir antes que las rutas con parámetros)
		protected.GET("/roles/:id/permisos", controllers.GetPermisosDeRol)
		protected.GET("/roles/:id/permisos_agrupados", controllers.GetPermisosByRolId)
		protected.GET("/roles/:id/permisos_con_asignacion", controllers.GetRolePermisosConEstadoAsignacion)

		// Rutas generales de roles (con parámetros)
		protected.GET("/roles/:id", controllers.GetRole)
		protected.PUT("/roles/:id", controllers.UpdateRole)
		protected.DELETE("/roles/:id", controllers.DeleteRole)

		// Endpoints para permisos
		protected.GET("/permisos", controllers.GetPermisos)
		protected.POST("/permisos", controllers.CreatePermiso)

		// Rutas específicas de permisos (deben ir antes que las rutas con parámetros)
		protected.GET("/permisos/:id/roles", controllers.GetRolesDePermiso)

		// Rutas generales de permisos (con parámetros)
		protected.GET("/permisos/:id", controllers.GetPermiso)
		protected.PUT("/permisos/:id", controllers.UpdatePermiso)
		protected.DELETE("/permisos/:id", controllers.DeletePermiso)

		//endpoint para categorias_permisos
		protected.GET("/categorias_permisos", controllers.GetCategoriasPermisos)
		protected.GET("/categorias_permisos/:id", controllers.GetCategoriaPermiso)
		protected.POST("/categorias_permisos", controllers.CreateCategoriaPermiso)
		protected.PUT("/categorias_permisos/:id", controllers.UpdateCategoriaPermiso)
		protected.DELETE("/categorias_permisos/:id", controllers.DeleteCategoriaPermiso)

		// Endpoints para role_tiene_permiso
		protected.GET("/roles_tienen_permisos", controllers.GetRolesTienenPermisos)
		protected.GET("/roles_tienen_permisos/:role_id/:permiso_id", controllers.GetRoleTienePermiso)
		protected.POST("/roles_tienen_permisos", controllers.CreateRoleTienePermiso)
		protected.DELETE("/roles_tienen_permisos/:role_id/:permiso_id", controllers.DeleteRoleTienePermiso)
		protected.POST("/roles/asignar_permisos", controllers.AsignarPermisosARol)
		protected.POST("/roles/desasignar_permisos", controllers.DesasignarPermisosARol)

		// ---------- Rutas de gestión de catálogos: Planteles --------------
		protected.GET("/planteles", gestioncatalogos.ObtenerPlanteles)       // Obtener todos los planteles
		protected.POST("/planteles", gestioncatalogos.CrearPlantel)          // Crear un nuevo plantel
		protected.PUT("/planteles/:id", gestioncatalogos.EditarPlantel)      // Editar un plantel existente
		protected.DELETE("/planteles/:id", gestioncatalogos.EliminarPlantel) // Eliminar un plantel si cumple las restricciones

		// ---------- Rutas de gestión de catálogos: Niveles Escolares --------------
		protected.GET("/niveles_escolares", gestioncatalogos.ObtenerNivelesEscolares)     // Obtener todos los niveles escolares o filtrados
		protected.POST("/niveles_escolares", gestioncatalogos.CrearNivelEscolar)          // Crear un nuevo nivel escolar
		protected.PUT("/niveles_escolares/:id", gestioncatalogos.EditarNivelEscolar)      // Editar un nivel escolar existente
		protected.DELETE("/niveles_escolares/:id", gestioncatalogos.EliminarNivelEscolar) // Eliminar un nivel escolar si cumple las restricciones

		// ---------- RUTAS DE GESTIÓN DE CATÁLOGOS: GRADOS --------------
		protected.GET("/grados", gestioncatalogos.ObtenerGrados)        // Obtener todos los grados registrados
		protected.POST("/grados", gestioncatalogos.InsertarGrado)       // Insertar un nuevo grado
		protected.PUT("/grados/:id", gestioncatalogos.EditarGrado)      // Editar un grado existente
		protected.DELETE("/grados/:id", gestioncatalogos.EliminarGrado) // Eliminar un grado solo si no tiene materias relacionadas

		// ---------- RUTAS DE GESTIÓN DE CATÁLOGOS: GRUPOS --------------
		protected.GET("/grupos", gestioncatalogos.ObtenerGrupos)        // Obtener todos los grupos
		protected.POST("/grupos", gestioncatalogos.InsertarGrupo)       // Crear un nuevo grupo
		protected.PUT("/grupos/:id", gestioncatalogos.EditarGrupo)      // Editar un grupo existente
		protected.DELETE("/grupos/:id", gestioncatalogos.EliminarGrupo) // Eliminar un grupo existente

		// ---------- RUTAS DE GESTIÓN DE CATÁLOGOS: GRADOS ACADÉMICOS --------------
		protected.GET("/grados_academicos", gestioncatalogos.ObtenerGradoAcademico)         // Obtener todos los grados académicos
		protected.POST("/grados_academicos", gestioncatalogos.InsertarGradoAcademico)       // Crear un nuevo grado académico
		protected.PUT("/grados_academicos/:id", gestioncatalogos.EditarGradoAcademico)      // Editar un grado académico existente
		protected.DELETE("/grados_academicos/:id", gestioncatalogos.EliminarGradoAcademico) // Eliminar un grado académico existente

		// ---------- RUTAS DE GESTIÓN DE CATÁLOGOS: ESTATUS LABORALES --------------
		protected.GET("/estatus_laborales", gestioncatalogos.ObtenerEstatusLaborales)         // Obtener todos los estatus laborales
		protected.POST("/estatus_laborales", gestioncatalogos.InsertarEstatusLaborales)       // Crear un nuevo estatus laboral
		protected.PUT("/estatus_laborales/:id", gestioncatalogos.EditarEstatusLaborales)      // Editar un estatus laboral existente
		protected.DELETE("/estatus_laborales/:id", gestioncatalogos.EliminarEstatusLaborales) // Eliminar un estatus laboral existente

		// ---------- RUTAS DE GESTIÓN DE CATÁLOGOS: ESTATUS EMPLEADOS --------------
		protected.GET("/estatus_empleados", gestioncatalogos.ObtenerEstatusEmpleados)        // Obtener todos los estatus de empleados
		protected.POST("/estatus_empleados", gestioncatalogos.InsertarEstatusEmpleado)       // Crear un nuevo estatus de empleado
		protected.PUT("/estatus_empleados/:id", gestioncatalogos.EditarEstatusEmpleado)      // Editar un estatus de empleado existente
		protected.DELETE("/estatus_empleados/:id", gestioncatalogos.EliminarEstatusEmpleado) // Eliminar un estatus de empleado existente

		// ---------- RUTAS DE GESTIÓN DE CATÁLOGOS: PUESTOS --------------
		protected.GET("/puestos", gestioncatalogos.ObtenerPuestos)        // Obtener todos los puestos
		protected.POST("/puestos", gestioncatalogos.InsertarPuesto)       // Crear un nuevo puesto
		protected.PUT("/puestos/:id", gestioncatalogos.EditarPuesto)      // Editar un puesto existente
		protected.DELETE("/puestos/:id", gestioncatalogos.EliminarPuesto) // Eliminar un puesto existente

		// ---------- RUTAS DE GESTIÓN DE USUARIOS: ESTUDIANTES --------------
		protected.GET("/estudiantes", gestionusuarios.ObtenerEstudiantes)        // Obtener todos los estudiantes con su usuario
		protected.POST("/estudiantes", gestionusuarios.InsertarEstudiante)       // Crear un estudiante (usuario + estudiante)
		protected.PUT("/estudiantes/:id", gestionusuarios.EditarEstudiante)      // Editar datos de un estudiante y su usuario
		protected.DELETE("/estudiantes/:id", gestionusuarios.EliminarEstudiante) // Eliminar a un estudiante y su usuario asociado

		// ---------- RUTAS DE GESTIÓN DE USUARIOS: PERSONAL --------------
		protected.GET("/personal", gestionusuarios.ObtenerPersonal)         // Obtener la lista de personal con su usuario asociado
		protected.POST("/personal", gestionusuarios.InsertarPersonal)       // Crear un nuevo personal y usuario asociado
		protected.PUT("/personal/:id", gestionusuarios.EditarPersonal)      // Editar datos de personal y su usuario
		protected.DELETE("/personal/:id", gestionusuarios.EliminarPersonal) // Eliminar un registro de personal y su usuario asociado

		// ---------- RUTAS DE GESTIÓN DE USUARIOS: TUTORES --------------
		protected.GET("/tutores", gestionusuarios.ObtenerTutores)       // Obtener la lista de tutores con su usuario asociado
		protected.POST("/tutores", gestionusuarios.InsertarTutor)       // Crear un tutor y su usuario asociado
		protected.PUT("/tutores/:id", gestionusuarios.EditarTutor)      // Editar los datos de un tutor y su usuario asociado
		protected.DELETE("/tutores/:id", gestionusuarios.EliminarTutor) // Eliminar un tutor y su usuario asociado
	}

	return r
}
