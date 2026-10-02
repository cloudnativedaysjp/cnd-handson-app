# idp

ユーザー登録・ログインとトークン発行を担う gRPC サービス（`handson-idp`）。user / session / role を置き換える（#104）。

- proto: `proto/idp/idp.proto`
- 契約テスト: `e2e/tests/contract/idp.spec.ts`

## コマンド

```bash
idp-service migrate   # テーブルを作る
idp-service server    # gRPC を PORT（既定 50051）で起動
```

## 環境変数

| 変数 | 用途 |
|---|---|
| `PORT` | gRPC のポート（既定 50051） |
| `DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_DB` | Postgres |
| `JWT_SECRET_KEY` | アクセストークンの署名鍵 |
