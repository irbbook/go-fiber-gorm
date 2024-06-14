package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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

	db.AutoMigrate(&Person{})
	fmt.Println("Migration Successfully")

	//setup fiber
	app := fiber.New()

	app.Get("/api/test", func(c *fiber.Ctx) error {
		return c.SendString("i use fresh")
	})

	personRoute := app.Group("/api/person")

	personRoute.Get("/", func(c *fiber.Ctx) error {
		result, err := getPeople(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(result)
	})

	personRoute.Get("/:id", func(c *fiber.Ctx) error {
		personId, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		result, err := getPerson(db, uint(personId))

		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		return c.JSON(result)
	})

	personRoute.Post("/", func(c *fiber.Ctx) error {
		currentPerson := new(Person)
		if err := c.BodyParser(currentPerson); err != nil {
			return fiber.ErrBadRequest
		}

		err := createPerson(db, currentPerson)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	personRoute.Put("/", func(c *fiber.Ctx) error {
		currentPerson := new(Person)
		if err := c.BodyParser(currentPerson); err != nil {
			return fiber.ErrBadRequest
		}

		result, err := updatePerson(db, currentPerson)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(result)
	})

	personRoute.Delete("/:id", func(c *fiber.Ctx) error {
		personId, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid id")
		}

		err = deletePerson(db, uint(personId))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	log.Fatal(app.Listen(":8080"))

}
