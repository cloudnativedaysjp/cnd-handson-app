# deploy

アプリ一式を k8s に入れる Helm chart です。chart は `helm/handson` にあります。

## 入るもの

| 名前 | 役割 | ポート |
|---|---|---|
| `handson-legacy` | フロントエンド。表形式の画面と REST API（color: blue） | `http` 8080、`metrics` 9464 |
| `handson-modern` | フロントエンド。カンバンの画面と REST API（color: green） | `http` 8080、`metrics` 9464 |
| `handson-idp` | ログインと JWT の発行。JWKS を HTTP で返す | `grpc` 50051、`http` 8080、`metrics` 9464 |
| `handson-project` | プロジェクト。タスクと状態の一覧をまとめる | `grpc` 50051、`metrics` 9464 |
| `handson-task` | タスク | `grpc` 50051、`metrics` 9464 |
| `handson-column` | タスクの状態（Python） | `grpc` 50051、`metrics` 9464 |
| `handson-postgres` | 全サービスが使う DB | `postgres` 5432 |

- Deployment 名と Service 名は同じです。`OTEL_SERVICE_NAME` にもこの名前を使います。
- DB の migrate は、各サービスの initContainer で流します。DB が起動するまでは失敗して再試行されます。
- Ingress、Gateway API、ServiceMonitor は、既定では作りません。使う環境に合わせて有効にします。

## 必要なもの

- Kubernetes 1.27 以上。gRPC の probe を使うためです。
- Helm 3 以上
- 使う機能によっては、次も要ります。
  - Ingress: ingress-nginx などの Ingress コントローラー
  - Gateway API: Gateway API の CRD と、Gateway コントローラー（cilium など）
  - ServiceMonitor: Prometheus Operator（kube-prometheus-stack など）

## 入れる

```bash
helm install handson deploy/helm/handson -n handson --create-namespace
kubectl -n handson wait --for=condition=Ready pod --all --timeout=300s
```

画面を開くには、handson-legacy を port-forward します。

```bash
kubectl -n handson port-forward svc/handson-legacy 8080:8080
```

http://localhost:8080 を開き、`demo@example.com` / `demo-password` でログインします。

## よく使う設定

### ハンズオンの章と組み合わせる

cnd-handson の章で入れたものに合わせて、次の値を指定します。

```yaml
otel:
  # opentelemetry の章の collector
  endpoint: http://trace-collector-collector.default:4318
serviceMonitor:
  # prometheus の章の kube-prometheus-stack が拾う
  enabled: true
entry:
  ingress:
    # cluster-create の章の ingress-nginx
    enabled: true
    hosts:
      legacy: legacy.example.com
      modern: modern.example.com
```

### Gateway API で legacy と modern に振り分ける

HTTPRoute の重みで、2 つの版に振り分けます。cilium の章のトラフィック分割と同じ形です。

```yaml
entry:
  gateway:
    enabled: true
    weights:
      legacy: 90
      modern: 10
```

- chart が Gateway（`gatewayClassName: cilium`）も作ります。
- 章で作った Gateway を使うときは、`create: false` にして `parentRefs` を書きます。

```yaml
entry:
  gateway:
    enabled: true
    create: false
    parentRefs:
      - name: color-gw
```

Gateway に IP を振るには、LoadBalancer が要ります。kind では cilium の章の L2 Announcement（`manifest/l2announcement.yaml`）で振れます。

### Argo CD で入れる

Argo CD は `helm template` で manifest を作るので、chart がクラスタの Secret を読めません。そのままだと、DB のパスワードと idp の署名鍵が同期のたびに作り直されます。Secret を先に作り、`secret.create: false` にします。

```bash
kubectl create namespace handson
kubectl -n handson create secret generic handson-secrets \
  --from-literal=DB_PASSWORD="$(openssl rand -hex 16)" \
  --from-literal=IDP_SIGNING_KEY="$(openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 | base64 | tr -d '\n')" \
  --from-literal=IDP_DEMO_PASSWORD=demo-password
```

Application の例です。

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: handson
  namespace: argo-cd
spec:
  project: default
  source:
    repoURL: https://github.com/cloudnativedaysjp/cnd-handson-app
    targetRevision: main
    path: deploy/helm/handson
    helm:
      valuesObject:
        secret:
          create: false
  destination:
    server: https://kubernetes.default.svc
    namespace: handson
  syncPolicy:
    automated: {}
```

### 外の DB を使う

```yaml
postgres:
  enabled: false
  host: db.example.com
  port: 5432
  database: handson
  user: handson
  password: <パスワード>
