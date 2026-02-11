CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    source VARCHAR(255),
    url TEXT,
    category TEXT[],
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    geom GEOGRAPHY(Point, 4326),
    published_at TIMESTAMP WITH TIME ZONE,
    relevance_score DOUBLE PRECISION,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || content)) STORED
);

CREATE INDEX IF NOT EXISTS article_geom_idx ON articles USING GIST (geom);
CREATE INDEX IF NOT EXISTS article_search_idx ON articles USING GIN (search_vector);
