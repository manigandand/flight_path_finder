package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"volumefi/api"
	"volumefi/config"
	appmiddleware "volumefi/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	adkapi "github.com/manigandand/adk/api"
	"github.com/rs/cors"
)

var (
	name    = "volume.finance"
	version = "1.0.0"
)

func main() {
	adkapi.InitService(name, version)

	router := chi.NewRouter()
	appcors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{
			"Origin", "Authorization", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Header", "Accept",
			"Content-Type", "X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin", "Origin",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// cross & loger middleware
	router.Use(appcors.Handler)
	router.Use(
		middleware.Logger,
		appmiddleware.Recoverer,
	)

	router.Route("/", api.Routes)

	interruptChan := make(chan os.Signal, 1)
	go func() {
		signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		// Block until we receive our signal.
		<-interruptChan

		log.Println("Shutting down db...")
		os.Exit(0)
	}()

	log.Println("Starting server on port:", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router); err != nil {
		log.Fatal(err)
	}
}
