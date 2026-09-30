# レビューの基準

- コメントは日本語で書く
- 基準はリポジトリ直下の `AGENTS.md` と `docs/conventions.md`。これに反する変更を優先して指摘する
- フォーマットは CI（biome / gofmt / black）が見るので指摘しない

## 重点的に見ること
- 生成コード（`*.pb.go`、`backend/*/gen/`）を手で編集していないか。proto を変えたら `make gen` の結果が含まれているか
- テストの削除・skip・期待値の変更が、PR 本文で理由付きで説明されているか
- 受け取った `context` を下流の gRPC・DB 呼び出しに渡しているか（`context.Background()` / `context.TODO()` を使っていないか）
- OTel の設定（endpoint、`service.name`）をコードに直書きしていないか。環境変数で渡しているか
- 設定値・接続情報・秘密値がコードや既定値に埋め込まれていないか
- gRPC health、SIGTERM での graceful shutdown、非 root 実行が守られているか
- `.env` やシークレットの値がコミット・ログ出力に含まれていないか
