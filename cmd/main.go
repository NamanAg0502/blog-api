package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/namanag0502/blog-api/pkg/db"
	"github.com/namanag0502/blog-api/pkg/routes"
	"github.com/namanag0502/blog-api/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	h := db.Init()
	defer func() {
		if err := h.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	r := routes.NewRouter(&h)

	port := utils.GetPort()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r.Router(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
