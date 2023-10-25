CREATE TABLE apod_metadata (
    id SERIAL PRIMARY KEY,
    explanation TEXT,
    media_type TEXT,
    service_version TEXT,
    title TEXT,
    hdurl TEXT,
    thumbnail_url TEXT,
    url TEXT,
    image_path TEXT,
    date TIMESTAMP
);
CREATE INDEX image_path_index ON apod_metadata (image_path);
CREATE INDEX date_index ON apod_metadata (date);