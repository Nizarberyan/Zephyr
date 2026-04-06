-- name: CreateScrobble :one
INSERT INTO scrobbles (
    listened_at,
    username,
    user_uuid,
    cover_art_id,
    duration_ms,
    artist_name,
    track_name,
    release_name
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING *;


-- name: GetAllScrobbles :many
SELECT *
FROM scrobbles
ORDER BY listened_at DESC;

-- name: GetTopArtists :many
SELECT artist_name, COUNT(*) as play_count
FROM scrobbles
WHERE listened_at >= $1 AND listened_at <= $2
  AND (username = $3 OR $3 = '')
GROUP BY artist_name
ORDER BY play_count DESC
LIMIT $4;

-- name: GetTopAlbums :many
SELECT release_name, artist_name, cover_art_id, COUNT(*) as play_count
FROM scrobbles
WHERE listened_at >= $1 AND listened_at <= $2
  AND (username = $3 OR $3 = '')
GROUP BY release_name, artist_name, cover_art_id
ORDER BY play_count DESC
LIMIT $4;

-- name: GetTopTracks :many
SELECT track_name, artist_name, release_name, cover_art_id, COUNT(*) as play_count
FROM scrobbles
WHERE listened_at >= $1 AND listened_at <= $2
  AND (username = $3 OR $3 = '')
GROUP BY track_name, artist_name, release_name, cover_art_id
ORDER BY play_count DESC
LIMIT $4;

-- name: GetRecentListens :many
SELECT *
FROM scrobbles
WHERE (username = $1 OR $1 = '')
ORDER BY listened_at DESC
LIMIT $2;

-- name: GetUserStats :one
SELECT
    COUNT(*) AS total_listens,
    COUNT(DISTINCT artist_name) AS unique_artists,
    COUNT(DISTINCT release_name) AS unique_albums,
    COUNT(DISTINCT track_name) AS unique_tracks
FROM scrobbles
WHERE (username = $1 OR $1 = '')
  AND listened_at >= $2 AND listened_at <= $3;

-- name: GetUserPlayTime :one
SELECT COALESCE(SUM(duration_ms), 0)::BIGINT AS total_play_time_ms
FROM scrobbles
WHERE (username = $1 OR $1 = '')
  AND listened_at >= $2 AND listened_at <= $3;

-- name: GetPlaysByHourOfDay :many
SELECT
    EXTRACT(HOUR FROM listened_at)::INT AS hour,
    COUNT(*) AS play_count
FROM scrobbles
WHERE (username = $1 OR $1 = '')
  AND listened_at >= $2 AND listened_at <= $3
GROUP BY hour
ORDER BY hour ASC;

-- name: GetPlaysByDayOfWeek :many
SELECT
    EXTRACT(DOW FROM listened_at)::INT AS day_of_week,
    COUNT(*) AS play_count
FROM scrobbles
WHERE (username = $1 OR $1 = '')
  AND listened_at >= $2 AND listened_at <= $3
GROUP BY day_of_week
ORDER BY day_of_week ASC;

-- name: GetPlayHistoryTimeline :many
SELECT time_bucket($1, listened_at) AS bucket, COUNT(*) AS play_count
FROM scrobbles
WHERE listened_at >= $2 AND listened_at <= $3
  AND (username = $4 OR $4 = '')
  AND (artist_name = $5 OR $5 = '')
GROUP BY bucket
ORDER BY bucket ASC;
