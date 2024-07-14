package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/PhilippElizarov/go_final_project/internal/model"
	"github.com/PhilippElizarov/go_final_project/internal/nextdate"
)

func CreateTable(db *sql.DB) {
	createSchedulerTableSQL := `CREATE TABLE scheduler (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,	
		"date" CHAR(8), 	
		"title" VARCHAR(128),
		"comment" TEXT,
		"repeat" VARCHAR(128) NULL
	  );
	  CREATE INDEX scheduler_date ON scheduler (date);`

	statement, err := db.Prepare(createSchedulerTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func DeleteTask(id string) error {
	var scheduler model.Scheduler
	db, err := sql.Open("sqlite3", model.DbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", scheduler.ID))
	if err != nil {
		return err
	}

	return nil
}

func DoneTask(id string) error {
	var scheduler model.Scheduler
	db, err := sql.Open("sqlite3", model.DbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
	if err != nil {
		return err
	}

	dateNow := time.Now().Format(model.TimeTemplate)
	dateNow_, err := time.Parse(model.TimeTemplate, dateNow)
	if err != nil {
		return err
	}

	if scheduler.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", scheduler.ID))
		if err != nil {
			return err
		}
	} else {
		scheduler.Date, err = nextdate.NextDate(dateNow_, scheduler.Date, scheduler.Repeat)
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", scheduler.Date),
			sql.Named("title", scheduler.Title),
			sql.Named("comment", scheduler.Comment),
			sql.Named("repeat", scheduler.Repeat),
			sql.Named("id", scheduler.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateTask(scheduler model.Scheduler) error {
	db, err := sql.Open("sqlite3", model.DbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", scheduler.ID))
	err = row.Scan(&scheduler.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", scheduler.Date),
		sql.Named("title", scheduler.Title),
		sql.Named("comment", scheduler.Comment),
		sql.Named("repeat", scheduler.Repeat),
		sql.Named("id", scheduler.ID))
	if err != nil {
		return err
	}

	return nil
}

func GetTaskByID(id string) (model.Scheduler, error) {
	var scheduler model.Scheduler

	db, err := sql.Open("sqlite3", model.DbFile)
	if err != nil {
		return scheduler, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
	if err != nil {
		return scheduler, err
	}

	return scheduler, nil
}

func PostTask(scheduler model.Scheduler) (model.Response, error) {
	var response model.Response

	db, err := sql.Open("sqlite3", model.DbFile)
	if err != nil {
		return response, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", scheduler.Date),
		sql.Named("title", scheduler.Title),
		sql.Named("comment", scheduler.Comment),
		sql.Named("repeat", scheduler.Repeat))
	if err != nil {
		return response, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return response, err
	}

	response.Id = strconv.FormatInt(id, 10)

	return response, nil
}

func GetTask(search string) (model.Tasks, model.Response, error) {
	var scheduler model.Scheduler
	var response model.Response
	var tasks model.Tasks
	var rows *sql.Rows

	db, err := sql.Open("sqlite3", model.DbFile)
	if err != nil {
		return tasks, response, err
	}
	defer db.Close()

	limit := 50

	if search == "" {
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
		if err != nil {
			return tasks, response, err
		}
	} else {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			search = `%` + search + `%`
			rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", limit))
			if err != nil {
				return tasks, response, err
			}
		} else {
			rows, err = db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
				sql.Named("date", date.Format(model.TimeTemplate)),
				sql.Named("limit", limit))
			if err != nil {
				return tasks, response, err
			}
		}
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
		if err != nil {
			return tasks, response, err
		}
		tasks.Tasks = append(tasks.Tasks, scheduler)
	}

	return tasks, response, nil
}
