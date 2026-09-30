.PHONY: all clean proto test server client

# ビルド出力先ディレクトリ
BIN_DIR := ./bin
# ビルド対象のGoソースコード
SERVER_SRC := ./cmd/server/main.go
CLIENT_SRC := ./cmd/client/main.go

# ビルド出力ファイル名
SERVER_BIN := $(BIN_DIR)/project-service
CLIENT_BIN := $(BIN_DIR)/project-client

# ビルド対象となるGoソースファイル
GO_FILES := $(shell find . -type f -name "*.go")

# デフォルトターゲット
all: $(SERVER_BIN) $(CLIENT_BIN)

# ディレクトリの作成
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# サーバーのビルド
$(SERVER_BIN): $(GO_FILES) | $(BIN_DIR)
	go build -o $(SERVER_BIN) $(SERVER_SRC)

# クライアントのビルド
$(CLIENT_BIN): $(GO_FILES) | $(BIN_DIR)
	go build -o $(CLIENT_BIN) $(CLIENT_SRC)

# protoファイルからGo言語のコードを生成
proto:
	@$(MAKE) -C ../.. gen

# テスト実行
test:
	go test ./...

# クリーン
clean:
	rm -rf $(BIN_DIR)

# サーバー起動（開発用）
server: $(SERVER_BIN)
	$(SERVER_BIN) server

# クライアント実行（開発用）
client: $(CLIENT_BIN)
	$(CLIENT_BIN)