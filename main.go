package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mata-codes/rssagg/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	var portString string = os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
	}
	fmt.Println("Port: ", portString)

	var dbURL string = os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cant connect to database: ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(db, 10, time.Minute)

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
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.handlerGetUser))
	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.handlerGetPostsUser))

	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feed_followID}", apiCfg.MiddlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	var server *http.Server = &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server Starting at port %s", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
