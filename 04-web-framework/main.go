package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Contact struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Phone   string
	Address string
	Email   string
	Website string
	Notes   string
}

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// This command creates and keeps update the database table related to the
	// contact Entity.
	db.AutoMigrate(&Contact{})

	r := gin.Default()

	r.POST("/contacts", func(c *gin.Context) {
		var contact Contact
		if err := c.ShouldBindJSON(&contact); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		saveContact(db, &contact)
		c.JSON(http.StatusCreated, contact)
	})

	r.PUT("/contacts/:id", func(c *gin.Context) {
		var contact Contact
		if err := c.ShouldBindJSON(&contact); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		contactId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			panic("failed parse the contactId")
		}
		updatedContact := updateContact(db, uint(contactId), contact)
		c.JSON(http.StatusOK, updatedContact)
	})

	r.DELETE("/contacts/:id", func(c *gin.Context) {
		contactId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			panic("failed parse the contactId")
		}
		deleteContact(db, uint(contactId))
		c.JSON(http.StatusNoContent, "")
	})

	r.GET("/contacts/:id", func(c *gin.Context) {
		contactId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			panic("failed parse the contactId")
		}
		contact := readContactById(db, uint(contactId))
		c.JSON(http.StatusOK, contact)
	})

	r.GET("/contacts", func(c *gin.Context) {
		allContacts := readAllContacts(db)
		c.JSON(http.StatusOK, allContacts)
	})

	r.Run()
}

func deleteContact(db *gorm.DB, contactId uint) {
	result := db.Delete(Contact{ID: contactId})
	if result.Error != nil {
		panic(fmt.Sprintf("Cannot delete contact with id '%d'", contactId))
	}
}

func updateContact(db *gorm.DB, contactId uint, contact Contact) (c Contact) {

	result := db.Model(Contact{ID: contactId}).First(&c)
	if result.Error != nil {
		panic(fmt.Sprintf("Cannot retrieve contact with id '%d'", contactId))
	}

	c.Address = contact.Address
	c.Email = contact.Email
	c.Name = contact.Name
	c.Notes = contact.Notes
	c.Website = contact.Website
	c.Phone = contact.Phone

	result = db.Save(&c)
	if result.Error != nil {
		panic(fmt.Sprintf("Cannot update contact with id '%d'", contactId))
	}
	return
}

func readContactById(db *gorm.DB, contactId uint) (result Contact) {
	db.Model(Contact{ID: contactId}).First(&result)
	return
}

func readAllContacts(db *gorm.DB) []Contact {
	var contacts []Contact
	result := db.Find(&contacts)
	if result.Error != nil {
		panic("Cannot list contacts")
	}
	return contacts
}

func saveContact(db *gorm.DB, contact *Contact) {
	result := db.Create(&contact)
	if result.Error != nil {
		panic("Cannot insert contact")
	}
}
