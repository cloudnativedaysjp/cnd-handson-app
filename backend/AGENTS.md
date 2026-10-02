# backend/AGENTS.md

計測と実行時の約束は `docs/conventions.md`。

- 新サービスは `backend/project` と同じ構成にする
  - `cmd/server/main.go`: 起動
  - `internal/<name>/{handler,service,repository,model}`: 層ごとに分離。テストは `service/test/`
  - `pkg/db`: DB 接続とマイグレーション
  - `Dockerfile` はサービスごと（ビルドコンテキストはリポジトリルート）
- Go はルートの `go.mod` 1 つ。gRPC の stub は `gen/go/<svc>`（`make gen` で生成）
- 他サービスの `internal` は import しない。共有するのは `gen/go/<svc>` の stub と `pkg/`（計測は `pkg/telemetry`）だけ
- `pkg/` を使うサービスは Dockerfile で `COPY pkg ./pkg` する
