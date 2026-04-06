package models

type ListenPayload struct {
	ListenType string `json:"listen_type"`
	Payload    []struct {
		ListenedAt        int64  `json:"listened_at"`
		NavidromeUsername string `json:"navidrome_username"`
		NavidromeUUID     string `json:"navidrome_uuid"`
		CoverArtID        string `json:"cover_art_id"`
		DurationMs        int    `json:"duration_ms"`
		TrackMetadata     struct {
			ArtistName  string `json:"artist_name"`
			TrackName   string `json:"track_name"`
			ReleaseName string `json:"release_name"`
		} `json:"track_metadata"`
	} `json:"payload"`
}
