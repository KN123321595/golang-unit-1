CREATE TABLE apod_metadata (
    id SERIAL PRIMARY KEY,
    explanation TEXT,
    media_type TEXT,
    service_version TEXT,
    title TEXT,
    hdurl TEXT,
    url TEXT,
    image_path TEXT,
    date TIMESTAMP
);
CREATE INDEX image_path_index ON apod_metadata (image_path);
CREATE INDEX date_index ON apod_metadata (date);

CREATE TABLE cron_logs (
    id SERIAL PRIMARY KEY,
    job_name TEXT,
    start_time TIMESTAMP,
    end_time TIMESTAMP 
);
CREATE INDEX job_name_index ON cron_logs (job_name);