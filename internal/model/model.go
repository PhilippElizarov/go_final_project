package model

const DbName string = "scheduler.db"
const TimeTemplate string = "20060102"

var DbFile string

type Scheduler struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
type Tasks struct {
	Tasks []interface{} `json:"tasks"`
}
