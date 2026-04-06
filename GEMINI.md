# GEMINI Context: Zephyr Development

## Core Mission
We are building a custom scrobbling ecosystem. The primary goal is to modify a Navidrome fork to emit enriched ListenBrainz-style scrobble payloads to our custom backend.

## Architectural Mandates
- **No External Metadata APIs:** All cover art and metadata must be sourced from the Navidrome instance via Subsonic APIs.
- **Time-Series Priority:** The backend uses TimescaleDB for high-performance charting and history.
- **SSO Focus:** User identification must be seamless, leveraging Navidrome's internal user IDs and names.

## Immediate Task: Navidrome Scrobbler Mutation
We need to modify the Navidrome source code (specifically the ListenBrainz scrobbler) to inject the following fields into the `payload` array of the JSON body:

1.  `navidrome_username`: The string username.
2.  `navidrome_uuid`: The internal user UUID.
3.  `cover_art_id`: The ID for the album/track cover art.
4.  `duration_ms`: Total track length in milliseconds.

### Target JSON Format
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

## Implementation Notes
- The scrobbler logic in Navidrome typically resides in its internal `scrobbler` or `listenbrainz` packages.
- We must ensure these fields are captured at the point of scrobble dispatch.
