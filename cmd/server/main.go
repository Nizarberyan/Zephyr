package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Nizarberyan/Zephyr/internal/api/handlers"
	"github.com/Nizarberyan/Zephyr/internal/config"
	"github.com/Nizarberyan/Zephyr/internal/database"
	"github.com/Nizarberyan/Zephyr/internal/services"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.Load()

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	defer db.Close()

	queries := database.New(db)
	subsonic := services.NewSubsonicService(cfg.NavidromeURL)
	h := handlers.NewHandler(queries, subsonic, cfg.JWTSecret)

	app := fiber.New()

	// Temporary logger
	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("🔍 INCOMING REQUEST: %s %s\n", c.Method(), c.OriginalURL())
		return c.Next()
	})

	// Public Routes
	app.Post("/login", h.Login)
	app.Post("/submit-listens", h.SubmitListens)

	// Protected Chart Routes
	charts := app.Group("/charts", h.JWTMiddleware)
	charts.Get("/top-artists", h.GetTopArtists)
	charts.Get("/top-albums", h.GetTopAlbums)
	charts.Get("/top-tracks", h.GetTopTracks)
	charts.Get("/recent", h.GetRecent)
	charts.Get("/timeline", h.GetHistoryTimeline)
	charts.Get("/listens", h.GetListens) // Moved listens inside charts for protection

	log.Fatal(app.Listen(cfg.Port))
}
