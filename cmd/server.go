package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	webDir = "./web"
	wg     sync.WaitGroup
)

func setupServer() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	addr := fmt.Sprintf(":%s", setupPort())

	log.Println("Starting the server on port", addr)

	server := &http.Server{Addr: addr}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Server stopped with error:", err)
		}
	}()

	//Обработчик корректного завершения сервера
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Stopping the server...")
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Println("Error while stopping server:", err)
	}
	wg.Wait()
	log.Println("Server stopped.")
}

func setupPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
