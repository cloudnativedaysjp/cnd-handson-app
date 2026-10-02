---
name: security-review
description: cnd-handson-app の変更を、実際に悪用可能な問題に絞ってセキュリティレビューする
---

`git diff origin/main...HEAD` を起点にセキュリティレビューする。

diff だけで判断せず、必要に応じて呼び出し元・呼び出し先・proto・repository・認証処理まで追跡する。
`gen/` は原則レビューせず、元の `proto/` を確認する。

## Priority

### Authentication / session

- JWT の生成・署名・検証・期限
- refresh token の生成・保存・検証・失効
- password hashing
- authentication bypass
- token / password / secret の漏洩

### Authorization

最優先で確認する。

resource ID を受け取る処理では、

authenticated user
→ membership / role
→ project / board / task 等の resource

まで辿り、そのユーザーに操作権限があることを確認する。

認証済みであることを authorization の代わりにしない。
IDOR / BOLA を重点的に確認する。

### Service boundaries

Frontend → BFF → gRPC services → DB の各境界を確認する。

内部 gRPC だからという理由だけで入力を trusted とみなさない。

特に以下を確認する。

- user ID / role を request から偽装できないか
- BFF での外部認証情報の検証から backend での principal（認証済みユーザー情報）の受領まで追跡し、backend が認証済み BFF から、クライアントや通信経路上の攻撃者が偽装・改ざんできない principal を受け取ることを確認する
- client が任意に指定した未検証の user ID / role を、認可の主体や権限の根拠に使っていないか
- service 間で authorization context が失われていないか

### Untrusted input

外部入力から sensitive operation まで data flow を追う。

SQL、shell、filesystem、URL、HTML、redirect、deserialization 等に到達する場合は、
injection、path traversal、SSRF、XSS 等の可能性を確認する。

危険そうな API が存在するだけでは finding にしない。

## Finding criteria

finding を報告する前に以下を確認する。

1. attacker-controlled な入力または状態がある
2. sensitive operation / resource まで到達できる
3. validation / authentication / authorization で防止されていない
4. 具体的な security impact がある

best practice の不足だけでは finding にしない。
推測だけの finding を出さない。

## Output

finding ごとに以下を出す。

### [Critical|High|Medium|Low] Title

`path:line`

- Attack path:
- Impact:
- Why exploitable:
- Fix:

最後に件数と reviewed areas、unverified areas をまとめる。
unverified areas には未確認領域と確認できなかった理由を記載し、未確認領域がなければ `none` と明記する。

finding がない場合は:

`No exploitable security findings found in the reviewed areas.`

`Unverified areas: [未確認領域と理由。なければ none].`
