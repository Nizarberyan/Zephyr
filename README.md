# Zephyr: Custom Maloja Alternative

## 1. Project Overview
We are building a custom, self-hosted music scrobbling and charting service (a Maloja alternative) with a heavy focus on multi-user support, seamless single sign-on (SSO), and rich mobile charting. 
To achieve seamless integration without requiring users to generate API keys, we are forking Navidrome to modify its native ListenBrainz scrobbling payload.

## 2. Tech Stack Architecture
* **Media Server:** Forked Navidrome (Go).
* **Custom Backend:** Go (Fiber framework) + TimescaleDB (for heavy time-series data).
* **Custom Frontend:** React Native (Expo) or Flutter (using `fl_chart`).
* **Authentication Mechanism:** Pass-through Subsonic API verification.
* **Image Serving Strategy:** No external API calls. The frontend will fetch images directly from Navidrome via the Subsonic `getCoverArt` endpoint.

## 3. Navidrome Source Modifications (Completed)
The Go source code of Navidrome has been modified to alter the ListenBrainz scrobbler payload to inject custom internal data.

**Behavior:**
When Navidrome completes a track, it sends a standard ListenBrainz JSON payload to the URL defined in `ND_LISTENBRAINZ_BASEURL`. 
The JSON payload now includes four new fields injected alongside `track_metadata`:
1.  `navidrome_username`: The string username of the user who played the track.
2.  `navidrome_uuid`: The internal database UUID of the user.
3.  `cover_art_id`: The ID of the album/track cover art.
4.  `duration_ms`: The total length/duration of the track in milliseconds.

**Expected JSON Output Format:**
```json
{
  "listen_type": "single",
  "payload": [
    {
      "listened_at": 1775166795,
      "navidrome_username": "example_user",
      "navidrome_uuid": "550e8400-e29b-41d4-a716-446655440000",
      "cover_art_id": "a1b2c3d4",
      "duration_ms": 210000,
      "track_metadata": {
        "artist_name": "Artist Name",
        "track_name": "Song Name",
        "release_name": "Album Name"
      }
    }
  ]
}
```
