package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/handlers"
	log "github.com/sirupsen/logrus"
)

var (
	webDir = "./web"
	wg     sync.WaitGroup
)

// Функция создаёт сервер и запускает сервер на заданном порте и регистрирует обработчики путей.
// Значение порта можно задавать в переменной окружения TODO_PORT
func SetupServer() {
	mux := http.NewServeMux()

	db, err := database.NewDB()
	if err != nil {
		log.WithError(err).Error("Ошибка создания и подключения к БД")
		return
	}
	defer db.CloseDB()

	h := handlers.GetHandler(db)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/signin", h.LoginHandler)
	mux.HandleFunc("/api/nextdate", h.NextDateHandler)
	mux.HandleFunc("/api/task", h.Auth(h.TaskHandler))
	mux.HandleFunc("/api/tasks", h.Auth(h.TasksHandler))
	mux.HandleFunc("/api/task/done", h.Auth(h.TaskDoneHandler))

	addr := fmt.Sprintf(":%s", setupPort())

	log.Infof("Запуск сервера на порте %s", addr)

	server := &http.Server{Addr: addr, Handler: mux}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Серевер остановлен")
		}
	}()

	// Обработчик корректного завершения сервера
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Info("Остановка сервера...")
	if err := server.Shutdown(context.TODO()); err != nil {
		log.WithError(err).Error("Серевер остановлен с ошибкой")
	}
	wg.Wait()
	log.Info("Сервер остановлен.")
}

// Функция возвращает порт, на котором будет запускаться сервер. Значение может быть задано в переменной
// окружения TODO_PORT
func setupPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
