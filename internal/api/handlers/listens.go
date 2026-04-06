package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Nizarberyan/Zephyr/internal/api/models"
	"github.com/Nizarberyan/Zephyr/internal/database"
	"github.com/Nizarberyan/Zephyr/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	Queries   *database.Queries
	Subsonic  *services.SubsonicService
	JWTSecret string
}

func NewHandler(queries *database.Queries, subsonic *services.SubsonicService, jwtSecret string) *Handler {
	return &Handler{
		Queries:   queries,
		Subsonic:  subsonic,
		JWTSecret: jwtSecret,
	}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password,omitempty"`
		Token    string `json:"token,omitempty"`
		Salt     string `json:"salt,omitempty"`
	}

	req := new(loginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ok, err := h.Subsonic.VerifySubsonic(req.Username, req.Password, req.Token, req.Salt)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Failed to verify credentials with Navidrome"})
	}

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Navidrome credentials"})
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"status":   "ok",
		"username": req.Username,
		"token":    t,
	})
}

func (h *Handler) JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Format: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("username", claims["username"])

	return c.Next()
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

func (h *Handler) GetTopArtists(c *fiber.Ctx) error {
	params, err := parseChartParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Enforce user isolation: use username from JWT
	params.Username = c.Locals("username").(string)

	topArtists, err := h.Queries.GetTopArtists(c.Context(), database.GetTopArtistsParams{
		ListenedAt:   params.From,
		ListenedAt_2: params.To,
		Username:     params.Username,
		Limit:        params.Limit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch top artists"})
	}
	return c.JSON(topArtists)
}

func (h *Handler) GetTopAlbums(c *fiber.Ctx) error {
	params, err := parseChartParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Enforce user isolation: use username from JWT
	params.Username = c.Locals("username").(string)

	topAlbums, err := h.Queries.GetTopAlbums(c.Context(), database.GetTopAlbumsParams{
		ListenedAt:   params.From,
		ListenedAt_2: params.To,
		Username:     params.Username,
		Limit:        params.Limit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch top albums"})
	}
	return c.JSON(topAlbums)
}

func (h *Handler) GetTopTracks(c *fiber.Ctx) error {
	params, err := parseChartParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Enforce user isolation: use username from JWT
	params.Username = c.Locals("username").(string)

	topTracks, err := h.Queries.GetTopTracks(c.Context(), database.GetTopTracksParams{
		ListenedAt:   params.From,
		ListenedAt_2: params.To,
		Username:     params.Username,
		Limit:        params.Limit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch top tracks"})
	}
	return c.JSON(topTracks)
}

func (h *Handler) GetHistoryTimeline(c *fiber.Ctx) error {
	params, err := parseChartParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Enforce user isolation: use username from JWT
	params.Username = c.Locals("username").(string)

	// Additional filter: artist name
	artistName := c.Query("artist", "")

	// Interval (TimescaleDB style, e.g., '1 day', '1 hour')
	interval := c.Query("interval", "1 day")

	history, err := h.Queries.GetPlayHistoryTimeline(c.Context(), database.GetPlayHistoryTimelineParams{
		TimeBucket:   interval,
		ListenedAt:   params.From,
		ListenedAt_2: params.To,
		Username:     params.Username,
		ArtistName:   artistName,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch play history timeline"})
	}
	return c.JSON(history)
}

func (h *Handler) GetRecent(c *fiber.Ctx) error {
	// Enforce user isolation: use username from JWT
	username := c.Locals("username").(string)
	limit := int32(c.QueryInt("limit", 20))

	recent, err := h.Queries.GetRecentListens(c.Context(), database.GetRecentListensParams{
		Username: username,
		Limit:    limit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch recent listens"})
	}
	return c.JSON(recent)
}

// chartParams helper for common query parameters
type chartParams struct {
	From     time.Time
	To       time.Time
	Username string
	Limit    int32
}

func parseChartParams(c *fiber.Ctx) (*chartParams, error) {
	// Default 'from' is 30 days ago if not provided
	fromStr := c.Query("from")
	from := time.Now().AddDate(0, 0, -30)
	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		} else {
			return nil, fmt.Errorf("invalid 'from' date format (expected RFC3339)")
		}
	}

	// Default 'to' is now
	toStr := c.Query("to")
	to := time.Now()
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		} else {
			return nil, fmt.Errorf("invalid 'to' date format (expected RFC3339)")
		}
	}

	return &chartParams{
		From:     from,
		To:       to,
		Username: c.Query("username", ""),
		Limit:    int32(c.QueryInt("limit", 10)),
	}, nil
}
