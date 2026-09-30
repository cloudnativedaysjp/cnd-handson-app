# Role Service

CloudNative Days Handson用のロールサービスのサンプルアプリケーションです。

---

## Features

- ロール作成
- ロール更新
- ロール情報取得
- ロール削除

---

## Requirements

- Python 3.12
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
## Quick Start

### ロール作成
#### コマンド
```bash
docker compose exec role-service python3 cmd/client/main.py create <name> <description>
```
#### 例
```bash
$ docker compose exec role-service python3 cmd/client/main.py create test test用のrole
Response from server: id=b99e3afb-dad5-4067-9c3d-883faf43ae04, name=test, description=test用のrole
```
### ロール更新
#### コマンド
```bash
docker compose exec role-service python3 cmd/client/main.py update <id> <name> <description>
```
#### 例
```bash
$ docker compose exec role-service python3 cmd/client/main.py update b99e3afb-dad5-4067-9c3d-883faf43ae04 update updateのtest用のrole
Response from server: id=b99e3afb-dad5-4067-9c3d-883faf43ae04, name=update, description=updateのtest用のrole
```

### ロール情報取得
```bash
docker compose exec role-service python3 cmd/client/main.py get <id>
```
#### 例
```bash
$ docker compose exec role-service python3 cmd/client/main.py get b99e3afb-dad5-4067
-9c3d-883faf43ae04 
Response from server: id=b99e3afb-dad5-4067-9c3d-883faf43ae04, name=update, description=updateのtest用のrole
```

### ユーザ削除
```bash
docker compose exec role-service python3 cmd/client/main.py delete <id>
```

#### 例
```bash
$ docker compose exec role-service python3 cmd/client/main.py delete b99e3afb-dad5-4
067-9c3d-883faf43ae04 
Response from server: Role with id b99e3afb-dad5-4067-9c3d-883faf43ae04 deleted successfully
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
│   └── role/           # ユーザー関連
├── pkg/                # 再利用可能なパッケージ
│   └── db/             # データベース関連
├── proto/              # proto関連
├── Dockerfile          # Dockerビルド設定
└── README.md           # このファイル
```

---



