CREATE TABLE apod_metadata (
    id SERIAL PRIMARY KEY,
    explanation text,
    media_type text,
    service_version text,
    title text,
    hdurl text,
    url text,
    image_path text,
    date date
);

CREATE TABLE cron_logs (
    id serial PRIMARY KEY,
    job_name text,
    start_time TIMESTAMP,
    end_time TIMESTAMP 
);