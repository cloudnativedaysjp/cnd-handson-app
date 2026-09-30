BUF := go run github.com/bufbuild/buf/cmd/buf@v1.57.0
GO_SERVICES := user session project task
PY_SERVICES := role column
PNPM := corepack pnpm

.PHONY: gen lint up down clean e2e

# Requires: go (runs buf via `go run`) and network access (buf remote plugins).
gen:
	$(BUF) generate $(foreach s,$(GO_SERVICES),--path proto/$(s))
	@for s in $(PY_SERVICES); do \
		rm -rf backend/$$s/gen; \
		$(BUF) generate --template buf.gen.python.yaml --path proto/$$s -o backend/$$s/gen || exit 1; \
	done

# Requires: go, docker (Go services are linted via their Dockerfile lint target).
lint:
	$(BUF) lint
	@for s in user session task; do $(MAKE) -C backend/$$s lint || exit 1; done
	cd backend/project && go vet ./...  # project has no Docker lint target
	@for s in $(PY_SERVICES); do $(MAKE) -C backend/$$s lint || exit 1; done

up:
	@test -f .env || cp .env.example .env
	docker compose up -d --build --wait --wait-timeout 300

down:
	docker compose down

clean:
	docker compose down -v

# Requires: node (corepack provides pnpm) and a running stack (`make up`).
e2e:
	cd e2e && $(PNPM) install --frozen-lockfile && $(PNPM) exec playwright install chromium && $(PNPM) test
