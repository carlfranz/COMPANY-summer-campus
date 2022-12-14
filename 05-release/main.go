package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Contact example
type Contact struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Phone   string
	Address string
	Email   string
	Website string
	Notes   string
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host  localhost:8080
// @schemes http

func initDB() *gorm.DB {
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
	return db
}

var db = initDB()

// This is the main application entry point
// it can be run with 'go run main.go'.
// to build the application we need to run
// the command 'go build'.
func main() {

	// This command creates and keeps update the database table related to the
	// contact Entity.
	db.AutoMigrate(&Contact{})

	r := gin.Default()
	r.Use(cors.Default())

	contacts := r.Group("/contacts")
	{
		contacts.POST("/", createContact)
		contacts.PUT(":id", updateContactById)
		contacts.DELETE(":id", deleteContactById)
		contacts.GET(":id", getContactById)
		contacts.GET("/", listContacts)
	}

	r.Run()
}

// CONTROLLERS
////////////////////////////////////////////////////////////////////////////////

// Create Contact godoc.
// @Summary      Create new idea.
// @Description  Creates a new contact
// @tags         Contact
// @Accept       json
// @Param        Body  body      Contact  true  "All the informations required to create a contact"
// @Success      201   {object}  Contact
// @Router       /contacts [post]
func createContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := saveContact(db, &contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, contact)
	}
}

// UpdateContact updates a a contact data.
// @Summary      Update contact.
// @Description  Update the contact informations
// @tags         Contact
// @Accept       json
// @Param        Body  body  Contact  true  "All the property of the contact"
// @Param 		 id    path int true "Contact ID"
// @Success      200
// @Router       /contacts/{id} [put]
func updateContactById(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contactId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	updatedContact, err := updateContact(db, uint(contactId), contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, updatedContact)
	}
}

// DeleteContact deletes a contact.
// @Summary      Request delete contact.
// @Description  Allows the deletion of a contact.
// @Param 		 id  path int true "Contact ID"
// @tags         Contact
// @Success      200
// @Router       /contacts/{id} [delete]
func deleteContactById(c *gin.Context) {
	contactId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = deleteContact(db, uint(contactId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, "")
}

// GetContact get all details about a contact.
// @Summary      Get contact details.
// @Description  Gets detailed info about a contact.
// @Param 		 id  path int true "Contact ID"
// @tags         Contact
// @Produce      json
// @Success      200  {object}  Contact
// @Router       /contacts/{id} [get]
func getContactById(c *gin.Context) {
	contactId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	contact, err := readContactById(db, uint(contactId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, contact)
	}
}

// GetAllContacts Get all contacts.
// @Summary      Get the Contacts.
// @Description  Returns all the contacts in the contact manager.
// @tags         Contact
// @Produce      json
// @Success      200  {object}  []Contact
// @Router       /contacts [get]
func listContacts(c *gin.Context) {
	allContacts, err := readAllContacts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, allContacts)
}

// DATABASE
func deleteContact(db *gorm.DB, contactId uint) error {
	result := db.Delete(Contact{}, Contact{ID: contactId})
	if result.RowsAffected != 1 {
		return fmt.Errorf("cannot delete contact with id '%d'", contactId)
	}
	return nil
}

func updateContact(db *gorm.DB, contactId uint, contact Contact) (c *Contact, err error) {

	result := db.Model(Contact{}).First(&c, Contact{ID: contactId})
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("cannot retrieve contact with id '%d'", contactId)
	}

	c.Address = contact.Address
	c.Email = contact.Email
	c.Name = contact.Name
	c.Notes = contact.Notes
	c.Website = contact.Website
	c.Phone = contact.Phone

	result = db.Save(&c)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot update contact with id '%d'", contactId)
	}
	return
}

func readContactById(db *gorm.DB, contactId uint) (contact *Contact, err error) {
	result := db.Model(Contact{}).First(&contact, Contact{ID: contactId})
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf(`no user found with id '%d'`, contactId)
	}
	return
}

func readAllContacts(db *gorm.DB) ([]Contact, error) {
	var contacts []Contact
	result := db.Find(&contacts)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot list contacts")
	}
	return contacts, nil
}

func saveContact(db *gorm.DB, contact *Contact) error {
	result := db.Create(&contact)
	if result.Error != nil {
		return fmt.Errorf(`error saving contact`)
	}
	return nil
}
