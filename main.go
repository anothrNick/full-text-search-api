package main

import (
	"log"
	"os"

	"github.com/anothrNick/full-text-search-api/database"
	"github.com/anothrNick/full-text-search-api/web"
	"github.com/gin-gonic/gin"
)

func main() {
	// connect to db
	db, err := database.NewPostgres(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PW"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB"),
	)

	if err != nil {
		log.Fatal(err)
	}

	// create HTTP handlers
	httpHandler := web.NewHandlers(db)

	// create router, set routes
	router := gin.Default()

	// serve HTTP routes
	web.SetRoutes(router, httpHandler)

	// run server
	router.Run(":5001")
}
