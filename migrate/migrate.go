package main

import (
	"flag"
	"log"

	"api-margaritai/config"
	"api-margaritai/database"
	"api-margaritai/models"
)

// insertarGenerosIniciales inserta los registros iniciales de género
func insertarGenerosIniciales() {
	generos := []models.Genero{
		{Nombre: "Masculino"},
		{Nombre: "Femenino"},
		{Nombre: "No especificado"},
	}

	for _, genero := range generos {
		// Verificar si el género ya existe
		var existingGenero models.Genero
		result := database.DB.Where("nombre = ?", genero.Nombre).First(&existingGenero)

		if result.Error != nil {
			// Si no existe, crearlo
			if err := database.DB.Create(&genero).Error; err != nil {
				log.Printf("Error insertando género %s: %v", genero.Nombre, err)
			} else {
				log.Printf("Género '%s' insertado exitosamente", genero.Nombre)
			}
		} else {
			log.Printf("Género '%s' ya existe, omitiendo", genero.Nombre)
		}
	}
}

// insertarEstatusEmpleadosIniciales inserta los registros iniciales de estatus de empleados
func insertarEstatusEmpleadosIniciales() {
	estatus := []models.EstatusEmpleado{
		{Titulo: "Activo"},
		{Titulo: "Inactivo"},
		{Titulo: "Suspendido"},
		{Titulo: "Terminado"},
	}

	for _, est := range estatus {
		var existing models.EstatusEmpleado
		result := database.DB.Where("titulo = ?", est.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&est).Error; err != nil {
				log.Printf("Error insertando estatus empleado %s: %v", est.Titulo, err)
			} else {
				log.Printf("Estatus empleado '%s' insertado exitosamente", est.Titulo)
			}
		} else {
			log.Printf("Estatus empleado '%s' ya existe, omitiendo", est.Titulo)
		}
	}
}

// insertarEstatusLaboralesIniciales inserta los registros iniciales de estatus laborales
func insertarEstatusLaboralesIniciales() {
	estatus := []models.EstatusLaboral{
		{Titulo: "Contratado"},
		{Titulo: "Por horas"},
		{Titulo: "Temporal"},
		{Titulo: "Pasantía"},
	}

	for _, est := range estatus {
		var existing models.EstatusLaboral
		result := database.DB.Where("titulo = ?", est.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&est).Error; err != nil {
				log.Printf("Error insertando estatus laboral %s: %v", est.Titulo, err)
			} else {
				log.Printf("Estatus laboral '%s' insertado exitosamente", est.Titulo)
			}
		} else {
			log.Printf("Estatus laboral '%s' ya existe, omitiendo", est.Titulo)
		}
	}
}

// insertarGradosAcademicosIniciales inserta los registros iniciales de grados académicos
func insertarGradosAcademicosIniciales() {
	grados := []models.GradoAcademico{
		{Titulo: "Licenciatura"},
		{Titulo: "Maestría"},
		{Titulo: "Doctorado"},
		{Titulo: "Técnico"},
		{Titulo: "Bachillerato"},
	}

	for _, grado := range grados {
		var existing models.GradoAcademico
		result := database.DB.Where("titulo = ?", grado.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&grado).Error; err != nil {
				log.Printf("Error insertando grado académico %s: %v", grado.Titulo, err)
			} else {
				log.Printf("Grado académico '%s' insertado exitosamente", grado.Titulo)
			}
		} else {
			log.Printf("Grado académico '%s' ya existe, omitiendo", grado.Titulo)
		}
	}
}

// insertarTiposContratosIniciales inserta los registros iniciales de tipos de contratos
func insertarTiposContratosIniciales() {
	tipos := []models.TipoContrato{
		{Titulo: "Tiempo completo"},
		{Titulo: "Medio tiempo"},
		{Titulo: "Por horas"},
		{Titulo: "Temporal"},
		{Titulo: "Pasantía"},
	}

	for _, tipo := range tipos {
		var existing models.TipoContrato
		result := database.DB.Where("titulo = ?", tipo.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&tipo).Error; err != nil {
				log.Printf("Error insertando tipo contrato %s: %v", tipo.Titulo, err)
			} else {
				log.Printf("Tipo contrato '%s' insertado exitosamente", tipo.Titulo)
			}
		} else {
			log.Printf("Tipo contrato '%s' ya existe, omitiendo", tipo.Titulo)
		}
	}
}

