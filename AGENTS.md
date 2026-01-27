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
- Database: Update this when DB is added (PostgreSQL/MySQL recommended)

## Code Style
- Imports: stdlib first, then external, then internal (grouped with blank lines)
- Naming: camelCase for unexported, PascalCase for exported
- Errors: Return errors, don't panic; wrap with `fmt.Errorf("context: %w", err)`
- Types: Use strong typing; avoid interface{}/any unless necessary
- Formatting: Use `gofmt`/`goimports`; max line ~100 chars

## Kasir API Implementation Workflow

**Completed Epic: kasir-api-12a** - Build Kasir (POS) REST API in Go from scratch

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
