package handlers

import "github.com/blackenkeeper/go_final_project/internal/database"

type Handler struct {
	Storage database.Storage
}

func GetHandler(database *database.Storage) *Handler {
	return &Handler{Storage: *database}
}
