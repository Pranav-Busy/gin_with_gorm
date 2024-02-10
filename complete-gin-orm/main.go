package main

import (
	"helper/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"helper/storage"
	"fmt"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(c *gin.Context) {

	book := Book{}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Request failed", "error": err.Error()})
       return 
	}

	err := r.DB.Create(&book).Error

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "could not create book", "error": err.Error()})
       return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "book has ben added "})
    
}
func (r *Repository) DeleteBook(c *gin.Context) {
	bookModel := models.Books{}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "id cannot be empty"})
        return ;
	}

	err := r.DB.Delete(bookModel, id).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not delete book ", "error": err.Error()})
         return ;
	}
	c.JSON(http.StatusOK, gin.H{"message": " book was deleted "})
    
}

func (r *Repository) GetBooks(c *gin.Context) {

	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not get books", "error": err.Error()})
        return ;
	}

	c.JSON(http.StatusOK, gin.H{"message": "books fetched sucessfully", "data": bookModels})
    
}
func (r *Repository) GetBookByID(c *gin.Context) {

	id := c.Param("id")
	bookModel := &models.Books{}
	if id == "" {
		c.JSON(http.StatusInternalServerError,gin.H{"message":"id cannot be nil"})	
		return ;
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"message": "could not get the book"})
		return 
	}
	c.JSON(http.StatusOK,gin.H{"message": "book id fetched successfully","data":bookModel})
	
}
func (r *Repository) SetupRoutes(app *gin.Engine) {

	api := app.Group("/api")
	api.POST("/create_books", r.CreateBook)
	api.DELETE("delete_book/:id", r.DeleteBook)
	api.GET("/get_books/:id", r.GetBookByID)
	api.GET("/books", r.GetBooks)
}
func main() {

	err := godotenv.Load(".env")

	if err != nil {

		log.Fatal(err)
	}

	app := gin.Default()

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := &Repository{
		DB: db,
	}

	r.SetupRoutes(app)

	app.Run()

}
