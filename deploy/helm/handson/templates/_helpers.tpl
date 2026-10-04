{{/* server と migrate で共通の、イメージ・コマンド・環境変数 */}}
{{- define "handson.container" -}}
{{- $v := .root.Values -}}
image: {{ printf "%s/%s:%s" $v.image.registry .key $v.image.tag }}
imagePullPolicy: {{ $v.image.pullPolicy }}
command: {{ toJson .svc.command }}
securityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop: ["ALL"]
env:
  - name: DB_HOST
    value: {{ ternary "handson-postgres" $v.postgres.host $v.postgres.enabled | quote }}
  - name: DB_PORT
    value: {{ ternary 5432 $v.postgres.port $v.postgres.enabled | quote }}
  - name: DB_DB
    value: {{ $v.postgres.database | quote }}
  - name: DB_USER
    value: {{ $v.postgres.user | quote }}
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: handson-secrets
        key: DB_PASSWORD
  {{- if eq .key "idp" }}
  - name: IDP_SIGNING_KEY
    valueFrom:
      secretKeyRef:
        name: handson-secrets
        key: IDP_SIGNING_KEY
  - name: IDP_ISS
    value: {{ $v.idp.iss | quote }}
  - name: IDP_AUD
    value: {{ $v.idp.aud | quote }}
  - name: IDP_DEMO_EMAIL
    value: {{ $v.idp.demo.email | quote }}
  - name: IDP_DEMO_PASSWORD
    valueFrom:
      secretKeyRef:
        name: handson-secrets
        key: IDP_DEMO_PASSWORD
  {{- end }}
  {{- range $k, $val := .svc.env }}
  - name: {{ $k }}
    value: {{ $val | quote }}
  {{- end }}
  - name: OTEL_SERVICE_NAME
    value: {{ .name }}
  - name: OTEL_RESOURCE_ATTRIBUTES
    value: service.namespace=handson
  - name: OTEL_EXPORTER_OTLP_PROTOCOL
    value: http/protobuf
  - name: OTEL_PROPAGATORS
    value: tracecontext,baggage
  {{- with $v.otel.endpoint }}
  - name: OTEL_EXPORTER_OTLP_ENDPOINT
    value: {{ . | quote }}
  {{- end }}
{{- end -}}
