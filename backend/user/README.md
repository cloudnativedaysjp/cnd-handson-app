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

リポジトリ直下の `.env.example` を `.env` にコピーし、リポジトリ直下で起動します（docker compose を使用）。

```bash
cp .env.example .env
make up     # 起動
make down   # 停止
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
User created: id:"55c9ed73-ec88-4ffc-8458-8e2f2ef52906" name:"test" email:"test@cloudnative.com" role_id:"00000000-0000-0000-0000-000000000000" created_at:{seconds:1745085211 nanos:51603000} updated_at:{seconds:1745085211 nanos:51604255}
```
### ユーザー更新
#### コマンド
```bash
./bin/user-client update-user <id> <name> <email> <password> <role_id>
```
#### 例
```bash
$ ./bin/user-client update-user "55c9ed73-ec88-4ffc-8458-8e2f2ef52906" "" "" "" "123e4567-e89b-12d3-a456-426614174000" 
User updated: id:"55c9ed73-ec88-4ffc-8458-8e2f2ef52906" name:"test" email:"test@cloudnative.com" role_id:"123e4567-e89b-12d3-a456-426614174000" created_at:{seconds:1745085211 nanos:51603000} updated_at:{seconds:1745085530 nanos:281918000}
```

### ユーザ情報取得
```bash
./bin/user-client get-user <id>
```
#### 例
```bash
$ ./bin/user-client get-user b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
User: id:"b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d" name:"test" email:"test1@cloudnative.com" role_id:"00000000-0000-0000-0000-000000000000" created_at:{seconds:1745086235 nanos:825482000} updated_at:{seconds:1745086235 nanos:825483000}
```

### ユーザ削除
```bash
./bin/user-client delete-user <id>
```

#### 例
```bash
$ ./bin/user-client delete-user b8228ab4-ff39-4568-b8f9-2eb1cb3cd59d
削除成功: true
```

### パスワード検証
```bash
./bin/user-client verify-password <email> <password>
```

#### 例
```bash
$ ./bin/user-client verify-password test@cloudnative.com test
Password valid: true
```
---

## gRPC Documentation

protoからのコード生成は、リポジトリ直下で `make gen` を実行してください。


## Project Structure

```
.
├── cmd/                # エントリーポイント
│   ├── client/         # clientのエントリーポイント
│   └── server/         # serverのエントリーポイント
├── internal/           # 内部ロジック
│   └── user/           # ユーザー関連
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



