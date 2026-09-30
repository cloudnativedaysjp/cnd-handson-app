# Task Service

CloudNative Days Handson用のタスクサービスのサンプルアプリケーションです。

---

## Features

- タスク作成
- タスク更新
- タスク情報取得
- タスク削除
- タスク一覧

---

## Requirements

- Go 1.24.1 以上
- PostgreSQL

---

## Setup

### 1. リポジトリをクローン
```bash
git clone https://github.com/cloudnativedaysjp/cnd-handson-app/
cd cnd-handson-app/backend/task
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
binにtask-serviceバイナリが作成されます。
## Quick Start
### タスクの作成
#### コマンド
```bash
go run cmd/client/main.go create-task <title> <description> <status> <column_id> <Assgnee_id>
```

#### 例
```bash
go run cmd/client/main.go create-task "Task Title" "Task Description" "TODO" 9d1e1d77-113a-8f18-7737-452806dbd3e5 1683b280-6d67-1753-e1f6-c6a9932f2a51
Task created: id:"1fd35ddd-7793-40cf-bf27-31a2937941b8"  title:"Task Title"  description:"Task Description"  status:"TODO"  start_time:{seconds:1747614680  nanos:738533000}  end_time:{seconds:1747614680  nanos:738533966}  column_id:"9d1e1d77-113a-8f18-7737-452806dbd3e5"  assignee_id:"1683b280-6d67-1753-e1f6-c6a9932f2a51"
```

### タスクの更新
#### コマンド
```bash
go run cmd/client/main.go update-task -id -title -description -status -column_id -assgnee_id
```
#### 例
```bash
$ go run cmd/client/main.go update-task -id 1fd35ddd-7793-40cf-bf27-31a2937941b8 -title "Updated Task Title" -description "Updated Task Description" -status "IN_PROGRESS"
Task updated: id:"1fd35ddd-7793-40cf-bf27-31a2937941b8"  title:"Updated Task Title"  description:"Updated Task Description"  status:"IN_PROGRESS"  start_time:{seconds:1747614680  nanos:738533000}  end_time:{seconds:1747614819  nanos:189074002}  column_id:"9d1e1d77-113a-8f18-7737-452806dbd3e5"  assignee_id:"1683b280-6d67-1753-e1f6-c6a9932f2a51"
```
### タスクの取得
#### コマンド
```bash
go run cmd/client/main.go get-task <task_id>
```
#### 例
```bash
$ go run cmd/client/main.go get-task 1fd35ddd-7793-40cf-bf27-31a2937941b8
Task: id:"1fd35ddd-7793-40cf-bf27-31a2937941b8"  title:"Updated Task Title"  description:"Updated Task Description"  status:"IN_PROGRESS"  start_time:{seconds:1747614680  nanos:738533000}  end_time:{seconds:1747614819  nanos:189074000}  column_id:"9d1e1d77-113a-8f18-7737-452806dbd3e5"  assignee_id:"1683b280-6d67-1753-e1f6-c6a9932f2a51"
```

### タスクの一覧表示
#### コマンド
```bash
go run cmd/client/main.go list-tasks <column_id> <assignee_id> <page> <page_size>
```
#### 例
```bash
$ go run cmd/client/main.go list-tasks 9d1e1d77-113a-8f18-7737-452806dbd3e5 1683b280-6d67-1753-e1f6-c6a9932f2a51 1 10
タスク一覧: 
ID: 1fd35ddd-7793-40cf-bf27-31a2937941b8, Title: Updated Task Title, Description: Updated Task Description, Status: IN_PROGRESS
Total Count: 1
```

### タスクの削除
#### コマンド
```bash
go run cmd/client/main.go delete-task <task_id>
```
#### 例
```bash
$ go run cmd/client/main.go delete-task 1fd35ddd-7793-40cf-bf27-31a2937941b8
削除成功: true
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
├── .devcontainer/      # DevContainer設定
├── Dockerfile          # Dockerビルド設定
├── docker-compose.yml  # Docker Compose設定
├── Makefile            # ビルド・テスト用Makefile
└── README.md           # このファイル
```

---



