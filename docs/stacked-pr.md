# stacked PR を squash マージした後の追従

サブ Issue を順に積んだ PR（stacked PR）では、下の PR を squash マージすると、上の PR が毎回コンフリクトする。下の PR のコミットが、main では 1 つの別のコミットになるため。

上の PR は下の PR の変更をすべて含んでいる。なので、コンフリクトは上の PR 側を採れば解決できる。ただし、結果が正しいかを必ず確かめる。

## いつ使うか

次の 2 つがそろっているときだけ使う。

- 下の PR が squash マージされ、上の PR の base が main に変わった、または下のブランチが rebase された
- 上の PR のブランチに、下の PR の変更がすべて入っている

`proto/` か `e2e/` がコンフリクトしたときは、この手順を使わずに止めて、人間に確認する（AGENTS.md の作業ルール）。

## 手順

1. 上の PR のブランチで、今の HEAD を控える

   ```bash
   old=$(git rev-parse HEAD)
   ```

2. 新しい base を merge する。base が main なら `origin/main`、下のブランチが rebase されたなら `origin/<下のブランチ>` を使う。rebase と force push はしない

   ```bash
   git fetch origin
   git merge --no-commit --no-ff origin/main
   ```

   `--no-commit` を付ける。コンフリクトが無くても、手順 4 で確かめるまでコミットしないため

3. コンフリクトしたファイルは、上の PR 側を採る

   ```bash
   git diff --name-only --diff-filter=U | while IFS= read -r f; do
     git checkout --ours -- "$f" && git add -- "$f"
   done
   ```

4. 控えた HEAD とツリーを比べる

   ```bash
   git diff --cached --stat "$old"
   ```

   差分は、新しい base にだけある変更でなければならない。たとえば、別の PR で main に入ったファイルだけが出る。上の PR のファイルが出たら、取り込み方が間違っている。merge を `git merge --abort` で取り消し、人間に確認する

5. コミットして push する。ツリーが変わっていなければ、テストをやり直す必要はない

   ```bash
   git commit --no-edit
   git push
   ```

6. PR 本文の「Stacked on #N」や「#N の後にマージする」を消す
