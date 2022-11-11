package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	"github.com/stelgkio/go-fiber-postgres/models"
	"github.com/stelgkio/go-fiber-postgres/storage"
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

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}
	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messasge": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"messasge": "books deleted successfully",
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messasge": "could not get the book",
		})
		return err

	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"messasge": "books id fetched successfully",
		"data":     bookModel,
	})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/get_books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_POST"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_DBNAME"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Connection to database fail")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Migration fail")
	}
	r := Repository{
		DB: db,
	}
	app := fiber.New()
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	r.SetupRoutes(app)
	app.Listen(":8080")

}
