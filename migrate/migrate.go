package main

import (
	"flag"
	"log"

	"api-margaritai/config"
	"api-margaritai/database"
	"api-margaritai/models"
	"api-margaritai/seeders"
)

func main() {
	// Definir flag para migrate fresh
	fresh := flag.Bool("fresh", false, "Eliminar todas las tablas y recrear desde cero")
	flag.Parse()

	config.LoadEnv()
	database.ConnectDB()

	if *fresh {
		log.Println("Ejecutando migrate fresh - eliminando todas las tablas...")

		// Eliminar todas las tablas en orden inverso (respetando dependencias)
		err := database.DB.Migrator().DropTable(
			&models.Tutor{},
			&models.Estudiante{},
			&models.Condicion{},
			&models.Contrato{},
			&models.Personal{},
			&models.RoleTienePermiso{},
			&models.Permiso{},
			&models.CategoriaPermiso{},
			&models.Rol{},
			&models.Grupo{},
			&models.Aula{},
			&models.NivelEscolar{},
			&models.Plantel{},
			&models.Direccion{},
			&models.User{},
			&models.Session{},
			&models.Genero{},
			&models.EstatusEmpleado{},
			&models.EstatusLaboral{},
			&models.Puesto{},
			&models.GradoAcademico{},
			&models.TipoContrato{},
			&models.Grado{},
			&models.Materia{},
		)
		if err != nil {
			log.Fatal("Error eliminando tablas: ", err)
		}

		log.Println("Tablas eliminadas exitosamente")
	}

	// Verificar si las tablas existen
	var usersTableExists bool
	var generosTableExists bool
	var rolesTableExists bool
	var gruposTableExists bool
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&usersTableExists)
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'generos')").Scan(&generosTableExists)
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'roles')").Scan(&rolesTableExists)
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'grupos')").Scan(&gruposTableExists)

	if !usersTableExists || !generosTableExists || !rolesTableExists || !gruposTableExists || *fresh {
		// Si las tablas no existen o estamos haciendo fresh, usar AutoMigrate
		log.Println("Creando tablas con AutoMigrate...")
		err := database.DB.AutoMigrate(
			// Tablas base (sin dependencias)
			&models.Genero{},
			&models.EstatusEmpleado{},
			&models.EstatusLaboral{},
			&models.Puesto{},
			&models.GradoAcademico{},
			&models.TipoContrato{},
			&models.CategoriaPermiso{},
			&models.Permiso{},
			&models.Rol{},
			&models.Aula{},
			&models.Grado{},
			&models.Materia{},
			// Tablas con dependencias
			&models.User{},
			&models.Session{},
			&models.Direccion{},
			&models.Plantel{},
			&models.NivelEscolar{},
			&models.Grupo{}, // <-- Asegurar que Grupo está incluido
			&models.Personal{},
			&models.Contrato{},
			&models.Condicion{},
			&models.Estudiante{},
			&models.Tutor{},
			&models.RoleTienePermiso{},
		)
		if err != nil {
			log.Fatal("Error migrating database: ", err)
		}
		log.Println("Database migrated successfully")

		// Insertar datos iniciales
		log.Println("Insertando datos iniciales...")
		seeders.InsertarGenerosIniciales()
		seeders.InsertarEstatusEmpleadosIniciales()
		seeders.InsertarEstatusLaboralesIniciales()
		seeders.InsertarGradosAcademicosIniciales()
		seeders.InsertarTiposContratosIniciales()
		seeders.InsertarPuestosIniciales()
		seeders.InsertarRolesIniciales()
		seeders.InsertarCategoriasPermisosIniciales()
		seeders.InsertarPermisosIniciales()
		seeders.AsignarPermisosAdministrador()
	} else {
		log.Println("Tabla users existe, ejecutando migración manual...")

		// Paso 1: Agregar columna curp como nullable si no existe
		var columnExists bool
		database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'curp')").Scan(&columnExists)

		if !columnExists {
			log.Println("Agregando columna curp como nullable...")
			err := database.DB.Exec("ALTER TABLE users ADD COLUMN curp VARCHAR(18)").Error
			if err != nil {
				log.Fatal("Error agregando columna curp: ", err)
			}
		}

		// Paso 2: Verificar si hay registros sin CURP
		var count int64
		database.DB.Model(&models.User{}).Where("curp IS NULL OR curp = ''").Count(&count)

		if count > 0 {
			log.Printf("Encontrados %d registros sin CURP. Eliminando registros sin CURP...", count)
			// Eliminar registros sin CURP
			result := database.DB.Where("curp IS NULL OR curp = ''").Delete(&models.User{})
			if result.Error != nil {
				log.Fatal("Error eliminando registros sin CURP: ", result.Error)
			}
			log.Printf("Eliminados %d registros sin CURP", result.RowsAffected)
		}

		// Paso 3: Agregar restricciones NOT NULL y UNIQUE a la columna curp
		log.Println("Agregando restricciones NOT NULL y UNIQUE a la columna curp...")
		err := database.DB.Exec("ALTER TABLE users ALTER COLUMN curp SET NOT NULL").Error
		if err != nil {
			log.Fatal("Error agregando restricción NOT NULL: ", err)
		}

		// Agregar índice único para curp
		err = database.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_curp ON users(curp)").Error
		if err != nil {
			log.Fatal("Error creando índice único para curp: ", err)
		}

		log.Println("Migración manual completada exitosamente")

		// Verificar e insertar datos iniciales si no existen
		log.Println("Verificando datos iniciales...")
		seeders.InsertarGenerosIniciales()
		seeders.InsertarEstatusEmpleadosIniciales()
		seeders.InsertarEstatusLaboralesIniciales()
		seeders.InsertarGradosAcademicosIniciales()
		seeders.InsertarTiposContratosIniciales()
		seeders.InsertarPuestosIniciales()
		seeders.InsertarRolesIniciales()
		seeders.InsertarCategoriasPermisosIniciales()
		seeders.InsertarPermisosIniciales()
		seeders.AsignarPermisosAdministrador()
	}
}
