# Contextual News Data Retrieval System - Inshorts

## Overview
- **Geo-spatial queries** (PostGIS)
- **Natural Language Understanding** (OpenAI/Gemini/Antropic LLM) for intent extraction
- **Trending News** (Redis Sorted Sets)
- **Clean Architecture** (Handlers -> UseCases -> Repositories)

## Stack
- **Go 1.22+**
- **Gin** (Web Framework)
- **PostgreSQL + PostGIS** (Storage)
- **Redis** (Caching/Trending)
- **OpenAI API** (Langchain)

## Setup

### Option 1: Docker (Recommended)
1. **Prerequisites**: Docker, Docker Compose
2. **Start Infrastructure & App**:
   ```bash
   docker-compose up --build
   ```
   The API will be available at `http://localhost:8080`.

### Option 2: Local Development
1. **Prerequisites**: Go 1.25+, PostgreSQL + PostGIS, Redis
2. **Start Dependencies**:
   Ensure PostgreSQL and Redis are running.
3. **Configure Environment**:
   Create a `.env` file or export variables:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=user
   export DB_PASSWORD=password
   export DB_NAME=newsdb
   export REDIS_ADDR=localhost:6379
   export LLM_PROVIDER=openai # or gemini, claude
   export OPENAI_API_KEY="your_key"
   ```
4. **Run Application**:
   ```bash
   go mod tidy
   go run cmd/api/main.go
   ```

## Endpoints

### News Retrieval
- `GET /api/v1/news?q=Tech+news`
  - Uses full-text search along with query lat=28.38&lng=77.12.
- `GET /api/v1/news/nearby`
  - Uses `st_dwithin` for geospatial search.
  - Automatically enriches response with an AI summary.

### Trending News
- `POST /api/v1/news/:id/view`
  - Records a view for a specific article (Article ID can be UUID).
  - Used to track trending articles.
- `GET /api/v1/news/trending`
  - Returns the top 10 trending articles based on view counts.
  - Data is fetched from Redis (cached IDs) and enriched from PostgreSQL.

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
