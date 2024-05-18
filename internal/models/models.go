package models

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type AnswerHandler struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}
