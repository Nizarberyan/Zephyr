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

