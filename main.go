package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yash5060/xenotask/models"

	"github.com/yash5060/xenotask/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	// "golang.org/x/text/date"
	"gorm.io/gorm"
)

// func abc (time.Time).Date() (year int, month time.Month, day int)

type Task struct {
	Title string `json:"title"`

	Description *string `json:"description"`
	Due_Date    *string `json:"age"`
	Status      *string `json:"status"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateTask(context *fiber.Ctx) error {
	task := Task{}

	err := context.BodyParser(&task)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&task).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create task"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "task has been added"})
	return nil
}

func (r *Repository) DeleteTask(context *fiber.Ctx) error {
	taskModel := models.Tasks{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(taskModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete task",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "task delete successfully",
	})
	return nil
}

func (r *Repository) GetTask(context *fiber.Ctx) error {
	taskModels := &[]models.Tasks{}

	err := r.DB.Find(taskModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get task"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "task fetched successfully",
		"data":    taskModels,
	})
	return nil
}

func (r *Repository) GetTaskByID(context *fiber.Ctx) error {

	id := context.Params("id")
	taskModel := &models.Tasks{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(taskModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the task"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "task id fetched successfully",
		"data":    taskModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/post_tasks", r.CreateTask)
	api.Delete("delete_task/:id", r.DeleteTask)
	api.Get("/get_task/:id", r.GetTaskByID)
	api.Get("/get_tasks", r.GetTask)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
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
	err = models.MigrateTask(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":3000")
}
