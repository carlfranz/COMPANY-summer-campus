package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// This is the main application entry point
// it can be run with 'go run main.go'.
// to build the application we need to run
// the command 'go build'.
func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
