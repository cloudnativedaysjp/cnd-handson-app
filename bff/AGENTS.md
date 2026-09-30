# bff/AGENTS.md

- OTel 計装は必須 (HTTP サーバー、gRPC クライアント)。送信先は `OTEL_EXPORTER_OTLP_ENDPOINT`
- ログは JSON、`trace_id` を含める
- gRPC クライアントは root の `proto/` から生成したコードを使う。生成コードは編集しない
