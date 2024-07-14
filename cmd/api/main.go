package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PhilippElizarov/go_final_project/internal/database"
	"github.com/PhilippElizarov/go_final_project/internal/model"
	"github.com/PhilippElizarov/go_final_project/internal/routes"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	var dir string

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir, exists := os.LookupEnv("TODO_DBFILE")
	if exists {
		appPath = dir
	}

	model.DbFile = filepath.Join(filepath.Dir(appPath), model.DbName)

	var install bool
	_, err = os.Stat(model.DbFile)
	if err != nil {
		install = true
	}

	if install {
		file, err := os.Create(model.DbFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()

		sqliteDatabase, _ := sql.Open("sqlite3", model.DbFile)
		defer sqliteDatabase.Close()
		database.CreateTable(sqliteDatabase)
	}

	router := routes.NewRouter()

	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		port = "7540"
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err.Error())
	}
}
