package models

type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type HandlerAnswer struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
