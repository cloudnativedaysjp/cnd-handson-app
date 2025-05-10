# Project Service

CloudNative Days Handson用のプロジェクトサービスのサンプルアプリケーションです。

## 概要

このプロジェクトサービスは、カンバン型タスク管理アプリケーションのバックエンドを構成する一部として、プロジェクト管理の機能を提供します。gRPCを使用したマイクロサービスとして実装されており、PostgreSQLデータベースを使用してデータを保存します。

---

## 機能

- プロジェクト作成
- プロジェクト更新
- プロジェクト情報取得
- プロジェクト一覧取得
- オーナー別プロジェクト一覧取得
- プロジェクト削除

---

## 実行環境

このサービスは以下の環境で実行できるように設計されています：

### ローカル開発環境

1. ローカルのGoランタイムで直接実行
2. Docker Composeを使用した開発環境

### 本番環境

1. Kubernetes (k8s) クラスタ上での実行

---

## 必要条件

- Go 1.24.3 以上
- PostgreSQL

---

## 環境セットアップ

### Go のインストール

```bash
# バージョン確認
go version

# インストール（Homebrewの場合）
brew install go@1.24.3
```

### 環境変数

サービスの設定は環境変数で行います：

| 環境変数 | 説明 | デフォルト値 |
|----------|------|-------------|
| DB_HOST | データベースホスト名 | localhost |
| DB_PORT | データベースポート | 5432 |
| DB_USER | データベースユーザー名 | postgres |
| DB_PASSWORD | データベースパスワード | postgres |
| DB_DB | データベース名 | postgres |
| PORT | gRPCサーバーポート | 50053 |
| GRPC_HOST | gRPCクライアント用ホスト | localhost |
| GRPC_PORT | gRPCクライアント用ポート | 50053 |

---

## セットアップと実行

### 1. リポジトリのクローン

```bash
git clone https://github.com/cloudnativedaysjp/cnd-handson-app/
cd cnd-handson-app/backend/project
```

### 2. 環境変数の設定

`.env` ファイルを作成し、必要な環境変数を設定：

```env
DB_PASSWORD=postgres
```

### 3. データベースマイグレーション

```bash
# プロジェクトルートディレクトリで実行
go run cmd/server/main.go migrate
```

マイグレーションに失敗した場合やテーブルをリセットしたい場合：

```bash
go run cmd/server/main.go reset
```

### 4. サーバーの起動

```bash
go run cmd/server/main.go server
```

### 5. クライアントからのテスト実行

別のターミナルで以下を実行：

```bash
# テストスクリプトを使用（全機能テスト）
./scripts/test-client.sh

# または個別コマンドで実行
go run cmd/client/main.go create-project "新しいプロジェクト" "これは新しいプロジェクトです" "123e4567-e89b-12d3-a456-426614174000"
go run cmd/client/main.go list-projects
```

---

## Docker Composeでの実行

docker-compose.yamlファイルにはプロジェクトサービス用の設定が含まれています：

```bash
# Docker Composeでビルド・起動
docker-compose up --build

# バックグラウンドで起動
docker-compose up -d

# サービスの停止
docker-compose down
```

---

## ビルド

### ローカルでビルド

```bash
make
```

---

## APIドキュメント

gRPC APIの詳細は[proto/README.md](proto/README.md)を参照してください。

### 利用可能なクライアントコマンド

- `create-project <name> <description> <owner_id>`: プロジェクトを作成
- `update-project <id> <name> <description>`: プロジェクトを更新
- `get-project <id>`: プロジェクト情報を取得
- `list-projects [owner_id]`: プロジェクト一覧を取得（オーナーID指定は任意）
- `delete-project <id>`: プロジェクトを削除

---

## テストの実行

```bash
# 全テストの実行
go test ./...
```

---

## プロジェクト構造

```
.
├── cmd/                # エントリーポイント
│   ├── client/         # クライアントのエントリーポイント
│   └── server/         # サーバーのエントリーポイント
├── internal/           # 内部ロジック
│   └── project/        # プロジェクト関連
│       ├── handler/    # gRPCハンドラー
│       ├── model/      # データモデル
│       ├── repository/ # データベースアクセス層
│       └── service/    # ビジネスロジック層
├── pkg/                # 再利用可能なパッケージ
│   └── db/             # データベース関連
├── proto/              # proto関連
├── scripts/            # テスト用スクリプト
├── Dockerfile          # Dockerビルド設定
├── Makefile            # ビルド・テスト用Makefile
└── README.md           # このファイル
```