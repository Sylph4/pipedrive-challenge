package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	internalHttp "github.com/sylph4/pipedrive-challenge/internal/gist/http"
	"github.com/sylph4/pipedrive-challenge/internal/gist/repository"
	"github.com/sylph4/pipedrive-challenge/internal/gist/service"
	"github.com/sylph4/pipedrive-challenge/storage"
)

func main() {
	dbConn, err := storage.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gistRepository := repository.NewGistRepository(dbConn)
	userRepository := repository.NewUserRepository(dbConn)
	gistService := service.NewGistService(gistRepository, userRepository)
	gistHandler := internalHttp.NewGistHandler(gistService, userRepository)

	mux := http.NewServeMux()
	mux.HandleFunc("/create-user", gistHandler.CreateUser)
	mux.HandleFunc("/delete-user", gistHandler.DeleteUser)
	mux.HandleFunc("/users", gistHandler.GetUsers)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("API listening at port :8080")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("shutting down the http server")
		}
	}()

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
