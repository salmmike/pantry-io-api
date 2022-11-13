package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Can't connect to postgresql %s", err)
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/db", getPantry(db))
	router.POST("/db", postPantry(db))
	router.GET("/create", createDeviceEntry(db))

	router.Run(":" + port)
}
