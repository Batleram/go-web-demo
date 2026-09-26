package main

import (
	"demo/core"
	"demo/todolist"
	todolistDatabase "demo/todolist/database"
	"fmt"
	"os"
	"path/filepath"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	host := "0.0.0.0:3000"
	if len(os.Args) < 3 {
		fmt.Printf("Host and port not specified, using defaults\n")
	}else {
		host = fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	}

	exepath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get exepath: %s\n", err.Error())
		os.Exit(1)
	}
	exedir := filepath.Dir(exepath)
	db, err := setupDatabase(exedir + "/test.sqlite3")
	if err != nil {
		fmt.Printf("Error opening the database: %s\n", err.Error())
		os.Exit(1)
	}

	err = migrateDatabase(db)
	if err != nil {
		fmt.Printf("Error migrating the database: %s\n", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := fiber.New()

	app.Use(logger.New())
	if mountTodoList(app, db) != nil {
		fmt.Println("ERROR: Failed to mount todolist app... SKIPPING")
	}

	app.Mount("/", core.RegisterRoutes())


	app.Listen(host)

}

func mountTodoList(app *fiber.App, db *sql.DB) error {
	todoDb, err := todolistDatabase.CreateTodoDatabase(db)
	if err != nil {
		return err
	}

	todolistRouter, err := todolist.RegisterRoutes(&todoDb)
	if err != nil {
		return err
	}

	app.Mount("/todolist", todolistRouter)
	return nil
}

func setupDatabase(connectionUrl string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", connectionUrl)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDatabase(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "sqlite3", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil

}
