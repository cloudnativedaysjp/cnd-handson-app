# Session Service

CloudNative Days Handson用のsessionサービスのサンプルアプリケーションです。

---

## Features

- アクセストークンの発行
- アクセストークンの検証
- リフレッシュトークンの発行
- リフレッシュトークンの検証
- リフレッシュトークンの無効化（logout）

---

## Requirements

- Go 1.24.1 以上
- PostgreSQL

---

## Setup

### 1. リポジトリをクローン
```bash
git clone https://github.com/cloudnativedaysjp/cnd-handson-app/
cd session-service
```

### 2. 環境変数を設定
リポジトリ直下の `.env.example` を `.env` にコピーし、以下のように設定してください。

```env
DB_HOST=db
DB_PORT=5432
DB_DB=your_db_name
DB_USER=your_db_user
DB_PASSWORD=your_db_password
JWT_SECRET_KEY=your_jwt_secret_key
```

### 3. マイグレーションを実行
```bash
go run cmd/server/main.go migrate
```

### 4. アプリケーションを起動
```bash
go run cmd/server/main.go serve
```
## AppのBuild
```bash
make
```
binにsession-serviceバイナリが作成されます。
## Quick Start
### Access Tokenの作成
#### コマンド
```bash
./bin/session-client generate-access-token <user_id>
```
#### 例
```bash
$ ./bin/session-client generate-access-token b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
AccessToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDUyMzYwNzAsImlkIjoiYjgyMjhhYjQtZmYzOS00NTY4LWI4ZjktMmViMWNiM2NkNTlkIn0.4uahAE_pPug8jEEnKRAQC3hysWAfl2nTX_rKokl7kB8
ExpiresAt: 1745236070
```
### Access Tokenの検証
#### コマンド
```bash
./bin/session-client verify-access-token <access_token>
```
#### 例
```bash
./bin/session-client validate-access-token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDUyMzYwNzAsImlkIjoiYjgyMjhhYjQtZmYzOS00NTY4LWI4ZjktMmViMWNiM2NkNTlkIn0.4uahAE_pPug8jEEnKRAQC3hysWAfl2nTX_rKokl7kB8
Valid: true
UserID: b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
Error: 
```
### Refresh Tokenの作成
#### コマンド
```bash
./bin/session-client generate-refresh-token <user_id>
```

#### 例
```bash
$ ./bin/session-client generate-refresh-token b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
RefreshToken: t3sBwEzVj7d9JJFY-3fbZ-Eh8psi0MfQd0xVVGDLpTQ=
ExpiresAt: 1747827790
```
### Refresh Tokenの検証
#### コマンド
```bash
./bin/session-client verify-refresh-token <refresh_token> <user_id>
```
#### 例
```bash
$ ./bin/session-client validate-refresh-token t3sBwEzVj7d9JJFY-3fbZ-Eh8psi0MfQd0xVVGDLpTQ= b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
Valid: true
```
### Refresh Tokenの無効化
#### コマンド
```bash
./bin/session-client revoke-refresh-token <refresh-token> <user-id>
```

#### 例
```bash
$ ./bin/session-client revoke-refresh-token t3sBwEzVj7d9JJFY-3fbZ-Eh8psi0MfQd0xVVGDLpTQ= b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
Success: true
```


## gRPC Documentation

gRPCの詳細な仕様は、[gRPC仕様書](proto/READEME.md)を参照してください。


## Project Structure

```
.
├── cmd/                # エントリーポイント
│   ├── client/         # clientのエントリーポイント
│   └── server/         # serverのエントリーポイント
├── internal/           # 内部ロジック
│   ├── auth/           # トークン関連
├── pkg/                # 再利用可能なパッケージ
│   ├── auth/           # 認証関連
│   └── db/             # データベース関連
├── proto/              # proto関連
├── Dockerfile          # Dockerビルド設定
├── docker-compose.yml  # Docker Compose設定
├── Makefile            # ビルド・テスト用Makefile
└── README.md           # このファイル
```

---



