BUF_VERSION := 1.57.0
# 版が違う buf や未設定の mise シムでは生成結果が CI とずれるため、版が一致するときだけ PATH の buf を使う
BUF ?= $(shell [ "$$(buf --version 2>/dev/null)" = "$(BUF_VERSION)" ] && echo buf || echo go run github.com/bufbuild/buf/cmd/buf@v$(BUF_VERSION))
UP_BUILD ?= --build
GO_SERVICES := user session project task
PY_SERVICES := role column
PNPM := pnpm

.PHONY: gen lint up down clean e2e contract

# Requires: buf on PATH, or go (falls back to `go run`) and network access (buf remote plugins).
gen:
	rm -rf gen/go
	$(BUF) generate
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

# 既知の署名鍵で起動させないため、JWT の鍵は生成する
.env:
	sed "s/^JWT_SECRET_KEY=$$/JWT_SECRET_KEY=$$(openssl rand -hex 32)/" .env.example > $@

up: .env
	docker compose up -d $(UP_BUILD) --wait --wait-timeout 300

down: .env
	docker compose down

clean: .env
	docker compose down -v

# Requires: node and pnpm (mise locally, corepack in CI) and a running stack (`make up`).
e2e:
	cd e2e && $(PNPM) install --frozen-lockfile && $(PNPM) exec playwright install chromium && $(PNPM) test

# Contract tests for #65 services (entry/idp/project/task). Expected to fail until they are implemented; not run in CI.
# Addresses come from env (defaults in e2e/tests/contract/lib/config.ts).
# SVC="idp project" で担当サービスの契約テストだけを実行する（ファイル名で絞り込む）
contract:
	cd e2e && $(PNPM) install --frozen-lockfile && $(PNPM) exec playwright install chromium && $(PNPM) test:contract $(SVC)
