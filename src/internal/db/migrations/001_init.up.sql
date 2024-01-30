CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.screen_opening_time(
    session_id String,
    platform String,
    timestamp DateTime,
    duration Int64,
    screen_name String
) ENGINE = MergeTree
PRIMARY KEY (session_id, timestamp);

CREATE TABLE IF NOT EXISTS analytics.request_time(
    session_id String,
    platform String,
    timestamp DateTime,
    duration Int64,
    request_url String
) ENGINE = MergeTree
PRIMARY KEY (session_id, timestamp);