// insertarPuestosIniciales inserta los registros iniciales de puestos
func insertarPuestosIniciales() {
	puestos := []models.Puesto{
		{Titulo: "Director", PagoXHr: 500.0},
		{Titulo: "Subdirector", PagoXHr: 400.0},
		{Titulo: "Coordinador", PagoXHr: 350.0},
		{Titulo: "Profesor", PagoXHr: 300.0},
		{Titulo: "Secretario", PagoXHr: 200.0},
		{Titulo: "Conserje", PagoXHr: 150.0},
	}

	for _, puesto := range puestos {
		var existing models.Puesto
		result := database.DB.Where("titulo = ?", puesto.Titulo).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&puesto).Error; err != nil {
				log.Printf("Error insertando puesto %s: %v", puesto.Titulo, err)
			} else {
				log.Printf("Puesto '%s' insertado exitosamente", puesto.Titulo)
			}
		} else {
			log.Printf("Puesto '%s' ya existe, omitiendo", puesto.Titulo)
		}
	}
}

// insertarRolesIniciales inserta los registros iniciales de roles
func insertarRolesIniciales() {
	roles := []models.Rol{
		{Nombre: "Administrador", Descripcion: "Acceso completo al sistema", ParaEstudiante: false, ParaPersonal: true},
		{Nombre: "Director", Descripcion: "Gestión de plantel", ParaEstudiante: false, ParaPersonal: true},
		{Nombre: "Profesor", Descripcion: "Gestión de grupos y estudiantes", ParaEstudiante: false, ParaPersonal: true},
		{Nombre: "Estudiante", Descripcion: "Acceso a información académica", ParaEstudiante: true, ParaPersonal: false},
		{Nombre: "Tutor", Descripcion: "Acceso a información del estudiante", ParaEstudiante: false, ParaPersonal: false},
	}

	for _, rol := range roles {
		var existing models.Rol
		result := database.DB.Where("nombre = ?", rol.Nombre).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&rol).Error; err != nil {
				log.Printf("Error insertando rol %s: %v", rol.Nombre, err)
			} else {
				log.Printf("Rol '%s' insertado exitosamente", rol.Nombre)
			}
		} else {
			log.Printf("Rol '%s' ya existe, omitiendo", rol.Nombre)
		}
	}
}

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
			&models.Session{}, // Añadido Session aquí
			&models.Genero{},
			&models.EstatusEmpleado{},
			&models.EstatusLaboral{},
			&models.Puesto{},
			&models.GradoAcademico{},
			&models.TipoContrato{},
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
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&usersTableExists)
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'generos')").Scan(&generosTableExists)
	database.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'roles')").Scan(&rolesTableExists)

	if !usersTableExists || !generosTableExists || !rolesTableExists || *fresh {
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
			// Tablas con dependencias
			&models.User{},
			&models.Session{},
			&models.Direccion{},
			&models.Plantel{},
			&models.NivelEscolar{},
			&models.Grupo{},
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
		insertarGenerosIniciales()
		insertarEstatusEmpleadosIniciales()
		insertarEstatusLaboralesIniciales()
		insertarGradosAcademicosIniciales()
		insertarTiposContratosIniciales()
		insertarPuestosIniciales()
		insertarRolesIniciales()
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
		insertarGenerosIniciales()
		insertarEstatusEmpleadosIniciales()
		insertarEstatusLaboralesIniciales()
		insertarGradosAcademicosIniciales()
		insertarTiposContratosIniciales()
		insertarPuestosIniciales()
		insertarRolesIniciales()
	}
}
