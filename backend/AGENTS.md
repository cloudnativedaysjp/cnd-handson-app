# backend/AGENTS.md

- OTel 計装は必須 (gRPC サーバーの trace / metrics)。送信先は `OTEL_EXPORTER_OTLP_ENDPOINT`
- ログは JSON、`trace_id` を含める
- 新サービスは `backend/project` と同じ構成にする
  - `cmd/server/main.go`: 起動
  - `internal/<name>/{handler,service,repository,model}`: 層ごとに分離。テストは `service/test/`
  - `pkg/db`: DB 接続とマイグレーション
  - `Dockerfile`、`go.mod` はサービスごと
- proto は root の `proto/` で管理。生成コードは編集しない
