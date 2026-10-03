# サービスの約束

サービス（入口を含む）のコードを書く・変える前に読む。

## 計測の約束（入口を含む全サービス）
- OTLP は `http/protobuf`（4318）。`OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf`
- endpoint はコードに書かず SDK の env 解決に任せる。gRPC exporter（4317）は使わない
- `OTEL_SERVICE_NAME` = デプロイ名（`handson-legacy` / `handson-modern` / `handson-idp` / `handson-project` / `handson-task` / `handson-column`）。コードに書かず env で渡す
- `OTEL_RESOURCE_ATTRIBUTES` に `service.namespace=handson`。入口は `app.variant` / `app.color` も付ける
- 伝播: `OTEL_PROPAGATORS=tracecontext,baggage`。起動時に TextMapPropagator を設定する
- gRPC は otelgrpc の stats handler（server / client 両方）
- 受けた ctx を下流の gRPC・DB に必ず渡す（`context.Background()` / `TODO()` 禁止）
- DB は otelsql（or otelpgx）で計装。引数はスパンに入れない
- 期待するスパン: 入口 HTTP → project gRPC → task gRPC → DB
- HTTP は otelhttp、スパン名は `METHOD route`
- サンプリングは既定（parentbased_always_on）
- collector が落ちても起動・リクエストを止めない
- ログ: stdout に JSON 1 行（ファイル出力しない）
- ログのキー固定: `time`(RFC3339) `level`(小文字) `msg` `service` `trace_id` `span_id`
- 入口のリクエストログは `variant` `color` `method` `path` `status` `duration_ms` も
- gRPC サービスもリクエストごとに `rpc.method` `code` を 1 行
- メトリクス: `/metrics` をポート 9464（ポート名 `metrics`）で常に公開（Prometheus scrape）
- `OTEL_METRICS_EXPORTER=otlp` のとき OTLP 送信も併用（MeterProvider に reader を 2 つ）
- メトリクス名・属性は OTel semantic conventions（stable）に従い独自名を作らない

## 実行時の約束
- 設定は環境変数のみ（`PORT`、`DB_*`、`OTEL_*`）。秘密値も env で受ける
- 新規サービスの gRPC は 50051（ポート名 `grpc`）、HTTP 入口は 8080（ポート名 `http`）
- gRPC health（`grpc.health.v1`）を実装し SERVING を返す
- ユーザー ID は gRPC metadata の `x-user-id` で渡す。JWT の検証は入口が行い、下流のサービスは `x-user-id` を使う
- HTTP 入口は `/healthz`（`/color` をヘルスチェックに流用しない）
- SIGTERM で GracefulStop / Shutdown
- マイグレーションは `<svc> migrate` サブコマンドに分離し、server 起動時に実行しない
- コンテナ: 非 root、ベースイメージのタグ固定、1 イメージ 1 バイナリ
