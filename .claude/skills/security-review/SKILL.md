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
- BFF の認証結果を backend が無条件に信用していないか
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

最後に件数と reviewed areas をまとめる。

finding がない場合は:

`No exploitable security findings found.`
