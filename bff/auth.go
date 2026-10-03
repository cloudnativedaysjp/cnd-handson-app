package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/singleflight"
)

// verifyFunc は JWT を検証し、sub（ユーザー ID）を返す
type verifyFunc func(ctx context.Context, token string) (string, error)

func newVerifier(jwksURL, issuer, audience string) verifyFunc {
	keys := &jwks{url: jwksURL, client: &http.Client{Timeout: 5 * time.Second, Transport: otelhttp.NewTransport(http.DefaultTransport)}}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
	)
	return func(ctx context.Context, token string) (string, error) {
		t, err := parser.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			kid, _ := t.Header["kid"].(string)
			return keys.key(ctx, kid)
		})
		if err != nil {
			return "", err
		}
		return t.Claims.GetSubject()
	}
}

// 知らない kid のトークンを送り続けて idp を叩かせないよう、取り直しの間隔を空ける
const jwksRefetchInterval = 30 * time.Second

type jwks struct {
	url    string
	client *http.Client
	group  singleflight.Group

	mu      sync.Mutex
	keys    map[string]*rsa.PublicKey
	fetched time.Time
}

// ロックはキャッシュの読み書きだけに使う。取得中に、キャッシュ済みの kid のリクエストを待たせないため
func (j *jwks) key(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	if k, ok, recent := j.cached(kid); ok {
		return k, nil
	} else if recent {
		return nil, fmt.Errorf("unknown kid %q", kid)
	}
	// 同時に来た取得は 1 回にまとめる。呼び出し元が切断しても取得は続け、ほかのリクエストに鍵を残す
	_, err, _ := j.group.Do("", func() (any, error) {
		keys, err := j.fetch(context.WithoutCancel(ctx))
		if err != nil {
			return nil, err
		}
		// 取得できたときだけ間隔を空ける。失敗（idp の起動前など）で 30 秒間すべて 401 にしないため
		j.mu.Lock()
		j.keys, j.fetched = keys, time.Now()
		j.mu.Unlock()
		return nil, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwks: %w", err)
	}
	if k, ok, _ := j.cached(kid); ok {
		return k, nil
	}
	return nil, fmt.Errorf("unknown kid %q", kid)
}

func (j *jwks) cached(kid string) (key *rsa.PublicKey, ok, recent bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	key, ok = j.keys[kid]
	return key, ok, time.Since(j.fetched) < jwksRefetchInterval
}

func (j *jwks) fetch(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, j.url, nil)
	if err != nil {
		return nil, err
	}
	res, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", res.StatusCode)
	}
	var body struct {
		Keys []struct{ Kty, Kid, N, E string }
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	keys := map[string]*rsa.PublicKey{}
	for _, k := range body.Keys {
		if k.Kty != "RSA" {
			continue
		}
		n, errN := base64.RawURLEncoding.DecodeString(k.N)
		e, errE := base64.RawURLEncoding.DecodeString(k.E)
		if err := errors.Join(errN, errE); err != nil {
			return nil, fmt.Errorf("key %q: %w", k.Kid, err)
		}
		keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
	}
	return keys, nil
}
