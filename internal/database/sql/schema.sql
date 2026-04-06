CREATE TABLE scrobbles (
                           id BIGSERIAL,
                           listened_at TIMESTAMPTZ NOT NULL,
                           username TEXT NOT NULL,
                           user_uuid TEXT NOT NULL,
                           cover_art_id TEXT NOT NULL,
                           duration_ms INT NOT NULL,
                           artist_name TEXT NOT NULL,
                           track_name TEXT NOT NULL,
                           release_name TEXT NOT NULL,
    -- Timescale requires the time column to be part of the primary key
                           PRIMARY KEY (id, listened_at)
);

-- Convert the standard table into a Timescale hypertable partitioned by time
SELECT create_hypertable('scrobbles', 'listened_at');