```

## 値の一覧

| 値 | 既定 | 意味 |
|---|---|---|
| `image.registry` | `ghcr.io/cloudnativedaysjp/cnd-handson-app` | イメージの取得元 |
| `image.tag` | `latest` | handson-legacy / handson-modern 以外のイメージのタグ |
| `image.pullPolicy` | `IfNotPresent` | |
| `otel.endpoint` | `""` | OTLP（http/protobuf）の送り先。空なら送らない。送れなくてもサービスは止まらない |
| `serviceMonitor.enabled` | `false` | `/metrics` を集める ServiceMonitor を作る |
| `secret.create` | `true` | `handson-secrets` を作る。`false` なら先に作った Secret を使う |
| `idp.iss` / `idp.aud` | `http://handson-idp:8080` / `handson` | JWT の iss と aud。handson-legacy / handson-modern も検証に同じ値を使う |
| `idp.signingKey` | `""` | JWT の署名鍵。base64 の PEM。空なら生成する |
| `idp.demo.email` / `idp.demo.password` | `demo@example.com` / `demo-password` | デモユーザー。どちらかを空にすると作らない |
| `postgres.enabled` | `true` | chart の Postgres を作る |
| `postgres.image` | `postgres:17` | |
| `postgres.storage` | `1Gi` | Postgres の PVC の大きさ |
| `postgres.host` / `postgres.port` | `""` / `5432` | 外の DB の接続先。`enabled: false` のときだけ使う |
| `postgres.database` / `postgres.user` | `handson` / `handson` | |
| `postgres.password` | `""` | 空なら生成する |
| `entry.variants` | legacy（blue）、modern（green） | フロントエンドの版と、ログやトレースに付ける color |
| `entry.ingress.enabled` | `false` | handson-legacy / handson-modern の Ingress を作る |
| `entry.ingress.className` | `nginx` | |
| `entry.ingress.hosts` | `legacy.example.com` / `modern.example.com` | 版ごとのホスト名 |
| `entry.gateway.enabled` | `false` | handson-legacy / handson-modern の HTTPRoute を作る |
| `entry.gateway.create` | `true` | Gateway も作る。`false` なら `parentRefs` の Gateway につなぐ |
| `entry.gateway.className` | `cilium` | chart が作る Gateway の GatewayClass |
| `entry.gateway.parentRefs` | `[]` | `create: false` のときにつなぐ Gateway |
| `entry.gateway.hostnames` | `[]` | HTTPRoute のホスト名。空なら絞らない |
| `entry.gateway.weights` | legacy 100、modern 0 | 版ごとの重み |
| `services` | idp、project、task、column | サービスごとのコマンドと環境変数。ふつうは変えない |

## 秘密の値

`handson-secrets` に次の 3 つが入ります。

- `DB_PASSWORD`
- `IDP_SIGNING_KEY`
- `IDP_DEMO_PASSWORD`

値を指定しなかったものは、`helm install` のときに生成します。`helm upgrade` では既存の Secret から引き継ぐので、値は変わりません。

## 手元のイメージで kind に入れる

ghcr に公開する前のコードを試すときは、ビルドしたイメージを kind に読み込ませます。

```bash
kind create cluster --name handson-app
docker compose build idp-service project-service task-service column-service handson-legacy handson-modern
REG=ghcr.io/cloudnativedaysjp/cnd-handson-app
for s in idp project task column; do
  docker tag "${s}:latest" "$REG/${s}:dev" && kind load docker-image "$REG/${s}:dev" --name handson-app
done
for v in legacy modern; do
  docker tag "handson:${v}" "$REG/handson:${v}" && kind load docker-image "$REG/handson:${v}" --name handson-app
done
helm install handson deploy/helm/handson -n handson --create-namespace --set image.tag=dev
```

## 契約テストで確かめる

compose と同じポートに port-forward すれば、`make contract` をそのまま使えます。

```bash
kubectl -n handson port-forward svc/handson-idp 50051:50051 8082:8080 &
kubectl -n handson port-forward svc/handson-project 50053:50051 &
kubectl -n handson port-forward svc/handson-task 50055:50051 &
kubectl -n handson port-forward svc/handson-column 50056:50051 &
IDP_ISS=http://handson-idp:8080 IDP_AUD=handson make contract SVC="idp task project column"
```

## 消す

```bash
helm uninstall handson -n handson
```

Postgres の PVC は残ります。データも消すなら、PVC も消します。

```bash
kubectl -n handson delete pvc -l app=handson-postgres
```

## 今の制限

- ghcr のイメージは非公開なので、認証なしでは pull できません。公開されるまでは、手元のイメージを使います。
- handson-legacy / handson-modern のイメージのタグは、版の名前（`legacy` / `modern`）で固定です。`image.tag` は効きません。
