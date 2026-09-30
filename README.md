# cnd-handson-app
『一日で学ぶクラウドネイティブ技術実践ハンズオン』by CloudNative Days 実行委員会 アプリケーションリポジトリ

## 参加のしかた

1. [#65](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/65) のサブ Issue から担当を選び、自分を Assignee にする
2. 自分のエージェント（Claude Code なら `/start-issue <番号>`、ほかのツールは `AGENTS.md` を読ませる）に Issue を渡す
3. 完了条件は Issue に書いてある。基本は `make up && make contract SVC=<svc>` が通ること
4. PR を出す。CodeRabbit の指摘は直すか、理由を返信して resolve する（resolve しないとマージできない）
5. マージ後、`e2e/contract-enabled.txt` に書いたサービスの契約テストは CI で守られる

セットアップ: `mise trust && mise install && make .env && make up && make e2e`
