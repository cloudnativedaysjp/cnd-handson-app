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

### 1. リポジトリをクローン
```bash
git clone https://github.com/cloudnativedaysjp/cnd-handson-app/
cd role-service
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
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/server/main.py migrate
```
もし失敗するorテーブルをリセットしたい場合
```bash
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/server/main.py reset
```

### 5. アプリケーションを起動
```bash
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/server/main.py server
```
---
## Quick Start

### ロール作成
#### コマンド
```bash
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py create <name> <description>
```
#### 例
```bash
$ PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py create test test用のrole
Response from server: id=b99e3afb-dad5-4067-9c3d-883faf43ae04, name=test, description=test用のrole
```
### ロール更新
#### コマンド
```bash
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py update <id> <name> <description>
```
#### 例
```bash
$ PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py update b99e3afb-dad5-4067-9c3d-883faf43ae04 update updateのtest用のrole
Response from server: id=b99e3afb-dad5-4067-9c3d-883faf43ae04, name=update, description=updateのtest用のrole
```

### ロール情報取得
```bash
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py get <id>
```
#### 例
```bash
$ PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py get b99e3afb-dad5-4067
-9c3d-883faf43ae04 
Response from server: id=b99e3afb-dad5-4067-9c3d-883faf43ae04, name=update, description=updateのtest用のrole
```

### ユーザ削除
```bash
PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py delete <id>
```

#### 例
```bash
$ PYTHONPATH=/workspaces/role python3 /workspaces/role/cmd/client/main.py delete b99e3afb-dad5-4
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
├── .devcontainer/      # DevContainer設定
├── Dockerfile          # Dockerビルド設定
└── README.md           # このファイル
```

---



