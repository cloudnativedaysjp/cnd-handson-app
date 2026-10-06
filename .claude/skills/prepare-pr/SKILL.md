---
name: prepare-pr
description: 変更をレビュー可能な PR にまとめるときに使う。PR の分割と依存関係、検証結果、提出、stacked PR の追従を扱う。既存ブランチからの提出にも使う。
---

## レビューの単位を決める

- PR は一つの変更意図を判断できる単位にする。機械的にファイル数や行数で分割しない。
- 独立した変更は別 PR にする。依存する変更は下のブランチから分岐し、直前の PR のブランチを base にする。
- 各 PR に目的と依存先を記載する。Issue には PR の一覧と取り込み順を残す。
- Issue 全体が終わらない途中の PR では `Closes` を使わず `Refs #N` とする。

## 変更に対応する検証を行う

- コード変更は関連する lint・単体テストを実行する。サービス間の挙動を変える場合は `make up && make e2e` も実行する。
- サービスを実装したら `make contract SVC=<svc>` を実行し、`e2e/contract-enabled.txt` に登録する。CI は登録されたサービスだけを検証する。
- proto を変更したら `make gen` と互換性検査の結果を確認する。意図的な破壊的変更は理由と影響を PR に書き、`breaking-proto` ラベルで検査を省略することを明示する。検査を通すためだけに付けない。
- 認証・認可・外部入力・サービス境界を変える場合は [security-review](../security-review/SKILL.md) を使う。実行時の信頼境界は `docs/conventions.md` と照合する。
- ドキュメントだけの変更は参照先と記述の整合性を確認する。アプリのテストを一律には要求しない。
- ローカルで実行できない検証は理由を記録する。同じ変更に対する CI の実行結果を確認できれば、その結果を使ってよい。skip や未実行は成功と区別する。

## 提出する

- コミットは Conventional Commits、英語。変更対象を確認してファイルごとに add する。
- 提出前に base の更新を確認する。main への追従は `git merge origin/main`。生成コードの競合は再生成し、それ以外は人間に確認する。stack の追従は下記を参照する。
- AGENTS.md の許可範囲で作業ブランチを push する。許可がなければコミットまでで止め、push 待ちと報告する。
- `.github/pull_request_template.md` と `docs/writing.md` に従って PR を作り、`ai-authored` を付ける。ラベルがなければ assignee に報告する。
- 検証した内容と結果、未実行の理由、人間に判断してほしい点を記載する。必要な検証が未確認なら draft にする。
- CI の結果を確認して PR に反映する。必要な検証が通っても、assignee の許可なくマージしない。

## stacked PR を追従する

- 下の PR が squash merge されたら、次の PR の base と依存先の記述を更新する。
- 取り込み手順は [docs/stacked-pr.md](../../../docs/stacked-pr.md) に従う。`proto/`・`e2e/` の競合は人間に確認する。
- ツリーが変わった場合は影響する検証をやり直す。Issue の取り込み状況を更新する。
