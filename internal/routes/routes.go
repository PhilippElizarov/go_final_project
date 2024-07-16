package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/PhilippElizarov/go_final_project/internal/database"
	"github.com/PhilippElizarov/go_final_project/internal/model"
	"github.com/PhilippElizarov/go_final_project/internal/nextdate"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))

	r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handleAddTask)
	r.Get("/api/tasks", handleGetTasks)
	r.Get("/api/task", handleGetTaskByID)
	r.Put("/api/task", handleUpdateTask)
	r.Post("/api/task/done", handleDoneTask)
	r.Delete("/api/task", handleDeleteTask)

	return r
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")

	err := database.TaskStorage.DeleteTask(id)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusOK)
}

func handleDoneTask(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")

	err := database.TaskStorage.DoneTask(id)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusOK)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var response model.Response
	var task model.Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		response.Error = "Задача не найдена"
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		response.Error = "Не указан заголовок задачи"
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow := time.Now().Format(model.TimeTemplate)

	if task.Date == "" {
		task.Date = dateNow
	}

	var date_ time.Time

	date_, err = time.Parse(model.TimeTemplate, task.Date)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow_, err := time.Parse(model.TimeTemplate, dateNow)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if date_.Compare(dateNow_) == -1 {
		if task.Repeat == "" {
			task.Date = dateNow_.Format(model.TimeTemplate)
		} else {
			task.Date, err = nextdate.NextDate(dateNow_, task.Date, task.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(w).Encode(&response)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	err = database.TaskStorage.UpdateTask(task)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusOK)
}

func handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var response model.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := database.TaskStorage.GetTaskByID(id)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&task)
	w.WriteHeader(http.StatusFound)
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var tasks model.Tasks

	search := r.URL.Query().Get("search")

	tasks, err := database.TaskStorage.GetTasks(search)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if tasks.Tasks == nil {
		tasks.Tasks = []interface{}{}
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer
	var response model.Response

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		response.Error = "Не указан заголовок задачи"
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow := time.Now().Format(model.TimeTemplate)

	if task.Date == "" {
		task.Date = dateNow
	}

	var date_ time.Time

	date_, err = time.Parse(model.TimeTemplate, task.Date)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow_, err := time.Parse(model.TimeTemplate, dateNow)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if date_.Compare(dateNow_) == -1 { // || date_.Compare(dateNow_) == 0 {
		if task.Repeat == "" {
			task.Date = dateNow_.Format(model.TimeTemplate)
		} else {
			task.Date, err = nextdate.NextDate(dateNow_, task.Date, task.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(w).Encode(&response)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	response, err = database.TaskStorage.AddTask(task)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusCreated)
}

func handleNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowDate, err := time.Parse(model.TimeTemplate, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	nextDate, err := nextdate.NextDate(nowDate, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}
