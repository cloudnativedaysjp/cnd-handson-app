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
| `IDP_SIGNING_KEY` | アクセストークン（RS256）の署名鍵。base64 の PEM（PKCS#1 / PKCS#8）。`make .env` が生成する |
| `IDP_ISS` / `IDP_AUD` | JWT の `iss` / `aud` |

アクセストークンは 15 分、リフレッシュトークンは 30 日有効。`Refresh` は使ったリフレッシュトークンを無効にし、新しい組を返す。
