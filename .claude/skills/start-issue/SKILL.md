---
name: start-issue
description: Issue や曖昧な開発依頼から着手するときに使う。担当・既存作業・目的を確認し、Issue に判断の履歴を残して作業環境を準備する。
---

## 目的と担当を揃える

- Issue があれば本文と親 Issue、関連 PR を読む。依頼だけなら関連 Issue を探し、なければ目的・完了の判断基準・未決事項をまとめて Issue に残す。
- `gh issue view {n} --json assignees`、`gh pr list --search "{n} in:body"`、`git ls-remote --heads origin "feat/{n}-*"` で既存の担当と作業を確認する。
- 自分たちの作業の続きなら既存のブランチ・PR を使う。他人の担当や並行作業と重なるなら、着手前に assignee と調整する。
- assignee は、作業の責任を持つ人にする。Issue を作った人とは限らない。誰が持つか分からなければ確認する。
- 実装を調べて方針を具体化し、Issue に着手と方針を残す。方針を変えた場合も判断理由を追記する。逐次の作業ログは不要。

## 作業環境を準備する

- 新規作業は `git fetch origin` 後、`origin/main` から `feat/{n}-{slug}` を作る。slug は英語で最大 3 語。
- worktree は `.claude/worktrees/{n}-{slug}` に作り、以後はその中で作業する。依存する変更は `prepare-pr` の分割方針に従い、依存先のブランチを起点にする。
- 必要なツールを `mise trust && mise install` で準備する。JS を変更する場合は pnpm で依存を入れる。
- Compose を使う場合は `make .env` で署名鍵を含む設定を生成する。既存の `.env` は上書きせず、値を出力しない。

実装は AGENTS.md と対象ディレクトリの約束に従う。検証と PR 提出は [.claude/skills/prepare-pr/SKILL.md](../prepare-pr/SKILL.md) を参照する。
