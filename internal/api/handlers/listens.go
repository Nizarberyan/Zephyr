package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/Nizarberyan/Zephyr/internal/api/models"
	"github.com/Nizarberyan/Zephyr/internal/database"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Queries *database.Queries
}

func NewHandler(queries *database.Queries) *Handler {
	return &Handler{Queries: queries}
}

func (h *Handler) SubmitListens(c *fiber.Ctx) error {
	payload := new(models.ListenPayload)

	// Parse the incoming JSON from Navidrome
	if err := c.BodyParser(payload); err != nil {
		fmt.Println("❌ JSON Parse Error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// ListenBrainz sends "playing_now" when a track starts. We don't save these to the DB.
	if payload.ListenType == "playing_now" {
		fmt.Println("🎧 User is playing a song (Skipping DB insert)")
		return c.JSON(fiber.Map{"status": "ok"})
	}

	// Insert each listen into the database
	ctx := c.Context()
	for _, listen := range payload.Payload {
		listenedAt := time.Unix(listen.ListenedAt, 0)

		inserted, err := h.Queries.CreateScrobble(ctx, database.CreateScrobbleParams{
			ListenedAt:  listenedAt,
			Username:    listen.NavidromeUsername,
			UserUuid:    listen.NavidromeUUID,
			CoverArtID:  listen.CoverArtID,
			DurationMs:  int32(listen.DurationMs),
			ArtistName:  listen.TrackMetadata.ArtistName,
			TrackName:   listen.TrackMetadata.TrackName,
			ReleaseName: listen.TrackMetadata.ReleaseName,
		})

		if err != nil {
			log.Printf("❌ Failed to insert scrobble: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save listen"})
		}
		fmt.Printf("✅ Inserted scrobble: %s by %s\n", inserted.TrackName, inserted.ArtistName)
	}

	// ListenBrainz expects a "status: ok" response
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) GetListens(c *fiber.Ctx) error {
	scrobbles, err := h.Queries.GetAllScrobbles(c.Context())
	if err != nil {
		log.Printf("❌ Failed to fetch scrobbles: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch scrobbles"})
	}
	return c.JSON(scrobbles)
}
