# AGENTS.md

## Build & Test Commands
- Build: `go build ./...`
- Test all: `go test ./...`
- Test single: `go test -v -run TestName ./path/to/package`
- Lint: `golangci-lint run`
- Format: `go fmt ./...`

## Architecture
- API: REST API for point-of-sale (kasir) system
- Structure: cmd/ (entrypoints), internal/ (business logic), pkg/ (shared libs)
- Database: PostgreSQL (Supabase) using database/sql with pgx driver
- Config: Viper for environment-based configuration (.env support)
- Pattern: Layered Architecture (Handler → Service → Repository → Model)
- DI: Manual dependency injection in cmd/api/main.go

## Code Style
- Imports: stdlib first, then external, then internal (grouped with blank lines)
- Naming: camelCase for unexported, PascalCase for exported
- Errors: Return errors, don't panic; wrap with `fmt.Errorf("context: %w", err)`
- Types: Use strong typing; avoid interface{}/any unless necessary
- Formatting: Use `gofmt`/`goimports`; max line ~100 chars

## Kasir API Implementation Workflow

**Completed Epic: kasir-api-12a** - Build Kasir (POS) REST API in Go from scratch
**In Progress Epic: kasir-api-layered** - Refactor to Layered Architecture with PostgreSQL

### Task Completion Order (Layered Architecture)
1. **kasir-api-layered.1**: Project Structure - Create cmd/, internal/, pkg/ directories
2. **kasir-api-layered.2**: Config Management - Viper setup with .env support
3. **kasir-api-layered.3**: Database Connection - PostgreSQL with pgx driver
4. **kasir-api-layered.4-8**: Category Layers - Model, Repository, Service, Handler, DI
5. **kasir-api-layered.9-14**: Product Layers - Model, Repository, Service, Handler, DI
6. **kasir-api-layered.15**: Challenge - Add category_id with JOIN (optional)
7. **kasir-api-layered.16**: Railway Deployment - Update config for database
8. **kasir-api-layered.17**: Testing & Quality Gates - Run build, format, lint

### Project Structure
```
cmd/api/main.go          # Entry point with DI wiring
internal/
  config/                # Viper configuration
  db/                    # PostgreSQL connection
  category/              # Category feature
    model/category.go
    repository/
      repository.go      # Interface
      postgres.go        # Implementation
    service/service.go
    handler/handler.go
  product/               # Product feature (same structure)
pkg/utils/               # Shared utilities
```

### Layered Architecture Flow
```
HTTP Request
    ↓
Handler (parse request, validate, call service)
    ↓
Service (business logic, orchestration)
    ↓
Repository (SQL queries, database access)
    ↓
Model (data structures)
    ↓
PostgreSQL (Supabase)
```

### Task Completion Order
1. **kasir-api-12a.1**: Project Setup - `go mod init kasir-api`, create main.go with all imports
2. **kasir-api-12a.2**: Produk Struct - Defined with ID, Nama (string), Harga (int), Stok (int), JSON tags
3. **kasir-api-12a.3**: Health Check - `GET /health` returns `{"status":"ok"}`
4. **kasir-api-12a.4**: GET/POST /api/produk - List all products, create new with auto-incremented ID
5. **kasir-api-12a.5**: GET /api/produk/{id} - Retrieve single product by ID
6. **kasir-api-12a.6**: PUT /api/produk/{id} - Update product (ID preserved)
7. **kasir-api-12a.7**: DELETE /api/produk/{id} - Delete product from slice
8. **kasir-api-12a.8**: Main function & Routing - http.HandleFunc for all endpoints, ListenAndServe on :8080
9. **kasir-api-12a.9**: Build binary - `go build -ldflags="-s -w" -o kasir-api` (5.3MB)

### Implementation Details
- **Router**: http.HandleFunc with default ServeMux (no frameworks)
- **Routes**: `/api/produk` (exact) for list/create, `/api/produk/` (prefix) for ID operations
- **Storage**: In-memory slice `[]Produk` initialized with 3 sample products
- **ID Parsing**: `strings.TrimPrefix` + `strconv.Atoi`
- **JSON Tags**: Lowercase (id, nama, harga, stok)
- **Harga Type**: `int` (Rupiah, no decimals)

### API Endpoints
| Method | Path | Status | Handler |
|--------|------|--------|---------|
| GET | /health | 200 | healthHandler |
| GET | /api/produk | 200 | getAllProduk |
| POST | /api/produk | 201 | createProduk |
| GET | /api/produk/{id} | 200/404 | getProdukByID |
| PUT | /api/produk/{id} | 200/404 | updateProduk |
| DELETE | /api/produk/{id} | 200/404 | deleteProduk |

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
- ALWAYS git commit for every task completions
