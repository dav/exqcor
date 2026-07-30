DIST := server/internal/webui/dist

.PHONY: web build dev-server dev-web test clean ensure-dist

# Build the SPA into the Go embed directory.
web:
	cd web && npm install --no-audit --no-fund && npm run build

# Full production build: frontend + single binary at server/exqcor.
build: web
	cd server && CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X github.com/dav/exqcor/server/internal/version.Version=$$(git describe --tags --always)" \
		-o exqcor ./cmd/exqcor

# go build needs dist/ to exist for //go:embed; create a placeholder when the
# frontend hasn't been built yet (fresh clone, backend-only work).
ensure-dist:
	@mkdir -p $(DIST)
	@test -f $(DIST)/index.html || printf '<!doctype html><title>Exqcor</title><p>Frontend not built. Run <code>make web</code>.</p>' > $(DIST)/index.html

# Dev loop: run these in two terminals.
#   make dev-server   — Go API on :8080 (rebuild+rerun on change yourself, or use entr)
#   make dev-web      — Vite on :5173 with /api proxied to :8080
dev-server: ensure-dist
	@mkdir -p tmp
	cd server && go run ./cmd/exqcor --no-open --db ../tmp/dev.db

dev-web:
	cd web && npm run dev

test: ensure-dist
	cd server && go test ./...

clean:
	rm -rf $(DIST) server/exqcor web/node_modules
