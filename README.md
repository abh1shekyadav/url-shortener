# URL Shortener Service

A robust and scalable URL shortener built in **Golang** using the **Gin** framework, with PostgreSQL for persistent storage and Redis for caching. This service generates unique Base62 short codes, tracks click counts, supports URL expiry, and includes Docker support for easy deployment.

---

## Features

- Generate unique Base62 short codes for URLs (e.g., `Ka9L2s`)
- Persistent storage with PostgreSQL
- Caching with Redis for fast URL retrieval
- Track click counts for each short URL
- Support for URL expiration dates
- Thread-safe and high-performance implementation
- Docker and Docker Compose support for easy setup and deployment

---

## Technologies

- Go (Golang)
- Gin Web Framework
- PostgreSQL
- Redis
- Docker & Docker Compose

---

## Setup Instructions

### Prerequisites

- Docker
- Docker Compose

### Running the Service

1. Clone the repository:

```bash
git clone https://github.com/<your-username>/url-shortener.git
cd url-shortener
```

2. Start the service with Docker Compose:

```bash
docker-compose up --build
```

This will start the URL shortener service along with PostgreSQL and Redis containers.

3. The server will be accessible at `http://localhost:8080`.

---

## Environment Variables

The following environment variables are configured via `docker-compose.yml`:

- `POSTGRES_USER` — PostgreSQL username
- `POSTGRES_PASSWORD` — PostgreSQL password
- `POSTGRES_DB` — PostgreSQL database name
- `REDIS_ADDR` — Redis server address
- `REDIS_PASSWORD` — Redis password (if any)
- `BASE_URL` — Base URL for generating short URLs (e.g., `http://localhost:8080`)

You can customize these variables in the `docker-compose.yml` file as needed.

---

## Database Table

The following PostgreSQL table is used to store shortened URLs, their metadata, and click statistics:

```sql
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(16) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    click_count INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

This schema ensures each short code is unique, supports optional expiry, and tracks click counts. The service automatically manages this table on startup if it does not exist.

---

## API Endpoints

### 1. Shorten a URL

- **POST** `/shorten`

#### Request Body

```json
{
  "url": "https://example.com",
  "expires_at": "2024-07-01T12:00:00Z"  // Optional, ISO8601 format
}
```

#### Response

```json
{
  "success": true,
  "data": {
    "short_code": "Ka9L2s",
    "short_url": "http://localhost:8080/Ka9L2s",
    "original_url": "https://example.com",
    "expires_at": "2024-07-01T12:00:00Z"
  }
}
```

---

### 2. Redirect to Original URL

- **GET** `/:code`

Redirects to the original URL associated with the given short code.

- Returns HTTP 404 if the code does not exist or the URL has expired.

---

### 3. Get URL Stats

- **GET** `/stats/:code`

#### Response

```json
{
  "success": true,
  "data": {
    "short_code": "Ka9L2s",
    "original_url": "https://example.com",
    "click_count": 42,
    "expires_at": "2024-07-01T12:00:00Z"
  }
}
```

---

## Notes

- URLs are stored persistently in PostgreSQL with caching in Redis for improved performance.
- URL expiry is enforced; expired URLs will not redirect and will return appropriate errors.
- Docker Compose simplifies running the entire stack locally.
- Customize environment variables in `docker-compose.yml` to suit your deployment needs.

---

## License

MIT License