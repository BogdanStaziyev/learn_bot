package config

import "os"

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	MigrateToVersion    string
	MigrationLocation   string
	FileStorageLocation string
	BotToken            string
}

func GetConfiguration() *Configuration {
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "migration"
	}

	migratioToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migratioToVersion = "latest"
	}

	return &Configuration{
		DatabaseName:        "bot",
		DatabaseHost:        "localhost:54322",
		DatabaseUser:        "postgres",
		DatabasePassword:    "password",
		MigrateToVersion:    migratioToVersion,
		MigrationLocation:   migrationLocation,
		FileStorageLocation: "file_storage",
		BotToken:            os.Getenv("BOT_TOKEN"),
	}
}
