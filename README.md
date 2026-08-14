# gofeed

`gofeed` is a concurrent RSS feed aggregator and RESTful API backend written in Go. It allows users to manage and follow RSS feeds, automatically fetching and saving new posts in the background using Go concurrency primitives.

---

## 🌟 Key Features

- **Concurrent Feed Scraping:** Background worker implementation using goroutines, `sync.WaitGroup`, and configurable `time.Ticker` intervals to fetch multiple RSS feeds concurrently.
- **Feed Parsing:** Parses RSS 2.0 XML structures, unescapes HTML entities, and handles multiple `pubDate` formats (RFC1123, RFC3339, etc.).
- **Authentication:** User registration and authentication by **Argon2id** password hashing and **JWT** for HTTP authorization.
- **Transactional DB Operations:** Uses PostgreSQL transactions to ensure atomicity when creating feeds and feed follow and when refreshing access token.
- **Type-Safe Database Access:** Type-safe SQL query generation using **sqlc** and schema migrations managed via **goose**.
- **Automated Test Suite:** Integration and unit test suite using `net/http/httptest` with mock servers and isolated database test environments.

---

## 🛠️ Tech Stack

- **Language:** [Go 1.22+](https://go.dev/)
- **HTTP Router:** [Chi Router](https://github.com/go-chi/chi)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **SQL Code Generation:** [sqlc](https://sqlc.dev/)
- **Database Migrations:** [goose](https://github.com/pressly/goose)
- **Password Hashing:** `golang.org/x/crypto/argon2`
- **Testing:** Standard `testing` package with `httptest`

---

## 🏗️ System Architecture

```mermaid
graph TD
    Client[HTTP Client] -->|REST API Requests| Router[Chi Router]
    Router -->|MiddlewareAuth| Handlers[HTTP Handlers]
    Handlers -->|sqlc Queries| DB[(PostgreSQL)]
    
    subgraph Background Service
        Scraper[Scraper Worker] -->|Fetch RSS XML| ExternalFeeds[External RSS Feeds]
        Scraper -->|Insert New Posts| DB
    end
```
## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) (v1.21 or later)
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [goose](https://github.com/pressly/goose) for database migrations


### Clone & Configure Environment

Clone the repository and create a `.env` file in the root directory:

```bash
git clone https://github.com/MhdFiras-3/gofeed.git
cd gofeed
```

`.env`:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=gofeed
TEST_DB_URL=postgres://postgres:yourpassword@localhost:5432/gofeed_test?sslmode=disable
JWT_SECRET=yoursecret
DUMMY_HASH="$argon2id$v=19$m=65536,t=1,p=12$LpTO8GOk8ajNAFczIs12uQ$e48btjY28JEWiIfEDYcJBjb1GLBZGqXoOyrClQ9EIr0"
```


Start PostgreSQL Container:

```bash
docker compose up -d
```
Apply schema migrations using goose and run the API server:
```bash
goose -dir sql/migrations postgres "postgres://postgres:yourpassword@localhost:5432/gofeed?sslmode=disable" up
go run cmd/server/main.go
```
To run :
```bash
go test -p 1 ./...
```
## API Endpoints Summary
#### Public Routes
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/register` | Register a new user account |
| `POST` | `/api/v1/login` | Authenticate user and receive JWT access token |
| `POST` | `/api/v1/refresh` | Refresh access token using refresh token |
| `POST` | `/api/v1/logout` | Revoke refresh token / logout |
| `GET` | `/api/v1/feeds` | Get all created RSS feeds |
| `GET` | `/api/v1/feeds/{feedID}` | Get a specific RSS feed by ID |

#### Protected Routes
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/me` | Retrieve current user profile |
| `PATCH` | `/api/v1/me` | Update user profile details |
| `DELETE` | `/api/v1/me` | Delete current user account |
| `POST` | `/api/v1/feeds` | Create a new RSS feed and automatically follow it |
| `GET` | `/api/v1/follows` | Get all feed follows for current user |
| `DELETE` | `/api/v1/follows/{feedID}` | Unfollow a feed by ID |
| `GET` | `/api/v1/posts` | Fetch RSS posts for followed feeds |
