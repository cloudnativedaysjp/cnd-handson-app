---
name: start-issue
description: Issue 番号を受け取り、worktree 作成から実装・make e2e・PR 作成まで進める
---

引数: Issue 番号 `{n}`

1. `gh issue view {n}` で読む。親 Issue があればそれも読む
2. `feat/{n}-{kebab-slug}` の worktree とブランチを作る (slug は英語 max 3 語)
3. AGENTS.md に従って実装する (proto first)
4. `make e2e` を通す。失敗をテスト側の削除・skip で回避しない
5. Conventional Commits (英語) でコミットする。ファイルは 1 つずつ `git add` する
6. 人間に push の確認を取る。承認後のみ push
7. `gh pr create` で PR を作る。`.github/pull_request_template.md` に沿って Issue へのリンクと `make e2e` の結果を書き、`--label ai-authored` を付ける
