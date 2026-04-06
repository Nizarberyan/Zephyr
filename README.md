# Zephyr: Custom Maloja Alternative & Navidrome Fork

Zephyr is a self-hosted music scrobbling and charting service designed for multi-user support, seamless SSO, and rich mobile visualization.

## Project Overview
This project consists of a modified Navidrome media server and a custom backend to handle advanced scrobbling data and time-series charting.

## Tech Stack
- **Media Server:** Forked Navidrome (Go)
- **Custom Backend:** Go (Fiber framework) + TimescaleDB
- **Custom Frontend:** React Native (Expo) or Flutter
- **Authentication:** Pass-through Subsonic API verification
- **Images:** Served directly from Navidrome via Subsonic `getCoverArt` endpoint

## Project Structure
- `cmd/server/`: Main entry point for the custom backend.
- `internal/api/`: API handlers and models.
- `internal/database/`: Database logic (SQLC + TimescaleDB).
- `internal/services/`: Business logic.
