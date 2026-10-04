# cnd-handson-app

『一日で学ぶクラウドネイティブ技術実践ハンズオン』（[cnd-handson](https://github.com/cloudnativedaysjp/cnd-handson)）で使うデモアプリです。CloudNative Days 実行委員会が開発しています。

## どんなアプリか

マイクロサービスで作ったタスク管理（カンバン）です。フロントエンドには 2 つの版があります。ハンズオンの canary やトラフィック分割の実演で、どちらに届いたかを見分けられます。

| 版 | 画面 | color |
|---|---|---|
| legacy | 社内システム風の表形式 | blue |
| modern | カンバンボード | green |

どちらの版も、`/color` で自分の color を返します。

## 構成

```
ブラウザ → handson-legacy / handson-modern → project → task   → Postgres
                 │                     └→ column → Postgres
                 └→ idp（ログイン） → Postgres
```

| サービス | 言語 | 役割 |
|---|---|---|
| フロントエンド（`bff/`、`frontend/`）。handson-legacy と handson-modern | Go、React | 画面と REST API。JWT を検証し、利用者の ID を下流に渡す |
| idp（`backend/idp`） | Go | ログインと JWT の発行 |
| project（`backend/project`） | Go | プロジェクト。タスクと状態の一覧をまとめる。所有者かどうかを確かめる |
| task（`backend/task`） | Go | タスク |
| column（`backend/column`） | Python | タスクの状態 |

- サービスどうしは gRPC で話します。proto は `proto/` にあります。
- どのサービスも、トレースとメトリクスを OpenTelemetry で出します。フロントエンドから DB まで、1 本のトレースでつながります。

## 手元で動かす

```bash
mise trust && mise install
make .env && make up   # 起動する。healthy になるまで待つ
make e2e               # 動作を確かめる
make down              # 止める。データは残る。消すなら make clean
```

起動したら、フロントエンドを開いてログインします。

- legacy: http://localhost:8080
- modern: http://localhost:8081
- ユーザー: `demo@example.com` / `demo-password`

http://localhost:5173 は、frontend の開発サーバーです。API に繋がらないので、ログインは 8080 か 8081 から行います。

### トレースを見る

```bash
COMPOSE_PROFILES=otel make up
```

http://localhost:3001 の Grafana を開き、Explore で Tempo を選ぶと、トレースが見られます。

## k8s に入れる

Helm chart を `deploy/helm/handson` に置いています。入れ方と設定は [deploy/README.md](deploy/README.md) を見てください。

## 今の状態

構築中です。

- ログインからタスクの表示までを 1 本通す作業は、[#65](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/65) で進めています。
- modern の画面は [#139](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/139) で作ります。
- ハンズオンの各章をこのアプリに移す作業は、[#66](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/66) です。

## ドキュメント

- [CONTRIBUTING.md](CONTRIBUTING.md): 参加のしかた
- [AGENTS.md](AGENTS.md): リポジトリの構成、コマンド、作業のルール
- [docs/conventions.md](docs/conventions.md): 計測と実行時の約束。サービスのコードを書く前に読む
- [docs/writing.md](docs/writing.md): Issue と PR の書き方
