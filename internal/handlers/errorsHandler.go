package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blackenkeeper/go_final_project/internal/models"
)

func ErrorsHandler(w http.ResponseWriter, err error, answer models.AnswerHandler) {
	log.Println("Catch an error:", err)
	answer.Error = err.Error()
	bodyPage, _ := json.Marshal(answer)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bodyPage)
}
