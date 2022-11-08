package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stelgkio/go-fiber-postgres/models"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(contex *fiber.Ctx) error {
	book := Book{}

	err := contex.BodyParser(&book)

	if err != nil {
		contex.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		contex.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}
	contex.Status(http.StatusOK).JSON(&fiber.Map{"message": "book added"})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {

	booksModels := &[]models.Books{}

	err := r.DB.Find(booksModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"messasge": ""})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"messasge": "books fetched successfully",
		"data":     booksModels,
	})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/get_books/", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Connection to database fail")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}
