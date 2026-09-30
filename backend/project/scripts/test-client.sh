#!/bin/bash
# 接続先（compose の公開ポート）
export PROJECT_ADDR=localhost:50053

# ランダムなUUID生成（オーナーID用）- macOS/Linuxの場合
OWNER_ID=$(uuidgen | tr '[:upper:]' '[:lower:]')

echo "=== プロジェクト作成テスト ==="
go run cmd/client/main.go create-project "テストプロジェクト" "これはテストプロジェクトです" $OWNER_ID
echo "作成したオーナーID: $OWNER_ID"

# 作成したプロジェクトのIDを取得（実際のシナリオでは保存されるID）
PROJECT_ID=$(go run cmd/client/main.go list-projects $OWNER_ID | grep "ID:" | head -1 | awk '{print $2}')

echo "=== プロジェクト取得テスト ==="
go run cmd/client/main.go get-project $PROJECT_ID

echo "=== プロジェクト更新テスト ==="
go run cmd/client/main.go update-project $PROJECT_ID "更新プロジェクト" "説明を更新しました"

echo "=== プロジェクト一覧テスト（全て） ==="
go run cmd/client/main.go list-projects

echo "=== プロジェクト一覧テスト（オーナー指定） ==="
go run cmd/client/main.go list-projects $OWNER_ID

echo "=== プロジェクト削除テスト ==="
go run cmd/client/main.go delete-project $PROJECT_ID

echo "=== 削除確認（一覧取得） ==="
go run cmd/client/main.go list-projects $OWNER_ID
