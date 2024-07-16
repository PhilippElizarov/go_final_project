package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/PhilippElizarov/go_final_project/internal/model"
	"github.com/PhilippElizarov/go_final_project/internal/nextdate"
)

const limit = 50

type TaskStore struct {
	Db *sql.DB
}

var TaskStorage *TaskStore

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

func (s TaskStore) DeleteTask(id string) error {
	task, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}

	_, err = s.Db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	return nil
}

func (s TaskStore) DoneTask(id string) error {
	task, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}

	dateNow := time.Now().Format(model.TimeTemplate)
	dateNow_, err := time.Parse(model.TimeTemplate, dateNow)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		_, err = s.Db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
		if err != nil {
			return err
		}
	} else {
		task.Date, err = nextdate.NextDate(dateNow_, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		_, err = s.Db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
			sql.Named("id", task.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s TaskStore) UpdateTask(task model.Task) error {
	_, err := s.GetTaskByID(task.ID)
	if err != nil {
		return err
	}

	_, err = s.Db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	return nil
}

func (s TaskStore) GetTaskByID(id string) (model.Task, error) {
	var task model.Task
	row := s.Db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (s TaskStore) AddTask(task model.Task) (model.Response, error) {
	var response model.Response
	res, err := s.Db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
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

func (s TaskStore) GetTasks(search string) (model.Tasks, error) {
	var task model.Task
	var tasks model.Tasks
	var rows *sql.Rows
	var err error

	if search == "" {
		rows, err = s.Db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
		if err != nil {
			return tasks, err
		}
	} else {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			search = `%` + search + `%`
			rows, err = s.Db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", limit))
			if err != nil {
				return tasks, err
			}
		} else {
			rows, err = s.Db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
				sql.Named("date", date.Format(model.TimeTemplate)),
				sql.Named("limit", limit))
			if err != nil {
				return tasks, err
			}
		}
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}
