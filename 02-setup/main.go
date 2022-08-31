package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// This is the main application entry point
// it can be run with 'go run main.go'.
// to build the application we need to run
// the command 'go build'.
func main() {
	// Database drivers use a connection string in order to
	// create a pool of connections to a database. Hence we
	// describe all the variables in a "Database Source name"
	// variable. The ssl parameter allow us to connect to our
	// local instance without a SSL certificate
	dsn := `host=localhost 
			user=gorm 
			password=gorm 
			dbname=contacts 
			port=5432 
			sslmode=disable 
			TimeZone=Europe/Rome`

	r := gin.Default()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// this is the definition on an endpoint. Whenever we receive a request with
	// GET html verb on the path http://localhost:8080/ping
	r.GET("/ping", func(c *gin.Context) {
		// we try to access the database in order to
		// see if the connection is ok
		tx := db.Begin()
		// do nothing
		tx.Commit()

		// then return a result using JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// running the router the application doesn't end his execution it keeps
	// the running until we press CTRL+C or we send OS signals.
	r.Run()
}
