# AGENTS.md

CloudNative Days ハンズオン用カンバンのマイクロサービスデモ。

## 構成
- `frontend/`: Vite/React (pnpm)
- `bff/`: Go。frontend と gRPC サービスの間
- `backend/{user,session,project,task}`: Go gRPC
- `backend/{role,column}`: Python gRPC
- `proto/`: 全サービスの proto (buf)。唯一の定義元
- `e2e/`: Playwright と契約テスト
- Postgres は docker compose で起動

## コマンド
- `mise install`: ツール導入 (go / python / node / pnpm / buf)
- `make gen`: proto から生成 (buf generate)
- `make lint`
- `make up` / `make down`: docker compose の起動 (healthy まで待つ) / 停止
- `make e2e`: Playwright + backend smoke

## 作業ルール
- proto first: proto を変更 → `make gen` → 実装
- 生成コードは編集しない
- 完了 = `make e2e` が通ること
- テストを削除・skip・弱めて通さない
- `proto/` と契約テスト (`e2e/`) の変更は人間の承認が必要 (CODEOWNERS)
- コミットは Conventional Commits、英語
- 履歴を書き換えない (rebase / amend / force push / reset --hard 禁止)

## 運用
- Issue の assignee = エージェントに依頼した人。動作確認とレビュー対応はその人が持つ
- ブランチ名: `feat/{issue番号}-{kebab-slug}`
- PR は `.github/pull_request_template.md` に沿い、ラベル `ai-authored` を付ける
- 同じ間違いが 2 回起きたら、コードだけでなく AGENTS.md を直す

## 入口サービス (handson-legacy / handson-modern)
- ビルド時の `VARIANT` で切替: `legacy` → color `blue`、`modern` → color `green`
- `/` で画面、`/color` で color を返す
- リクエストごとに VARIANT と color を含むログを 1 行出す
- OTel SDK を組み込み、`OTEL_EXPORTER_OTLP_ENDPOINT` に送る
