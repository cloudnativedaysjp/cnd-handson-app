BUF_VERSION := 1.57.0
# 版が違う buf や未設定の mise シムでは生成結果が CI とずれるため、版が一致するときだけ PATH の buf を使う
BUF ?= $(shell [ "$$(buf --version 2>/dev/null)" = "$(BUF_VERSION)" ] && echo buf || echo go run github.com/bufbuild/buf/cmd/buf@v$(BUF_VERSION))
UP_BUILD ?= --build
GO_SERVICES := user session idp project task
PY_SERVICES := column
PNPM := pnpm

.PHONY: gen lint up down clean e2e contract

# Requires: buf on PATH, or go (falls back to `go run`) and network access (buf remote plugins).
gen:
	rm -rf gen/go
	$(BUF) generate
	@# Python はサービスと同じ版の grpcio-tools で生成し、BSR のリモートプラグインに依存しない
	@for s in $(PY_SERVICES); do \
		rm -rf backend/$$s/gen && mkdir -p backend/$$s/gen && \
		uvx --from grpcio-tools==1.71.0 python -m grpc_tools.protoc -I proto \
			--python_out=backend/$$s/gen --grpc_python_out=backend/$$s/gen proto/$$s/$$s.proto || exit 1; \
	done

# Requires: go, docker (Go services are linted via their Dockerfile lint target).
lint:
	$(BUF) lint
	@for s in user session idp task; do $(MAKE) -C backend/$$s lint || exit 1; done
	cd backend/project && go vet ./...  # project has no Docker lint target
	@for s in $(PY_SERVICES); do $(MAKE) -C backend/$$s lint || exit 1; done

# 既知の署名鍵で起動させないため、JWT の鍵は生成する
.env:
	sed -e "s/^JWT_SECRET_KEY=$$/JWT_SECRET_KEY=$$(openssl rand -hex 32)/" \
		-e "s|^IDP_SIGNING_KEY=$$|IDP_SIGNING_KEY=$$(openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 2>/dev/null | base64 | tr -d '\n')|" \
		.env.example > $@

up: .env
	docker compose up -d $(UP_BUILD) --wait --wait-timeout 300

down: .env
	docker compose down

clean: .env
	docker compose down -v

# Requires: node and pnpm (mise locally, corepack in CI) and a running stack (`make up`).
e2e:
	cd e2e && $(PNPM) install --frozen-lockfile && $(PNPM) exec playwright install chromium && $(PNPM) test

# Contract tests for #65 services. CI runs only the services listed in e2e/contract-enabled.txt.
# Addresses come from env (defaults in e2e/tests/contract/lib/config.ts).
# SVC="idp project" で担当サービスの契約テストだけを実行する（ファイル名で絞り込む）
contract:
	cd e2e && $(PNPM) install --frozen-lockfile && $(PNPM) exec playwright install chromium && $(PNPM) test:contract $(SVC)
