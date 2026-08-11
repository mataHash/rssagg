package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load()
	var portString string = os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
	}
	fmt.Println("Port: ", portString)

	var router chi.Router = chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	var v1Router chi.Router = chi.NewRouter()

	v1Router.Get("/healthz", handler_readiness)
	v1Router.Get("/err", handler_err)

	router.Mount("/v1", v1Router)

	var server *http.Server = &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server Starting at port %s", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
