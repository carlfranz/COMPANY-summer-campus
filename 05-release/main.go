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

var db *gorm.DB

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
	r.Use(cors.Default())

	contacts := r.Group("/contacts")
	{
		contacts.POST("/", createContact)
		contacts.PUT(":id", updateContactById)
		contacts.DELETE("/contacts/:id", deleteContactById)
		contacts.GET("/contacts/:id", getContactById)
		contacts.GET("/contacts", listContacts)
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
	saveContact(db, &contact)
	c.JSON(http.StatusCreated, contact)
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
		panic("failed parse the contactId")
	}
	updatedContact := updateContact(db, uint(contactId), contact)
	c.JSON(http.StatusOK, updatedContact)
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
		panic("failed parse the contactId")
	}
	deleteContact(db, uint(contactId))
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
		panic("failed parse the contactId")
	}
	contact := readContactById(db, uint(contactId))
	c.JSON(http.StatusOK, contact)
}

// GetAllContacts Get all contacts.
// @Summary      Get the Contacts.
// @Description  Returns all the contacts in the contact manager.
// @tags         Contact
// @Produce      json
// @Success      200  {object}  []Contact
// @Router       /contacts [get]
func listContacts(c *gin.Context) {
	allContacts := readAllContacts(db)
	c.JSON(http.StatusOK, allContacts)
}

// DATABASE
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
