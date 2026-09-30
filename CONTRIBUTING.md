# 参加のしかた

## セットアップ

```bash
mise trust && mise install && make .env && make up && make e2e
```

## 進め方

1. [#65](https://github.com/cloudnativedaysjp/cnd-handson-app/issues/65) などのサブ Issue から担当を選び、自分を Assignee にする
2. 自分のエージェントに Issue を渡す（Claude Code なら `/start-issue <番号>`、ほかのツールは `AGENTS.md` を読ませる）
3. Issue の完了条件を満たす。サービスの実装なら `make up && make contract SVC=<svc>` が通ること
4. PR を出す。CodeRabbit の指摘は直すか、理由を返信して resolve する（resolve しないとマージできない）
5. 実装が終わったサービスは `e2e/contract-enabled.txt` に追加する。以後、そのサービスの契約テストは CI の必須チェックで守られる

## ルール

作業のルール（proto first、完了条件、コミット、禁止事項）は [`AGENTS.md`](AGENTS.md)、サービスの約束（計測・実行時）は [`docs/conventions.md`](docs/conventions.md) にまとめている。人もエージェントも同じルールで進める。

新しい Issue は、テンプレート（サービス実装・章の移行・改善と整備）から作る。
