---
name: サービス実装
about: 契約テストを通すサービスの実装
title: "<svc> サービスを実装する"
---

親: #65

## 触る範囲
- `backend/<svc>/`（実装）、`docker-compose.yaml`（サービスの追加）
- `proto/` は変更しない（変更が必要なら先に Issue で相談）

## ポートと env
- gRPC: 50051（ポート名 `grpc`）、HTTP: 8080（ポート名 `http`）
- `OTEL_SERVICE_NAME=handson-<svc>`、DB は `DB_*`、ユーザー ID は metadata の `x-user-id`

## 仕様
- proto: `proto/<svc>/`
- 契約テスト: `e2e/tests/contract/<svc>.spec.ts`
- 約束: `docs/conventions.md`

## 完了条件
- [ ] `make up && make contract SVC=<svc>` が通る
- [ ] `e2e/contract-enabled.txt` に `<svc>` を追加し、CI が通る
- [ ] `make lint`、単体テスト、`make e2e` が通る
