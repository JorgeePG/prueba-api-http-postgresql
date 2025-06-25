package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "api_db"
)

// RunMigrations ejecuta los scripts SQL para crear o actualizar la estructura de la base de datos
func RunMigrations() error {
	log.Info().Msg("Iniciando proceso de migraciones") // NUEVO

	// Construir la cadena de conexión
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Debug().Str("connStr", connStr).Msg("Cadena de conexión construida") // NUEVO

	// Abrir la conexión a la base de datos
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Err(err).Msg("Error al abrir la conexión a la base de datos")
		return fmt.Errorf("error al abrir la conexión a la base de datos: %v", err)
	}
	defer db.Close()

	// Probar la conexión
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar a la base de datos")
		return fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	log.Info().Msg("Conexión a la base de datos establecida correctamente")

	// Ejecutar los scripts de migración
	// Comprobar varias rutas relativas posibles para encontrar las migraciones
	possiblePaths := []string{
		"infraestructure/db/migrations",       // Si se ejecuta desde la raíz del proyecto
		"../infraestructure/db/migrations",    // Si se ejecuta desde cmd/
		"../../infraestructure/db/migrations", // Si se ejecuta desde otro subdirectorio
	}

	var migrationsPath string
	var pathExists bool

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			migrationsPath = path
			pathExists = true
			break
		}
	}

	if !pathExists {
		log.Error().Msg("No se pudo encontrar la carpeta de migraciones en las rutas relativas especificadas")
		return fmt.Errorf("no se pudo encontrar la carpeta de migraciones en ninguna de las rutas relativas")
	}

	log.Info().Msgf("Buscando migraciones en: %s", migrationsPath)

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error al leer la carpeta de migraciones")
		return fmt.Errorf("error al leer la carpeta de migraciones: %v", err)
	}

	log.Info().Int("count", len(files)).Msg("Archivos encontrados en carpeta de migraciones") // NUEVO

	migrationCount := 0 // NUEVO: contador
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationCount++ // NUEVO
			log.Info().Msgf("Ejecutando migración: %s", file.Name())

			filePath := filepath.Join(migrationsPath, file.Name())
			sqlScript, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatal().Err(err).Msgf("Error al leer el archivo de migración %s", file.Name())
				return fmt.Errorf("error al leer el archivo de migración %s: %v", file.Name(), err)
			}

			_, err = db.Exec(string(sqlScript))
			if err != nil {
				log.Fatal().Err(err).Msgf("Error al ejecutar la migración %s", file.Name())
				return fmt.Errorf("error al ejecutar la migración %s: %v", file.Name(), err)
			}

			log.Info().Msgf("Migración ejecutada correctamente: %s", file.Name())
		} else {
			log.Debug().Str("file", file.Name()).Msg("Archivo ignorado (no es .sql)") // NUEVO
		}
	}

	log.Info().
		Int("total_migrations", migrationCount).
		Msg("Todas las migraciones se ejecutaron correctamente") // MEJORADO
	return nil
}
