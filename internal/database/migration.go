package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"bot_money/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Migrate(conf *config.Configuration) error {
	if conf.MigrateToVersion == "" {
		return nil
	}

	migrationsPath := conf.MigrationLocation

	_, err := os.Stat(migrationsPath)
	if err != nil {
		log.Printf("Migrate: failed to open '%s' folder, trying to create...\n", migrationsPath)
		err = os.Mkdir(migrationsPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	urlParams := config.GetConfiguration()
	urlString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		urlParams.DatabaseUser,
		urlParams.DatabasePassword,
		urlParams.DatabaseHost,
		urlParams.DatabaseName,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		urlString)
	if err != nil {
		return err
	}

	dbVersion, err := strconv.Atoi(conf.MigrateToVersion)
	if err == nil {
		log.Printf("Migrate: starting migration to version %v\n", dbVersion)
		err = m.Migrate(uint(dbVersion))
		if err != nil {
			log.Printf("Migrate: failed migration to version %v\n", dbVersion)
			log.Printf("Migration table will be forcing to version %v\n You should clean your data base from wrong tables and then start server mith 'MIGRATE=latest' enviroment variable!", dbVersion)
			err = m.Force(dbVersion)
		}
	} else {
		log.Println("Migrate: starting migration to the latest version")
		err = m.Up()
	}
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Migrate: no changes found")
			return nil
		}
		log.Println("file://" + migrationsPath)
		return err
	}
	log.Println("Migrate: migrations are done successfuly")
	return nil
}
