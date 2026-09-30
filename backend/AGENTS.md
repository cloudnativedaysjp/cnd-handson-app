# backend/AGENTS.md

計測と実行時の約束はルート AGENTS.md。

- 新サービスは `backend/project` と同じ構成にする
  - `cmd/server/main.go`: 起動
  - `internal/<name>/{handler,service,repository,model}`: 層ごとに分離。テストは `service/test/`
  - `pkg/db`: DB 接続とマイグレーション
  - `Dockerfile`、`go.mod` はサービスごと
