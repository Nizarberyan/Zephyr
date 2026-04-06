package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Nizarberyan/Zephyr/internal/api/handlers"
	"github.com/Nizarberyan/Zephyr/internal/database"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	defer db.Close()

	queries := database.New(db)
	h := handlers.NewHandler(queries)

	app := fiber.New()

	// Temporary logger
	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("🔍 INCOMING REQUEST: %s %s\n", c.Method(), c.OriginalURL())
		return c.Next()
	})

	// Routes
	app.Post("/submit-listens", h.SubmitListens)
	app.Get("/listens", h.GetListens)

	log.Fatal(app.Listen(":3002"))
}
