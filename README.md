

# URL Shortener (Go + Gin)

A simple URL shortener built in **Golang** using the **Gin** framework.  
This project is structured to be **phase-wise**, allowing step-by-step learning and gradual enhancement.

---

## Features (Phase-wise)

### Phase 1: Core Basics
- Shorten a long URL (`POST /shorten`)
- Retrieve the original URL (`GET /:code`)
- Track click counts for each short URL (`GET /stats/:code`)
- Thread-safe in-memory storage using `sync.Mutex`
- Simple numeric short codes (`1`, `2`, `3`...)

### Phase 2: Unique Short URLs
- Random **Base62 short codes** for URLs (e.g., `Ka9L2s`)
- Collision checking to ensure uniqueness

### Upcoming Phases
- PostgreSQL integration for persistence
- Redis caching for faster reads
- Scalable and distributed setup
- Analytics and monitoring
- User-defined/custom short URLs

---

## Installation & Run

1. Clone the repo:

```bash
git clone https://github.com/<your-username>/url-shortener.git
cd url-shortener
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the server:

```bash
go run main.go
```

Server will start at `http://localhost:8080`.

---

## API Endpoints

### 1. Shorten a URL
- **POST** `/shorten`  
- Request body:

```json
{
  "url": "https://example.com"
}
```

- Response:

```json
{
  "short_url": "http://localhost:8080/1"
}
```

---

### 2. Redirect
- **GET** `/:code`  
- Example: `http://localhost:8080/1` → redirects to `https://example.com`

---

### 3. Get URL Stats
- **GET** `/stats/:code`  
- Response:

```json
{
  "short_code": "1",
  "original_url": "https://example.com",
  "click_count": 3
}
```

---

## Notes
- Currently, URLs are stored in-memory; all data is lost when the server restarts.
- Future phases will include persistent storage and caching for production readiness.

---

## License
MIT License