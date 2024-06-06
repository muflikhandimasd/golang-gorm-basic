package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	logFiber "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Product struct {
	gorm.Model
	Name  string `gorm:"unique" json:"name"`
	Code  string `gorm:"unique"  json:"code"`
	Price uint   `json:"price"`
}

func main() {
	fmt.Println("JALAN")
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Product{})

	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	app.Use(logFiber.New(
		logFiber.Config{
			Format:     "${time} - ${ip} - ${status} - ${method} ${path} - ${latency}\n",
			TimeFormat: "02-Jan-2006",
			TimeZone:   "Local",
		}))

	app.Use(cors.New())

	app.Get("/api/products", func(ctx *fiber.Ctx) error {
		var status int
		var message string
		var data []Product

		status = fiber.StatusOK
		message = "Berhasil Get Products"

		res := db.Find(&data)

		if err = res.Error; err != nil {
			fmt.Println("err get products", err)
			status = fiber.StatusInternalServerError
			message = "Internal Server Error"
		}

		jsonMap := fiber.Map{}
		jsonMap["status"] = status
		jsonMap["data"] = data
		jsonMap["message"] = message
		return ctx.Status(status).JSON(jsonMap)
	})

	app.Post("/api/products", func(ctx *fiber.Ctx) error {
		var status int
		var message string
		var data interface{}

		prod := Product{}

		if err = ctx.BodyParser(&prod); err != nil {
			status = fiber.StatusBadRequest
			message = "Permintaan tidak valid"
		} else {
			res := db.Create(&prod)
			if err = res.Error; err != nil {
				status = fiber.StatusInternalServerError
				fmt.Println("ERR CREATE", err)
				message = "INTERNAL SERVER ERROR"
			} else {
				status = fiber.StatusOK
				data = res.RowsAffected
			}
		}

		jsonMap := fiber.Map{}

		jsonMap["status"] = status
		jsonMap["data"] = data
		jsonMap["message"] = message
		return ctx.Status(status).JSON(jsonMap)

	})

	app.Put("/api/products/:id", func(ctx *fiber.Ctx) error {
		var status int
		var message string
		var data interface{}

		prod := Product{}

		if err = ctx.BodyParser(&prod); err != nil {
			status = fiber.StatusBadRequest
			message = "Permintaan tidak valid"
		} else {
			id := ctx.Params("id")
			res := db.Model(&Product{}).Where("id = ?", id).Updates(&prod)
			if err = res.Error; err != nil {
				status = fiber.StatusInternalServerError
				fmt.Println("ERR CREATE", err)
				message = "INTERNAL SERVER ERROR"
			} else {
				status = fiber.StatusOK
				data = res.RowsAffected
			}
		}

		jsonMap := fiber.Map{}

		jsonMap["status"] = status
		jsonMap["data"] = data
		jsonMap["message"] = message
		return ctx.Status(status).JSON(jsonMap)
	})

	app.Delete("/api/products/:id", func(ctx *fiber.Ctx) error {
		var status int
		var message string
		var data interface{}

		id := ctx.Params("id")
		res := db.Delete(&Product{}, id)
		if err = res.Error; err != nil {
			status = fiber.StatusInternalServerError
			fmt.Println("ERR CREATE", err)
			message = "INTERNAL SERVER ERROR"
		} else {
			status = fiber.StatusOK
			data = res.RowsAffected
		}

		jsonMap := fiber.Map{}

		jsonMap["status"] = status
		jsonMap["data"] = data
		jsonMap["message"] = message
		return ctx.Status(status).JSON(jsonMap)
	})

	log.Fatal(app.Listen(":3000"))
}
