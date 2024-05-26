package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/blackenkeeper/go_final_project/internal/handlers"
)

var (
	webDir = "./web"
	wg     sync.WaitGroup
)

// Функция создаёт сервер и запускает сервер на заданном порте и регистрирует обработчики путей.
// Значение порта можно задавать в перменной окружения TODO_PORT
func SetupServer() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/signin", handlers.LoginHandler)
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", handlers.Auth(handlers.TaskHandler))
	mux.HandleFunc("/api/tasks", handlers.Auth(handlers.TasksHandler))
	mux.HandleFunc("/api/task/done", handlers.Auth(handlers.TaskDoneHandler))

	addr := fmt.Sprintf(":%s", setupPort())

	log.Println("Запуск сервера на порте", addr)

	server := &http.Server{Addr: addr, Handler: mux}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Серевер остановлен по причине:", err)
		}
	}()

	//Обработчик корректного завершения сервера
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Остановка сервера...")
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Println("Серевер остановлен по причине:", err)
	}
	wg.Wait()
	log.Println("Сервер остановлен.")
}

// Функция возвращат порт, на котором будет запускаться сервер. Значение может быть задано в переменной
// окружения TODO_PORT
func setupPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
