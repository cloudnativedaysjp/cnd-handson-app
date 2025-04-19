# User Service

CloudNative Days Handson用のユーザサービスのサンプルアプリケーションです。

---

## Features

- ユーザ作成
- ユーザ更新
- ユーザ情報取得
- ユーザ削除
- パスワード検証

---

## Requirements

- Go 1.24.1 以上
- PostgreSQL

---

## Setup

### 1. リポジトリをクローン
```bash
git clone https://github.com/cloudnativedaysjp/cnd-handson-app/
cd user-service
```

### 2. devContainerの環境変数を設定
`./devcontainr/.env` ファイルを作成し、以下のように設定してください（`.env.dummy` を参考にしてください）。

```env
DB_HOST=db
DB_PORT=5432
DB_DB=your_db_name
DB_USER=your_db_user
DB_PASSWORD=your_db_password
JWT_SECRET_KEY=your_jwt_secret_key
```

### 3. devcontainer起動
```bash
devcontainer open
```

### 4. マイグレーションを実行
```bash
go run cmd/server/main.go migrate
```
もし失敗するorテーブルをリセットしたい場合
```bash
go run cmd/server/main.go reset
```

### 5. アプリケーションを起動
```bash
go run cmd/server/main.go server
```
---
## AppのBuild
```bash
make
```
binにuser-serviceバイナリが作成されます。
## Quick Start

### ユーザー作成
#### コマンド
```bash
./bin/user-client create-user <name> <email> <password> <role_id>
```
#### 例
```bash
$ ./bin/user-client create-user test test@gmail.com test 123e4567-e89b-12d3-a456-426614174000
User created: id:"55c9ed73-ec88-4ffc-8458-8e2f2ef52906" name:"test" email:"test@ss.com" role_id:"00000000-0000-0000-0000-000000000000" created_at:{seconds:1745085211 nanos:51603000} updated_at:{seconds:1745085211 nanos:51604255}
```
### ユーザー更新
#### コマンド
```bash
update-user <id> <name> <email> <password> <role_id>
```
#### 例
```bash
$ ./bin/user-client update-user "55c9ed73-ec88-4ffc-8458-8e2f2ef52906" "" "" "" "123e4567-e89b-12d3-a456-426614174000" 
User updated: id:"55c9ed73-ec88-4ffc-8458-8e2f2ef52906" name:"test" email:"test@ss.com" role_id:"123e4567-e89b-12d3-a456-426614174000" created_at:{seconds:1745085211 nanos:51603000} updated_at:{seconds:1745085530 nanos:281918000}
```

### ユーザ削除
```bash

```

### ユーザ情報取得
```bash

```

---

## API Documentation

APIの詳細な仕様は、[OpenAPI仕様書](proto/READEME.md)を参照してください。Redocを使用してブラウザで確認することもできます。



## Project Structure

```
.
├── cmd/                # エントリーポイント
│   └── server/
│       └── main.go     # メインアプリケーション
├── internal/           # 内部ロジック
│   ├── user/           # ユーザー関連
│   └── refresh/        # リフレッシュトークン関連
├── pkg/                # 再利用可能なパッケージ
│   ├── auth/           # 認証関連
│   └── db/             # データベース関連
├── proto/              # proto関連
├── .devcontainer/      # DevContainer設定
├── Dockerfile          # Dockerビルド設定
├── docker-compose.yml  # Docker Compose設定
├── Makefile            # ビルド・テスト用Makefile
└── README.md           # このファイル
```

---



