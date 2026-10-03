#!/usr/bin/env bash
# 変更されたファイルから、走らせるジョブの種類を判定し、key=true|false を出力する。
# ワークフローは結果でジョブを if で飛ばす（skipped は必須チェックでも合格扱い）。fetch-depth 0 で checkout してから呼ぶ
set -euo pipefail

keys="go python ts proto e2e"
all() { for k in $keys; do echo "$k=true"; done; }

# main への push は、壊れていないことを毎回確かめるためすべて走らせる
if [ "${GITHUB_EVENT_NAME:-}" = push ] && [ "${GITHUB_REF:-}" = refs/heads/main ]; then
  all
  exit 0
fi

if [ -n "${CHANGED_FILES+x}" ]; then
  files=$CHANGED_FILES
else
  git fetch --no-tags --quiet origin main
  files=$(git diff --name-only "$(git merge-base HEAD origin/main)" HEAD)
fi
while read -r f; do echo "changed: $f" >&2; done <<<"$files"

has() { grep -qE "$1" <<<"$files"; }
flag() { if has "$1"; then echo "$2=true"; else echo "$2=false"; fi; }

# CI の定義やツールの版が変わったら、何が壊れるか分からないのですべて走らせる
if has '^(\.github/|mise\.toml$)'; then
  all
  exit 0
fi
flag '^(backend/(idp|project|task)/|bff/|pkg/|gen/go/|go\.(mod|sum)$)' go
flag '^backend/column/' python
flag '^(frontend/|package\.json$|pnpm-(lock|workspace)\.yaml$|biome\.json$)' ts
flag '^(proto/|gen/|backend/column/gen/|buf(\.gen)?\.yaml$|Makefile$)' proto
flag '^(backend/|bff/|frontend/|proto/|gen/|pkg/|e2e/|docker-compose\.yaml$|Makefile$|go\.(mod|sum)$|\.env\.example$|package\.json$|pnpm-(lock|workspace)\.yaml$)' e2e
