package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	MainResponse "github.com/irbbook/go/fiber-gorm/response"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = "3306"
	database = "cine_sysnc"
	username = "root"
	password = "P@ssw0rd"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to Connect Database")
	}

	fmt.Println("Connected Database Successfuly")

	db.AutoMigrate(&Person{}, &User{})
	fmt.Println("Migration Successfully")

	//setup fiber
	app := fiber.New()

	app.Get("/api/test", func(c *fiber.Ctx) error {
		response := MainResponse.Custom(fiber.StatusOK, "I use fresh", nil)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	personRoute := app.Group("/api/person")

	personRoute.Get("/", func(c *fiber.Ctx) error {
		result, err := getPeople(db)
		if err != nil {
			response := MainResponse.InternalServerError(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		response := MainResponse.SuccessWithData(result)
		return c.Status(fiber.StatusOK).JSON(response)

	})

	personRoute.Get("/:id", func(c *fiber.Ctx) error {
		personId, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			response := MainResponse.BadRequest(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		result, err := getPerson(db, uint(personId))

		if err != nil {
			response := MainResponse.NotFound(err.Error())
			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		response := MainResponse.SuccessWithData(result)
		return c.JSON(response)
	})

	personRoute.Post("/", func(c *fiber.Ctx) error {
		currentPerson := new(Person)
		if err := c.BodyParser(currentPerson); err != nil {
			response := MainResponse.BadRequest(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err := createPerson(db, currentPerson)

		if err != nil {
			response := MainResponse.InternalServerError(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		response := MainResponse.Success()
		return c.Status(fiber.StatusCreated).JSON(response)
	})

	personRoute.Put("/", func(c *fiber.Ctx) error {
		currentPerson := new(Person)
		if err := c.BodyParser(currentPerson); err != nil {
			return fiber.ErrBadRequest
		}

		result, err := updatePerson(db, currentPerson)

		if err != nil {
			response := MainResponse.InternalServerError(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		response := MainResponse.SuccessWithData(result)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	personRoute.Delete("/:id", func(c *fiber.Ctx) error {
		personId, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			response := MainResponse.BadRequest(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err = deletePerson(db, uint(personId))

		if err != nil {
			response := MainResponse.InternalServerError(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		response := MainResponse.InternalServerError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	})

	log.Fatal(app.Listen(":8080"))

}
