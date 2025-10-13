// routes/api.go
package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"api-margaritai/controllers"
	"api-margaritai/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		// Ruta para validar el token
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

		// Endpoints para roles
		protected.GET("/roles", controllers.GetRoles)    // Obtener todos los roles
		protected.POST("/roles", controllers.CreateRole) // Crear un nuevo rol

		// Rutas específicas de roles (deben ir antes que las rutas con parámetros)
		protected.GET("/roles/:id/permisos", controllers.GetPermisosDeRol) // Obtener todos los permisos de un rol

		// Rutas generales de roles (con parámetros)
		protected.GET("/roles/:id", controllers.GetRole)       // Obtener un rol específico
		protected.PUT("/roles/:id", controllers.UpdateRole)    // Actualizar un rol existente
		protected.DELETE("/roles/:id", controllers.DeleteRole) // Eliminar un rol

		// Endpoints para permisos
		protected.GET("/permisos", controllers.GetPermisos)    // Obtener todos los permisos
		protected.POST("/permisos", controllers.CreatePermiso) // Crear un nuevo permiso

		// Rutas específicas de permisos (deben ir antes que las rutas con parámetros)
		protected.GET("/permisos/:id/roles", controllers.GetRolesDePermiso) // Obtener todos los roles de un permiso

		// Rutas generales de permisos (con parámetros)
		protected.GET("/permisos/:id", controllers.GetPermiso)       // Obtener un permiso específico
		protected.PUT("/permisos/:id", controllers.UpdatePermiso)    // Actualizar un permiso existente
		protected.DELETE("/permisos/:id", controllers.DeletePermiso) // Eliminar un permiso

		//endpoint para categorias_permisos
		protected.GET("/categorias_permisos", controllers.GetCategoriasPermisos)         // Obtener todas las categorías de permisos
		protected.GET("/categorias_permisos/:id", controllers.GetCategoriaPermiso)       // Obtener una categoría de permiso específica
		protected.POST("/categorias_permisos", controllers.CreateCategoriaPermiso)       // Crear una nueva categoría de permiso
		protected.PUT("/categorias_permisos/:id", controllers.UpdateCategoriaPermiso)    // Actualizar una categoría de permiso existente
		protected.DELETE("/categorias_permisos/:id", controllers.DeleteCategoriaPermiso) // Eliminar una categoría de permiso

		// Endpoints para role_tiene_permiso
		protected.GET("/roles_tienen_permisos", controllers.GetRolesTienenPermisos)                         // Obtener todas las relaciones rol-permiso
		protected.GET("/roles_tienen_permisos/:role_id/:permiso_id", controllers.GetRoleTienePermiso)       // Obtener una relación específica rol-permiso
		protected.POST("/roles_tienen_permisos", controllers.CreateRoleTienePermiso)                        // Crear una nueva relación rol-permiso
		protected.DELETE("/roles_tienen_permisos/:role_id/:permiso_id", controllers.DeleteRoleTienePermiso) // Eliminar una relación rol-permiso

	}

	return r
}
