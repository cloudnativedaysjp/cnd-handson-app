# cnd-handson-app

『一日で学ぶクラウドネイティブ技術実践ハンズオン』（[cnd-handson](https://github.com/cloudnativedaysjp/cnd-handson)）で使うデモアプリです。CloudNative Days 実行委員会が開発しています。

## どんなアプリか

マイクロサービスで作ったカンバンです。入口には 2 つの版があり、ハンズオンの canary やトラフィック分割の実演で見分けられるようにしています。

- **legacy**：表形式の画面（color: blue）
- **modern**：カンバンボード（color: green）

## 今の状態

構築中です。ログインからタスクの表示までを 1 本通す作業を [#65](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/65) で進めています。ハンズオンの各章への移行は [#66](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/66) です。

## 動かし方

```bash
mise trust && mise install
make .env && make up   # 起動（healthy になるまで待つ）
make e2e               # 動作確認
make down              # 停止（データは残る。消すなら make clean）
```

起動後、http://localhost:5173 で画面を確認できます（現在は開発中の frontend）。

トレースを見るときは `COMPOSE_PROFILES=otel make up` で起動します。http://localhost:3001 の Grafana を開き、Explore で Tempo を選ぶとトレースが見られます。

## 参加する方へ

[CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。
