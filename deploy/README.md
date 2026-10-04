# deploy

k8s に入れる Helm chart。`helm/handson` に全サービスが入っている。

## 入れる

```bash
helm install handson deploy/helm/handson -n handson --create-namespace
kubectl -n handson wait --for=condition=Ready pod --all --timeout=300s
```

- イメージは `ghcr.io/cloudnativedaysjp/cnd-handson-app/<サービス名>:latest` から取る
- DB のパスワードと idp の署名鍵は、install 時に生成する。upgrade では同じ値を使い続ける
- 入口のイメージは `<registry>/handson:legacy` と `<registry>/handson:modern`。版はビルド時に決まるので、`image.tag` は効かない

よく変える値:

| 値 | 意味 |
|---|---|
| `otel.endpoint` | トレースとメトリクスの送り先（OTLP http/protobuf）。例: `http://trace-collector-collector.default:4318` |
| `image.registry` / `image.tag` | イメージの取得元 |
| `postgres.enabled=false` と `postgres.host` | chart の Postgres を使わず、外の DB に繋ぐ |
| `entry.ingress.enabled` と `entry.ingress.hosts` | 入口の Ingress を作る。ホスト名は版ごとに決める |
| `serviceMonitor.enabled` | Prometheus Operator が `/metrics` を集める ServiceMonitor を作る |
| `secret.create=false` | Secret を作らず、先に作った `handson-secrets` を使う。Argo CD で入れるときに使う |

## Argo CD で入れる

Argo CD は `helm template` で描画するので、`lookup` が効かない。生成した値が同期のたびに変わってしまうので、Secret は先に作り、`secret.create=false` にする。

```bash
kubectl create namespace handson
kubectl -n handson create secret generic handson-secrets \
  --from-literal=DB_PASSWORD="$(openssl rand -hex 16)" \
  --from-literal=IDP_SIGNING_KEY="$(openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 | base64 | tr -d '\n')" \
  --from-literal=IDP_DEMO_PASSWORD=demo-password
```

## 手元のイメージで kind に入れる

ghcr に公開する前のコードを試すときは、ビルドしたイメージを kind に読み込ませる。

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

## 契約テストを流す

compose と同じポートに port-forward すれば、`make contract` をそのまま使える。

```bash
kubectl -n handson port-forward svc/handson-idp 50051:50051 8082:8080 &
kubectl -n handson port-forward svc/handson-project 50053:50051 &
kubectl -n handson port-forward svc/handson-task 50055:50051 &
kubectl -n handson port-forward svc/handson-column 50056:50051 &
IDP_ISS=http://handson-idp:8080 IDP_AUD=handson make contract SVC="idp task project column"
```
