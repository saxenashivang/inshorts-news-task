# Contextual News Data Retrieval System - Inshorts

## Overview
- **Geo-spatial queries** (PostGIS)
- **Natural Language Understanding** (Langchain) for intent extraction
- **Trending News** (Redis Sorted Sets)
- **Clean Architecture** (Handlers -> UseCases -> Repositories)

## Stack
- **Go 1.22+**
- **Gin** (Web Framework)
- **PostgreSQL + PostGIS** (Storage)
- **Redis** (Caching/Trending)
- **Langchain** (Gemini Tested)

## Setup

1. **Prerequisites**: Docker, Go 1.22+
2. **Start Infrastructure**:
   ```bash
   docker-compose up -build
   ```
3. **Run Migrations**:
   (You can use a tool like `migrate` or manually run the SQL in `migrations/`)
   ```bash
   cat migrations/001_schema.sql | docker exec -i inshorts-task-postgres-1 psql -U user -d news_db
   ```

## Endpoints

- `GET /api/v1/news?q=Tech+news`
  - Uses full-text search.
- `GET /api/v1/news?q=News+near+me&lat=37.77&lng=-122.41`
  - Uses `st_dwithin` for geospatial search.
  - Automatically enriches response with an AI summary.

## Architecture

- `cmd/api`: Entrypoint
- `internal/core/entity`: Domain Models
- `internal/core/port`: Interfaces (Ports)
- `internal/core/usecase`: Business Logic
- `internal/adapter`: Implementations (SQL, Redis, HTTP)

## Design Decisions

- **Clean Architecture**: Decouples business logic from external frameworks.
- **PostGIS**: Efficient geospatial queries (`ST_DWithin`).
- **Redis Sorted Sets**: Efficient leaderboard for trending news.
- **LLM as Parser**: Transforms natural language into structured SQL queries (Intent/Entities).
- **Graceful Shutdown**: (Implied in production ready setup, though minimal in this demo).
