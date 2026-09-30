# Column Service

CloudNative Days Handson用のカラムサービスのサンプルアプリケーションです。

---

## Features

- カラム作成
- カラム更新
- カラム情報取得
- カラム削除

---

## Requirements

- Python 3.12
- PostgreSQL

---

## Setup

### 1. リポジトリをクローン
```bash
git clone https://github.com/cloudnativedaysjp/cnd-handson-app/
cd column-service
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
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/server/main.py migrate
```
もし失敗するorテーブルをリセットしたい場合
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/server/main.py reset
```

### 5. アプリケーションを起動
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/server/main.py server
```
---
## Quick Start

### カラム作成
#### コマンド
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py create <name> <board_id>
```
#### 例
```bash
$ PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py create test b99e3afb-dad5-4067-9c3d-883faf43ae04
Response from server: id=09dca8cc-8e84-4e31-a11d-9fbc51fc82b7, name=test, board_id=b99e3afb-dad5-4067-9c3d-883faf43ae04
```
### カラム更新
#### コマンド
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py update <id> <name> <board_id>
```
#### 例
```bash
$ PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py update 09dca8cc-8e84-4e31-a11d-9fbc51fc82b7 updateのtest用のcolumn b99e3afb-dad5-4067-9c3d-883faf43ae04
.Response from server: id=09dca8cc-8e84-4e31-a11d-9fbc51fc82b7, name=updateのtest用のcolumn, board_id=b99e3afb-dad5-4067-9c3d-883faf43ae04
```

### カラム情報取得
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py get <id>
```
#### 例
```bash
$ PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py get 09dca8cc-8e84-4e31-a11d-9fbc51fc82b7
Response from server: id=09dca8cc-8e84-4e31-a11d-9fbc51fc82b7, name=updateのtest用のcolumn, board_id=b99e3afb-dad5-4067-9c3d-883faf43ae04
```

### カラム一覧取得
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py list　<board_id> <page> <page_size>
```

#### 例
```bash
$ PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py list  b99e3afb-dad5-4067-9c3d-883faf43ae04 1 10
Response from server: columns {
  id: "5e4a870a-6870-47f0-bcd6-6a4c3622200e"
  name: "test"
  board_id: "b99e3afb-dad5-4067-9c3d-883faf43ae04"
}
total_count: 1
```

### カラム削除
```bash
PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py delete <id>
```

#### 例
```bash
$ PYTHONPATH=/workspaces/column python3 /workspaces/column/cmd/client/main.py delete 09dca8cc-8e84-4e31-a11d-9fbc51fc82b7
Response from server: Column with id 09dca8cc-8e84-4e31-a11d-9fbc51fc82b7 deleted successfully
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
│   └── column/           # カラム関連
├── pkg/                # 再利用可能なパッケージ
│   └── db/             # データベース関連
├── proto/              # proto関連
├── .devcontainer/      # DevContainer設定
├── Dockerfile          # Dockerビルド設定
└── README.md           # このファイル
```

---



